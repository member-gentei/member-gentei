package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/bwmarrin/discordgo"
	"github.com/member-gentei/member-gentei/gentei/async"
	"github.com/member-gentei/member-gentei/gentei/bot/roles"
	"github.com/member-gentei/member-gentei/gentei/bot/templates"
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
			Uint64("guildID", guildID).
			Str("guildName", gc.Name).
			Logger()
		logger.Info().Msg("joined Guild")
		// update Guild info in ~5s to avoid clashing with first time registration
		go func() {
			<-time.NewTimer(time.Second * 5).C
			err = b.db.Guild.UpdateOneID(guildID).
				SetName(gc.Name).
				SetIconHash(gc.Icon).
				Exec(context.Background())
			if err != nil {
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
			Uint64("guildID", guildID).
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

func (b *DiscordBot) StartPSApplier(parentCtx context.Context, sub *pubsub.Subscription) {
	var (
		pCtx, cancel = context.WithCancel(parentCtx)
	)
	b.cancelPSApplier = cancel
	go func() {
		defer cancel()
		err := sub.Receive(pCtx, func(ctx context.Context, m *pubsub.Message) {
			if typeAttribute := m.Attributes["type"]; typeAttribute != string(async.ApplyMembershipType) {
				log.Warn().Str("typeAttribute", typeAttribute).Msg("non apply-membership message made it past the filter?")
				m.Ack()
				return
			}
			var message async.ApplyMembershipPSMessage
			err := json.Unmarshal(m.Data, &message)
			if err != nil {
				log.Warn().Str("data", string(m.Data)).Msg("acking message that cannot be decoded as JSON")
				m.Ack()
				return
			}
			if message.DeleteUserID != "" {
				var (
					userIDStr   = message.DeleteUserID.String()
					userID, err = strconv.ParseUint(message.DeleteUserID.String(), 10, 64)
				)
				if err != nil {
					log.Err(err).
						Str("unparsedUserID", userIDStr).
						Msg("error decoding UserID as uint64")
					m.Ack()
				}
				logger := log.With().Uint64("userID", userID).Logger()
				err = b.revokeMembershipsByUserID(ctx, userID)
				if err != nil {
					logger.Err(err).Msg("error revoking all memberships before deletion")
					return
				}
				// now actually delete the user
				err = b.db.User.DeleteOneID(userID).Exec(ctx)
				if err != nil {
					return
				}
				m.Ack()
				// best-effort attempt at sending the user deletion DM
				ch, err := b.session.UserChannelCreate(userIDStr)
				if err != nil {
					logger.Err(err).Msg("error creating UserChannel to inform of deletion")
					return
				}
				msg, err := b.session.ChannelMessageSend(ch.ID, templates.PlaintextUserDeleted)
				if err != nil {
					logger.Err(err).
						Msg("error sending deletion confirmation message")
				} else {
					logger.Info().
						Interface("messageMetadata", msg).
						Msg("sent deletion confirmation message")
				}
				return
			} else if message.Gained {
				err = b.grantMemberships(ctx, b.db, message.UserMembershipID)
				if err != nil {
					log.Err(err).Msg("error granting memberships")
					return
				}
			} else {
				err = b.revokeMemberships(ctx, b.db, message.UserMembershipID)
				if err != nil {
					log.Err(err).Msg("error revoking memberships")
					return
				}
			}
			m.Ack()
		})
		if err != nil {
			log.Err(err).Msg("bot PSApplier crashed?")
		}
	}()
}

func (b *DiscordBot) Close() error {
	if b.cancelPSApplier != nil {
		b.cancelPSApplier()
	}
	return b.session.Close()
}

func (b *DiscordBot) applyRole(ctx context.Context, guildID, roleID, userID uint64, add bool) error {
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
	)
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
	log.Err(err).
		Uint64("guildID", guildID).Uint64("userID", userID).Uint64("roleID", roleID).
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
