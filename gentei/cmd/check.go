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
	"github.com/member-gentei/member-gentei/gentei/membership"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var (
	flagCheckDisabled        bool
	flagCheckUserID          uint64
	flagCheckChannelIDs      []string
	flagCheckRefreshAllUsers bool
	flagCheckRefreshOnly     bool
	flagCheckNoEnforce       bool
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
		if options != nil && flagCheckDisabled {
			options.CheckDisabledChannels = true
		}
		if flagCheckUserID != 0 {
			log.Debug().Msg("checking single user")
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
			if errors.Is(err, apis.ErrInvalidGrant) {
				log.Fatal().Err(err).Msg("token expired/revoked, user should be deleted")
			} else if err != nil {
				log.Fatal().Err(err).Msg("error checking memberships for user")
			}
			log.Info().
				Interface("checkResults", results).
				Str("userID", strconv.FormatUint(flagCheckUserID, 10)).
				Msg("check results")
			if flagCheckNoEnforce {
				return
			}
			membership.HookMembershipChanges(db, async.NewPubSubMembershipChangeHandler(ctx, asyncTopic))
			err = membership.SaveMemberships(ctx, db, flagCheckUserID, results)
			if err != nil {
				log.Err(err).Msg("error saving memberships for single user check")
			}
		} else if options != nil {
			log.Fatal().Msg("narrowing options not supported without --uid")
		} else {
			// otherwise, refresh all guilds and check all stale users
			userDeleteChan := make(chan uint64, 10)
			go func() {
				defer close(userDeleteChan)
				if flagCheckNoEnforce {
					for range userDeleteChan {
					}
				} else {
					for userID := range userDeleteChan {
						var (
							userIDStr = strconv.FormatUint(userID, 10)
							logger    = log.With().Str("userID", userIDStr).Logger()
						)
						err = async.PublishGeneralMessage(ctx, asyncTopic, async.GeneralPSMessage{
							UserDelete: &async.DeleteUserMessage{
								UserID: json.Number(userIDStr),
								Reason: "Discord token invalid/expired",
							},
						})
						if err != nil {
							logger.Fatal().Err(err).Msg("error publishing delete message")
						} else {
							logger.Info().Msg("issued delete for user")
						}
					}
				}
			}()
			opts := membership.DefaultPerformCheckOptions()
			if flagCheckRefreshAllUsers {
				log.Info().Msg("refreshing all UserGuildEdges")
				opts.StaleThreshold = 0
			} else {
				log.Info().Msg("refreshing stale UserGuildEdges")
			}
			opts.Enforce = !flagCheckNoEnforce
			err := membership.PerformCheckBatches(ctx, db, discordConfig, ytConfig, userDeleteChan, opts)
			if err != nil {
				log.Fatal().Err(err).Msg("error performing checks")
			}
			if flagCheckRefreshOnly {
				return
			}
			err = async.PublishApplyMembershipMessage(ctx, asyncTopic, async.ApplyMembershipPSMessage{
				EnforceAll: &async.EnforceAllMessage{
					DryRun: flagCheckNoEnforce,
					Reason: "periodic role enforcement",
				},
			})
			if err != nil {
				log.Fatal().Err(err).Msg("error publishing apply message")
			}
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
	flags.BoolVar(&flagCheckDisabled, "check-disabled", false, "check against disabled channels")
	flags.BoolVar(&flagCheckRefreshAllUsers, "refresh-all-users", false, "refresh tokens for non-stale users")
	flags.BoolVar(&flagCheckRefreshOnly, "refresh-only", false, "only refresh Discord information")
	flags.BoolVar(&flagCheckNoEnforce, "no-enforce", false, "do not effect membership changes")
}
