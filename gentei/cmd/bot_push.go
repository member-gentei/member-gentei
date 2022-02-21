package cmd

import (
	"context"
	"errors"

	"github.com/member-gentei/member-gentei/gentei/bot"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	flagPushEA     bool
	flagPushGlobal bool
)

// botPushCmd represents the push command
var botPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Pushes new commands, global or otherwise.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !(flagPushEA && flagPushGlobal) {
			return errors.New("must specify at least one of --early-access and --global")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ctx = context.Background()
			db  = mustOpenDB(ctx)
		)
		dgBot, err := bot.New(db, flagBotToken, getYouTubeConfig())
		if err != nil {
			log.Fatal().Err(err).Msg("error creating Discord bot")
		}
		err = dgBot.PushCommands(flagPushGlobal, flagPushEA)
		if err != nil {
			log.Fatal().Err(err).Msg("error pushing commands")
		}
	},
}

func init() {
	botCmd.AddCommand(botPushCmd)
	flags := botPushCmd.Flags()
	flags.BoolVar(&flagPushEA, "early-access", false, "push early access commands")
	flags.BoolVar(&flagPushGlobal, "global", false, "push global/prod commands")
}
