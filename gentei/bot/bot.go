package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"cloud.google.com/go/pubsub"
	"github.com/bwmarrin/discordgo"
	"github.com/member-gentei/member-gentei/gentei/async"
	"github.com/member-gentei/member-gentei/gentei/bot/roles"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/membership"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

type DiscordBot struct {
	session *discordgo.Session

	db              *ent.Client
	rut             *roles.RoleUpdateTracker
	qch             *membership.QueuedChangeHandler
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

func (b *DiscordBot) Start(prod, upsertCommands bool) (err error) {
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
	// load membership.ChangeHandler
	qch, changeHandler := membership.NewQueuedChangeHandler(10)
	membership.HookMembershipChanges(b.db, changeHandler)
	b.qch = qch
	if err = b.session.Open(); err != nil {
		return fmt.Errorf("error opening discordgo session: %w", err)
	}
	// load pubsub
	// load early access commands
	if !upsertCommands {
		log.Info().Msg("skipping command upsert")
		return nil
	}
	for _, guildID := range earlyAccessGuilds {
		log.Info().Str("guildID", guildID).Msg("loading early access command")
		_, err = b.session.ApplicationCommandCreate(b.session.State.User.ID, guildID, earlyAccessCommand)
		if err != nil {
			return fmt.Errorf("error loading early access command to guild '%s': %w", guildID, err)
		}
	}
	if globalCommand.Name == "" {
		return nil
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
			if message.Gained {
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
		applyCtx, cancelApplyCtx = context.WithCancel(ctx)
		guildIDStr               = strconv.FormatUint(guildID, 10)
		roleIDStr                = strconv.FormatUint(roleID, 10)
		userIDStr                = strconv.FormatUint(userID, 10)
	)
	b.rut.TrackHook(guildIDStr, userIDStr, func(gmu *discordgo.GuildMemberUpdate) (remove bool) {
		if add {
			for _, roleID := range gmu.Roles {
				if roleID == roleIDStr {
					cancelApplyCtx()
					return true
				}
			}
		} else {
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
	err := result.Error
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
