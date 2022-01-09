package bot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/rs/zerolog/log"
)

type DiscordBot struct {
	session *discordgo.Session
	db      *ent.Client
}

func New(db *ent.Client, token string) (*DiscordBot, error) {
	session, err := discordgo.New(fmt.Sprintf("Bot %s", token))
	if err != nil {
		return nil, fmt.Errorf("error creating discordgo session: %w", err)
	}
	return &DiscordBot{
		session: session,
		db:      db,
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
		case "info":
			b.handleInfo(ctx, i)
		case "manage":
			b.handleManage(ctx, i)
		}
	})
	if err = b.session.Open(); err != nil {
		return fmt.Errorf("error opening discordgo session: %w", err)
	}
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

func (b *DiscordBot) Close() error {
	return b.session.Close()
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
