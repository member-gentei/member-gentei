package cmd

import (
	"context"

	"github.com/member-gentei/member-gentei/gentei/apis"
	"github.com/member-gentei/member-gentei/gentei/membership"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var (
	flagCheckUserID      uint64
	flagCheckChannelIDs  []string
	flagCheckRefreshOnly bool
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check users' memberships.",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ctx           = context.Background()
			db            = mustOpenDB(ctx)
			discordConfig = getDiscordConfig()
			ytConfig      = getYouTubeConfig()
			results       *membership.CheckResultSet
			options       *membership.CheckForUserOptions
			err           error
		)
		for _, chID := range flagCheckChannelIDs {
			if options == nil {
				options = &membership.CheckForUserOptions{}
			}
			options.ChannelIDs = append(options.ChannelIDs, chID)
		}
		if flagCheckUserID != 0 {
			// refresh guilds for user
			var (
				ts    oauth2.TokenSource
				token *oauth2.Token
			)
			ts, err = apis.GetRefreshingDiscordTokenSource(ctx, db, discordConfig, flagCheckUserID)
			if err != nil {
				log.Fatal().Err(err).Msg("error getting Discord TokenSource for single user")
			}
			token, err = ts.Token()
			if err != nil {
				log.Fatal().Err(err).Msg("error getting Discord token for single user")
			}
			_, _, err = membership.RefreshUserGuildEdges(ctx, db, token, flagCheckUserID)
			if err != nil {
				log.Fatal().Err(err).Msg("error refreshing Discord Guild edges for single user")
			}
			if flagCheckRefreshOnly {
				return
			}
			// check single user ID
			results, err = membership.CheckForUser(ctx, db, ytConfig, flagCheckUserID, options)
			if err != nil {
				log.Fatal().Err(err).Msg("error checking memberships for user")
			}
			log.Info().
				Interface("checkResults", results).
				Uint64("userID", flagCheckUserID).
				Msg("check results")
			err = membership.SaveMemberships(ctx, db, flagCheckUserID, results)
		} else {
			// otherwise, refresh all guilds and check all stale
			err = membership.RefreshAllUserGuildEdges(ctx, db, discordConfig)
			if err != nil {
				log.Fatal().Err(err).Msg("error refreshing Discord Guild edges")
			}
			if flagCheckRefreshOnly {
				return
			}
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
	flags.BoolVar(&flagCheckRefreshOnly, "refresh-only", false, "only refresh Discord information")
}
