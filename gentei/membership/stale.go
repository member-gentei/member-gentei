package membership

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/member-gentei/member-gentei/gentei/apis"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/member-gentei/member-gentei/gentei/ent/predicate"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
)

type CheckStaleOptions struct {
	// StaleThreshold is used in a <= comparison to the last stored membership check time.
	StaleThreshold time.Duration
	// UserPredicates overrides predicates that consider a user's memberships stale.
	UserPredicates []predicate.User
	// AdditionalUserPredicates allows specifying additional constraints on user checks.
	AdditionalUserPredicates []predicate.User
	// NoSave skips saving UserMembership edges.
	NoSave bool
	// MaxWorkers sets the maximum amount of users to check and process simultaneously.
	MaxWorkers int
	// The maximum amount of users to check
	TotalLimit int
}

var DefaultCheckStaleOptions = &CheckStaleOptions{
	StaleThreshold: time.Hour * 12,
	MaxWorkers:     100,
}

func CheckStale(
	ctx context.Context,
	db *ent.Client,
	youtubeConfig *oauth2.Config,
	options *CheckStaleOptions,
) error {
	if options == nil {
		options = DefaultCheckStaleOptions
	}
	staleThreshold := options.StaleThreshold
	if options.StaleThreshold > 0 {
		staleThreshold *= -1
	}
	staleUserPredicates := options.UserPredicates
	if staleUserPredicates != nil {
		log.Info().Msg("stale user predicates overriden")
	} else {
		staleUserPredicates = []predicate.User{
			user.HasGuildsWith(
				guild.HasYoutubeTalents(),
			),
			user.YoutubeIDNotNil(),
			user.LastCheckLTE(time.Now().Add(staleThreshold)),
		}
	}
	if options.AdditionalUserPredicates != nil {
		staleUserPredicates = append(staleUserPredicates, options.AdditionalUserPredicates...)
		log.Info().Msg("appending additional stale user predicates")
	}
	var (
		totalStaleCount int
		eg              errgroup.Group
	)
	if options.MaxWorkers > 0 {
		eg.SetLimit(options.MaxWorkers)
	} else {
		eg.SetLimit(DefaultCheckStaleOptions.MaxWorkers)
	}
	log.Info().Msg("beginning refresh of stale users")
	for {
		staleUserIDs, err := db.User.Query().
			Where(staleUserPredicates...).
			Limit(1000).
			IDs(ctx)
		if err != nil {
			return err
		}
		if len(staleUserIDs) == 0 {
			break
		}
		totalStaleCount += len(staleUserIDs)
		for _, u := range staleUserIDs {
			userID := u // capture
			eg.Go(func() error {
				// TODO: https://github.com/member-gentei/member-gentei/issues/92
				results, err := CheckForUser(ctx, db, youtubeConfig, userID, nil)
				if err != nil {
					return fmt.Errorf("error checking memberships for user '%d': %w", userID, err)
				}
				if options.NoSave {
					return nil
				}
				err = SaveMemberships(ctx, db, userID, results)
				if err != nil {
					return fmt.Errorf("error saving memberships for user '%d': %w", userID, err)
				}
				return nil
			})
		}
		if err := eg.Wait(); err != nil {
			return err
		}
		log.Info().Int("count", totalStaleCount).Msg("refreshed stale user batch of <=1000")
		if options.TotalLimit > 0 && totalStaleCount > options.TotalLimit {
			log.Info().Int("limit", options.TotalLimit).Int("count", totalStaleCount).Msg("reached CheckStale limit, returning early")
			return nil
		}
	}
	log.Info().Int("count", totalStaleCount).Msg("refreshed stale users")
	return nil
}

// RefreshAllUserGuildEdges refreshes guild edges for all registered users. Returns a slice of userIDs that could not be refreshed and a count of all users.
func RefreshAllUserGuildEdges(ctx context.Context, db *ent.Client, discordConfig *oauth2.Config) ([]uint64, int, error) {
	// refresh everyone's tokens
	failed, _, total, err := refreshUserGuildEdgesWithPredicates(ctx, db, discordConfig, -1, nil)
	return failed, total, err
}

// RefreshStaleUserGuildEdgesrefreshes guild edges for all registered users below a freshness threshold. Returns a slice of userIDs that could not be refreshed and a count of all users.
func RefreshStaleUserGuildEdges(ctx context.Context, db *ent.Client, discordConfig *oauth2.Config, staleThreshold time.Duration) ([]uint64, int, error) {
	staleBefore := time.Now().Add(-staleThreshold)
	failed, _, total, err := refreshUserGuildEdgesWithPredicates(ctx, db, discordConfig, -1, user.LastCheckLT(staleBefore))
	return failed, total, err
}

