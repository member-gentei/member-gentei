package cmd

import (
	"context"

	"github.com/member-gentei/member-gentei/gentei/membership"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	flagCheckUserID     uint64
	flagCheckChannelIDs []string
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check users' memberships",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if flagYouTubeClientID == "" {
			log.Fatal().Msgf("env var %s must not be empty", envNameYouTubeClientID)
		}
		if flagYouTubeClientSecret == "" {
			log.Fatal().Msgf("env var %s must not be empty", envNameYouTubeClientSecret)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ctx      = context.Background()
			db       = mustOpenDB(ctx)
			ytConfig = getYouTubeConfig()
			results  *membership.CheckResultSet
			err      error
		)
		if flagCheckUserID != 0 {
			// check single user ID
			results, err = membership.CheckForUser(ctx, db, ytConfig, flagCheckUserID, nil)
			if err != nil {
				log.Fatal().Err(err).Msg("error checking memberships for user")
			}
			log.Info().
				Interface("checkResults", results).
				Uint64("userID", flagCheckUserID).
				Msg("check results")
			err = membership.SaveMemberships(ctx, db, flagCheckUserID, results)
		} else {
			// otherwise, check all stale
			err = membership.CheckStale(ctx, db, ytConfig, nil)
		}
		if err != nil {
			log.Fatal().Err(err).Msg("error saving memberships")
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	flags := checkCmd.Flags()
	flags.Uint64Var(&flagCheckUserID, "uid", 0, "check only this user")
	flags.StringSliceVar(&flagCheckChannelIDs, "channel", nil, "check user(s) against these channels")
}
