package membership

import (
	"context"
	"fmt"
	"time"

	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/predicate"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

type PerformCheckOpts struct {
	BatchSize      int
	StaleThreshold time.Duration
	Enforce        bool
}

func DefaultPerformCheckOptions() *PerformCheckOpts {
	return &PerformCheckOpts{
		BatchSize:      100,
		StaleThreshold: time.Hour * 12,
		Enforce:        true,
	}
}

func PerformCheckBatches(
	ctx context.Context,
	db *ent.Client,
	discordConfig, ytConfig *oauth2.Config,
	userDeleteChan chan<- uint64,
	opts *PerformCheckOpts,
) error {
	var staleUserPredicates []predicate.User
	if opts.StaleThreshold != 0 {
		staleUserPredicates = []predicate.User{user.LastCheckLT(time.Now().Add(-opts.StaleThreshold))}
	}
	for {
		// refresh tokens
		failedUserIDs, validUserIDs, totalStaleCount, err := refreshUserGuildEdgesWithPredicates(
			ctx,
			db,
			discordConfig,
			opts.BatchSize,
			staleUserPredicates...,
		)
		if err != nil {
			return fmt.Errorf("error refreshing User Guild edges: %w", err)
		}
		for _, userID := range failedUserIDs {
			userDeleteChan <- userID
		}
		log.Info().
			Int("failed", len(failedUserIDs)).
			Int("succeeded", len(validUserIDs)).
			Msg("refreshed guild edges")
		if totalStaleCount == 0 {
			break
		}
		err = CheckStale(ctx, db, ytConfig, &CheckStaleOptions{
			StaleThreshold: opts.StaleThreshold,
			AdditionalUserPredicates: []predicate.User{
				user.IDIn(validUserIDs...),
			},
			MaxWorkers: 100,
		})
		if err != nil {
			return fmt.Errorf("error checking memberships: %w", err)
		}
		uncheckableCount, err := db.User.Update().
			Where(
				append(
					staleUserPredicates,
					user.IDIn(validUserIDs...),
				)...,
			).
			SetLastCheck(time.Now()).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("error updating LastCheck for uncheckable users: %w", err)
		}
		if uncheckableCount > 0 {
			log.Info().Int("count", uncheckableCount).Msg("updated LastCheck for uncheckable users")
		}
		if len(failedUserIDs)+len(validUserIDs) < opts.BatchSize {
			break
		}
	}
	return nil
}