// refreshUserGuildEdgesWithPredicates is the inner implementation of all stale refreshes. Returns a slice of userIDs that could not be refreshed, a slice of those that succeeded, and a count of all users.
func refreshUserGuildEdgesWithPredicates(ctx context.Context, db *ent.Client, discordConfig *oauth2.Config, max int, predicates ...predicate.User) ([]uint64, []uint64, int, error) {
	var (
		userTokensInvalid []uint64
		userTokensValid   []uint64
		totalCount        int
		after             uint64
		userTokensMutex   = &sync.Mutex{}
	)
	totalCount, err := db.User.Query().
		Where(predicates...).
		Count(ctx)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("error getting total count of stale users: %w", err)
	}
	log.Info().Int("total", totalCount).Msg("current stale user count")
	const defaultPageSize = 400
	pageSize := defaultPageSize
	if max > 0 && max < defaultPageSize {
		pageSize = max
	}
	processedCounter := &atomic.Int64{}
	for {
		userIDs, err := db.User.Query().
			Where(append(
				predicates,
				user.IDGT(after),
			)...).
			Order(ent.Asc(user.FieldID)).
			Limit(pageSize).
			IDs(ctx)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("error paginating user IDs: %w", err)
		}
		var eGroup errgroup.Group
		eGroup.SetLimit(10) // I set this to 100 one time and got my airbnb IP banned, so don't go that high
		for i := range userIDs {
			userID := userIDs[i]
			eGroup.Go(func() error {
				defer processedCounter.Add(1)
				logger := log.With().Str("userID", strconv.FormatUint(userID, 10)).Logger()
				logger.Debug().Msg("getting Discord token for refresh")
				ts, err := apis.GetRefreshingDiscordTokenSource(ctx, db, discordConfig, userID)
				if err != nil {
					logger.Warn().Err(err).Msg("error creating Discord TokenSource, skipping")
					return nil
				}
				token, err := ts.Token()
				if err != nil {
					if strings.Contains(err.Error(), "429") {
						logger.Warn().Err(err).Msg("Encountered 429 for Discord TokenSource, skipping")
						return nil
					}
					logger.Warn().Err(err).Msg("error getting Discord token for user, will revoke all roles")
					userTokensMutex.Lock()
					userTokensInvalid = append(userTokensInvalid, userID)
					userTokensMutex.Unlock()
					return nil
				}
				logger.Debug().Msg("refreshing UserGuildEdges")
				added, removed, err := RefreshUserGuildEdges(ctx, db, token, userID)
				if err != nil {
					var restErr *discordgo.RESTError
					if errors.As(err, &restErr) {
						if restErr.Response.StatusCode >= 500 {
							logger.Warn().Err(err).Msg("Discord API server error, skipping user")
						} else {
							logger.Warn().Err(err).Msg("error using Discord token for user, will revoke all roles")
							userTokensMutex.Lock()
							userTokensInvalid = append(userTokensInvalid, userID)
							userTokensMutex.Unlock()
						}
						return nil
					}
					if errors.Is(err, context.DeadlineExceeded) || strings.Contains(err.Error(), "Client.Timeout exceeded while awaiting headers") {
						logger.Warn().Err(err).Msg("skipping refresh for user due to API request timeout")
						return nil
					}
					logger.Err(err).Msg("error refreshing guilds for user")
					return err
				}
				if len(added)+len(removed) > 0 {
					logger.Info().
						Strs("addedGuildIDs", uints64ToStrs(added)).
						Strs("removedGuildIDs", uints64ToStrs(removed)).
						Msg("refreshed with changes")
				} else {
					logger.Debug().Int64("progress", processedCounter.Load()+1).Msg("refreshed with no changes")
				}
				userTokensMutex.Lock()
				userTokensValid = append(userTokensValid, userID)
				userTokensMutex.Unlock()
				return nil
			})
		}
		if err = eGroup.Wait(); err != nil {
			return userTokensInvalid, userTokensValid, totalCount, err
		}
		if len(userIDs) == 0 {
			log.Debug().Msg("no users to refresh, breaking out early")
			break
		}
		after = userIDs[len(userIDs)-1]
		log.Info().
			Int64("count", processedCounter.Load()).
			Int("total", totalCount).
			Msg("refresh progress")
		if len(userIDs) < pageSize {
			break
		}
		if max > 0 && int(processedCounter.Load()) >= max {
			break
		}
	}
	if len(userTokensInvalid) > 0 {
		log.Info().Int("count", len(userTokensInvalid)).
			Msg("failed to refresh some Discord tokens")
	}
	return userTokensInvalid, userTokensValid, totalCount, nil
}

// Refreshes guilds for all registered users.
func RefreshUserGuildEdges(
	ctx context.Context,
	db *ent.Client,
	token *oauth2.Token,
	userID uint64,
) (added []uint64, removed []uint64, err error) {
	svc, err := discordgo.New(fmt.Sprintf("Bearer %s", token.AccessToken))
	if err != nil {
		err = fmt.Errorf("error creating discordgo.Session: %w", err)
		return
	}
	userGuilds, err := svc.UserGuilds(0, "", "", false)
	if err != nil {
		err = fmt.Errorf("error getting UserGuilds: %w", err)
		return
	}
	guildIDs := make([]uint64, len(userGuilds))
	for i, dg := range userGuilds {
		guildID, convErr := strconv.ParseUint(dg.ID, 10, 64)
		if convErr != nil {
			err = convErr
			return
		}
		guildIDs[i] = guildID
	}
	// add guilds
	addGuildIDs, err := db.Guild.Query().
		Where(
			guild.IDIn(guildIDs...),
			guild.Not(
				guild.HasMembersWith(user.ID(userID)),
			),
		).
		IDs(ctx)
	if err != nil {
		err = fmt.Errorf("error getting Guilds to add: %w", err)
		return
	}
	// remove guilds
	removePredicates := []predicate.Guild{
		guild.HasMembersWith(user.ID(userID)),
	}
	if len(guildIDs) > 0 {
		removePredicates = append(removePredicates, guild.IDNotIn(guildIDs...))
	}
	removeGuildIDs, err := db.Guild.Query().
		Where(removePredicates...).
		IDs(ctx)
	if err != nil {
		err = fmt.Errorf("error getting Guilds to remove: %w", err)
		return
	}
	// actually do it
	err = db.User.UpdateOneID(userID).
		AddGuildIDs(addGuildIDs...).
		RemoveGuildIDs(removeGuildIDs...).
		Exec(ctx)
	return addGuildIDs, removeGuildIDs, err
}

func uints64ToStrs(input []uint64) []string {
	output := make([]string, len(input))
	for i, n := range input {
		output[i] = strconv.FormatUint(n, 10)
	}
	return output
}
