/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/member-gentei/member-gentei/gentei/bot"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	envNameDiscordBotToken = "DISCORD_BOT_TOKEN"
)

var (
	flagBotToken         string
	flagBotProd          bool
	flagNoUpsertCommands bool
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
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ctx = context.Background()
			db  = mustOpenDB(ctx)
		)
		genteiBot, err := bot.New(db, flagBotToken)
		if err != nil {
			log.Fatal().Err(err).Msg("erorr creating bot.DiscordBot")
		}
		if err = genteiBot.Start(flagBotProd, !flagNoUpsertCommands); err != nil {
			log.Fatal().Err(err).Msg("error starting bot")
		}
		log.Info().Msg("bot started")
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
	flags.BoolVar(&flagNoUpsertCommands, "no-upsert-commands", false, "disable pushing new commands")
}
