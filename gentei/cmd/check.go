package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"cloud.google.com/go/pubsub"
	"github.com/bwmarrin/discordgo"
	"github.com/member-gentei/member-gentei/gentei/apis"
	"github.com/member-gentei/member-gentei/gentei/async"
	"github.com/member-gentei/member-gentei/gentei/ent/predicate"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/member-gentei/member-gentei/gentei/membership"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var (
	flagCheckUserID      uint64
	flagCheckChannelIDs  []string
	flagCheckRefreshOnly bool
	flagCheckNoEnforce   bool
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check users' memberships.",
	PreRun: func(cmd *cobra.Command, args []string) {
		if flagPubSubTopic == "" && !flagCheckNoEnforce {
			log.Fatal().Msgf("env var %s must not be empty", envNamePubSubTopic)
		}
	},
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
		ps, err := pubsub.NewClient(ctx, flagGCPProjectID)
		if err != nil {
			log.Fatal().Err(err).Msg("error calling pubsub.NewClient")
		}
		asyncTopic := ps.Topic(flagPubSubTopic)
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
				var restErr *discordgo.RESTError
				if errors.As(err, &restErr) {
					log.Warn().Err(err).Msg("error using Discord token for user")
				}
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
				Str("userID", strconv.FormatUint(flagCheckUserID, 10)).
				Msg("check results")
			err = membership.SaveMemberships(ctx, db, flagCheckUserID, results)
		} else {
			// otherwise, refresh all guilds and check all stale users
			var (
				failedUserIDs []uint64
				userCount     int
			)
			failedUserIDs, userCount, err = membership.RefreshAllUserGuildEdges(ctx, db, discordConfig)
			if err != nil {
				log.Fatal().Err(err).Msg("error refreshing Discord Guild edges")
			}
			log.Info().
				Int("total", userCount).
				Int("succeeded", userCount-len(failedUserIDs)).
				Msg("refreshed guild edges")
			if !flagCheckNoEnforce {
				for _, userID := range failedUserIDs {
					var (
						userIDStr = strconv.FormatUint(userID, 10)
						logger    = log.With().Str("userID", userIDStr).Logger()
					)
					err = async.PublishGeneralMessage(ctx, asyncTopic, async.GeneralPSMessage{
						UserDelete: &async.UserDeleteMessage{
							UserID: json.Number(userIDStr),
							Reason: "Discord token invalid",
						},
					})
					if err != nil {
						logger.Fatal().Err(err).Msg("error publishing delete message")
					} else {
						logger.Info().Msg("issued delete for user")
					}
				}
			}
			if flagCheckRefreshOnly {
				return
			}
			var excludeToBeDeleted []predicate.User
			if len(failedUserIDs) > 0 {
				excludeToBeDeleted = append(excludeToBeDeleted, user.IDNotIn(failedUserIDs...))
			}
			err = membership.CheckStale(ctx, db, ytConfig, &membership.CheckStaleOptions{
				StaleThreshold:           membership.DefaultCheckStaleOptions.StaleThreshold,
				AdditionalUserPredicates: excludeToBeDeleted,
			})
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
	flags.BoolVar(&flagCheckNoEnforce, "no-enforce", false, "do not effect membership changes")
}
