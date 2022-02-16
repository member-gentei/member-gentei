package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"cloud.google.com/go/pubsub"
	"github.com/member-gentei/member-gentei/gentei/bot"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	envNameDiscordBotToken = "DISCORD_BOT_TOKEN"
)

var (
	flagBotToken       string
	flagBotProd        bool
	flagUpsertCommands bool
)

// botCmd represents the bot command
var botCmd = &cobra.Command{
	Use:   "bot",
	Short: "Runs the Discord bot",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		flagBotToken = os.Getenv(envNameDiscordBotToken)
		if flagBotToken == "" {
			return fmt.Errorf("must specify env var '%s'", envNameDiscordBotToken)
		}
		if flagYouTubeClientID == "" {
			log.Fatal().Msgf("env var %s must not be empty", envNameYouTubeClientID)
		}
		if flagYouTubeClientSecret == "" {
			log.Fatal().Msgf("env var %s must not be empty", envNameYouTubeClientSecret)
		}
		if flagPubSubSubscription == "" {
			log.Fatal().Msgf("env var %s must not be empty", envNamePubSubSubscription)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ctx = context.Background()
			db  = mustOpenDB(ctx)
		)
		ps, err := pubsub.NewClient(ctx, flagGCPProjectID)
		if err != nil {
			log.Fatal().Err(err).Msg("error calling pubsub.NewClient")
		}
		genteiBot, err := bot.New(db, flagBotToken, getYouTubeConfig())
		if err != nil {
			log.Fatal().Err(err).Msg("error creating bot.DiscordBot")
		}
		if err = genteiBot.Start(flagBotProd, flagUpsertCommands); err != nil {
			log.Fatal().Err(err).Msg("error starting bot")
		}
		log.Info().Msg("bot started")
		genteiBot.StartPSApplier(ctx, ps.Subscription(flagPubSubSubscription))
		defer genteiBot.Close()
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)
		<-stop
		log.Info().Msg("gracefully shutting down")
	},
}

func init() {
	rootCmd.AddCommand(botCmd)
	flags := botCmd.Flags()
	flags.BoolVar(&flagBotProd, "prod", false, "pushes global commands")
	flags.BoolVar(&flagUpsertCommands, "upsert-commands", false, "push new commands")
}
