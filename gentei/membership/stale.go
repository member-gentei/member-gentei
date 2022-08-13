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
}

var DefaultCheckStaleOptions = &CheckStaleOptions{
	StaleThreshold: time.Hour * 12,
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
	)
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
		for _, userID := range staleUserIDs {
			// TODO: https://github.com/member-gentei/member-gentei/issues/92
			results, err := CheckForUser(ctx, db, youtubeConfig, userID, nil)
			if err != nil {
				return fmt.Errorf("error checking memberships for user '%d': %w", userID, err)
			}
			if options.NoSave {
				continue
			}
			err = SaveMemberships(ctx, db, userID, results)
			if err != nil {
				return fmt.Errorf("error saving memberships for user '%d': %w", userID, err)
			}
			log.Info().Int("count", totalStaleCount).Msg("refreshed stale user batch of <=1000")
		}
	}
	log.Info().Int("count", totalStaleCount).Msg("refreshed stale users")
	return nil
}

// RefreshAllUserGuildEdges refreshes guild edges for all registered users. Returns a slice of userIDs that could not be refreshed and a count of all users.
func RefreshAllUserGuildEdges(ctx context.Context, db *ent.Client, discordConfig *oauth2.Config) ([]uint64, int, error) {
	// refresh everyone's tokens
	return refreshUserGuildEdgesWithPredicates(ctx, db, discordConfig, nil)
}

// RefreshStaleUserGuildEdges
func RefreshStaleUserGuildEdges(ctx context.Context, db *ent.Client, discordConfig *oauth2.Config, staleThreshold time.Duration) ([]uint64, int, error) {
	staleBefore := time.Now().Add(-staleThreshold)
	return refreshUserGuildEdgesWithPredicates(ctx, db, discordConfig, user.LastCheckLT(staleBefore))
}

// refreshUserGuildEdgesWithPredicates is the inner implementation of all stale refreshes.
func refreshUserGuildEdgesWithPredicates(ctx context.Context, db *ent.Client, discordConfig *oauth2.Config, predicates ...predicate.User) ([]uint64, int, error) {
	var (
		userTokensInvalid      []uint64
		totalCount             int
		after                  uint64
		userTokensInvalidMutex = &sync.Mutex{}
	)
	const pageSize = 400
	for {
		processedCounter := &atomic.Int64{}
		userIDs, err := db.User.Query().
			Where(append(
				predicates,
				user.IDGT(after),
			)...).
			Order(ent.Asc(user.FieldID)).
			Limit(pageSize).
			IDs(ctx)
		if err != nil {
			return nil, 0, fmt.Errorf("error paginating user IDs: %w", err)
		}
		var eGroup errgroup.Group
		eGroup.SetLimit(10)
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
					logger.Warn().Err(err).Msg("error getting Discord token for user, will revoke all roles")
					userTokensInvalidMutex.Lock()
					userTokensInvalid = append(userTokensInvalid, userID)
					userTokensInvalidMutex.Unlock()
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
							userTokensInvalidMutex.Lock()
							userTokensInvalid = append(userTokensInvalid, userID)
							userTokensInvalidMutex.Unlock()
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
					logger.Debug().Msg("refreshed with no changes")
				}
				return nil
			})
		}
		if err = eGroup.Wait(); err != nil {
			return userTokensInvalid, totalCount, err
		}
		totalCount += len(userIDs)
		if len(userIDs) < pageSize {
			break
		}
		after = userIDs[len(userIDs)-1]
	}
	if len(userTokensInvalid) > 0 {
		log.Info().Int("count", len(userTokensInvalid)).
			Msg("failed to refresh some Discord tokens")
	}
	return userTokensInvalid, totalCount, nil
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
	userGuilds, err := svc.UserGuilds(0, "", "")
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
