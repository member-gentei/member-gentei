package bot

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/member-gentei/member-gentei/gentei/bot/roles"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

type DiscordBot struct {
	session *discordgo.Session

	db              *ent.Client
	rut             *roles.RoleUpdateTracker
	cancelPSApplier context.CancelFunc
	youTubeConfig   *oauth2.Config

	// roleRWMutex is held 'W' by role-wide operations, like daily enforcement. 'R' is held by all other operations that 'W' would overrun or interrupt.
	roleRWMutex *roles.DefaultMapRWMutex
}

func New(db *ent.Client, token string, youTubeConfig *oauth2.Config) (*DiscordBot, error) {
	session, err := discordgo.New(fmt.Sprintf("Bot %s", token))
	if err != nil {
		return nil, fmt.Errorf("error creating discordgo session: %w", err)
	}
	rut := roles.NewRoleUpdateTracker(session)
	return &DiscordBot{
		session:       session,
		db:            db,
		rut:           rut,
		youTubeConfig: youTubeConfig,
		roleRWMutex:   roles.NewDefaultMapRWMutex(),
	}, nil
}

func (b *DiscordBot) Start(prod bool) (err error) {
	// register handlers on bot start
	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var (
			ctx, cancel    = context.WithCancel(context.Background())
			appCommandData = i.ApplicationCommandData()
		)
		defer cancel()
		log.Debug().Interface("appCommandData", appCommandData).Send()
		// subcommand
		subcommand := appCommandData.Options[0]
		switch subcommand.Name {
		case "check":
			b.handleCheck(ctx, i)
		case "info":
			b.handleInfo(ctx, i)
		case "manage":
			b.handleManage(ctx, i)
		}
	})
	// guild metadata updates
	b.session.AddHandler(func(s *discordgo.Session, gc *discordgo.GuildCreate) {
		guildID, err := strconv.ParseUint(gc.ID, 10, 64)
		if err != nil {
			log.Err(err).Str("unparsedGuildID", gc.ID).Msg("error parsing joined guild ID as uint64")
			return
		}
		logger := log.With().
			Str("guildID", strconv.FormatUint(guildID, 10)).
			Str("guildName", gc.Name).
			Logger()
		logger.Info().Msg("joined Guild")
		// update Guild info opportunistically
		go func() {
			<-time.NewTimer(time.Second * 5).C
			err = b.db.Guild.UpdateOneID(guildID).
				SetName(gc.Name).
				SetIconHash(gc.Icon).
				Exec(context.Background())
			if err != nil && !ent.IsNotFound(err) {
				logger.Err(err).Msg("error updating on GUILD_CREATE")
				return
			}
		}()
	})
	b.session.AddHandler(func(s *discordgo.Session, gu *discordgo.GuildUpdate) {
		guildID, err := strconv.ParseUint(gu.ID, 10, 64)
		if err != nil {
			log.Err(err).Str("unparsedGuildID", gu.ID).Msg("error parsing joined guild ID as uint64")
			return
		}
		logger := log.With().
			Str("guildID", strconv.FormatUint(guildID, 10)).
			Str("guildName", gu.Name).
			Logger()
		// update if guild and info exists
		err = b.db.Guild.UpdateOneID(guildID).
			SetName(gu.Name).
			SetIconHash(gu.Icon).
			Exec(context.Background())
		if ent.IsNotFound(err) {
			logger.Warn().Msg("got update for Guild not in database")
		} else if err != nil {
			logger.Err(err).Msg("error updating GUILD_UPDATE")
		} else {
			logger.Info().Msg("updated Guild")
		}
	})
	// register intents (new for v8 gateway)
	b.session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMembers
	if err = b.session.Open(); err != nil {
		return fmt.Errorf("error opening discordgo session: %w", err)
	}
	return nil
}

func (b *DiscordBot) Close() error {
	if b.cancelPSApplier != nil {
		b.cancelPSApplier()
	}
	return b.session.Close()
}

func (b *DiscordBot) applyRole(ctx context.Context, guildID, roleID, userID uint64, add bool, auditReason string) error {
	var (
		guildIDStr = strconv.FormatUint(guildID, 10)
		roleIDStr  = strconv.FormatUint(roleID, 10)
		userIDStr  = strconv.FormatUint(userID, 10)
	)
	// first, check if we even need to do this
	member, err := b.session.GuildMember(guildIDStr, userIDStr)
	if err != nil {
		return fmt.Errorf("error calling GuildMember: %w", err)
	}
	var roleExists bool
	for _, existingRoleID := range member.Roles {
		if existingRoleID == roleIDStr {
			roleExists = true
			break
		}
	}
	if (roleExists && add) || (!roleExists && !add) {
		log.Debug().Msg("no change required")
		return nil
	}
	var (
		applyCtx, cancelApplyCtx = context.WithCancel(ctx)
		mutex                    = b.roleRWMutex.GetOrCreate(roleIDStr)
	)
	mutex.RLock()
	defer mutex.RUnlock()
	b.rut.TrackHook(guildIDStr, userIDStr, func(gmu *discordgo.GuildMemberUpdate) (removeHook bool) {
		if add {
			// check this update for the target role that should exist
			for _, roleID := range gmu.Roles {
				if roleID == roleIDStr {
					cancelApplyCtx()
					return true
				}
			}
		} else {
			// check that the role does not exist
			for _, roleID := range gmu.Roles {
				if roleID == roleIDStr {
					return false
				}
			}
			cancelApplyCtx()
			return true
		}
		return
	})
	result := <-roles.ApplyRole(applyCtx, b.session, guildID, userID, roleID, add)
	err = result.Error
	if errors.Is(err, context.Canceled) {
		err = nil
	}
	if err == nil {
		b.auditLog(ctx, guildID, userID, roleID, add, auditReason)
	}
	log.Err(err).
		Str("guildID", strconv.FormatUint(guildID, 10)).Str("userID", strconv.FormatUint(userID, 10)).Str("roleID", strconv.FormatUint(roleID, 10)).
		Int("attempts", result.Attempts).
		Msg("role apply attempt finished")
	return err
}

func getMessageAttributionIDs(i *discordgo.InteractionCreate) (guildID, userID uint64, err error) {
	userID, err = strconv.ParseUint(i.Member.User.ID, 10, 64)
	if err != nil {
		err = fmt.Errorf("error decoding Member.User.ID as uint64: %w", err)
		return
	}
	guildID, err = strconv.ParseUint(i.GuildID, 10, 64)
	if err != nil {
		err = fmt.Errorf("error decoding GuildID as uint64: %w", err)
		return
	}
	return
}
