package membership

import (
	"context"
	"fmt"
	"time"

	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"golang.org/x/oauth2"
)

type CheckStaleOptions struct {
	// StaleThreshold is used in a <= comparison to the last stored membership check time.
	StaleThreshold time.Duration
	// MembershipChangeHook gets called whnen a user experiences a change in channel membership.
	MembershipChangeHook func(userID uint64, results *CheckResultSet) error
}

var DefaultCheckStaleOptions = &CheckStaleOptions{
	StaleThreshold: time.Hour * 12,
}

func CheckStale(ctx context.Context, db *ent.Client, youtubeConfig *oauth2.Config, options *CheckStaleOptions) error {
	if options == nil {
		options = DefaultCheckStaleOptions
	}
	staleThreshold := options.StaleThreshold
	if options.StaleThreshold > 0 {
		staleThreshold *= -1
	}
	for {
		staleUserIDs, err := db.User.Query().
			Where(
				user.HasGuildsWith(
					guild.HasYoutubeTalents(),
				),
				user.LastCheckLTE(time.Now().Add(staleThreshold)),
			).
			Limit(1000).
			IDs(ctx)
		if err != nil {
			return err
		}
		if len(staleUserIDs) == 0 {
			break
		}
		for _, userID := range staleUserIDs {
			// TODO: https://github.com/member-gentei/member-gentei/issues/92
			results, err := CheckForUser(ctx, db, youtubeConfig, userID, nil)
			if err != nil {
				return fmt.Errorf("error checking memberships for user '%d': %w", userID, err)
			}
			if options.MembershipChangeHook != nil && (len(results.Lost) > 0 || len(results.Gained) > 0) {
				err = options.MembershipChangeHook(userID, results)
				if err != nil {
					return fmt.Errorf("error calling MembershipChangeHook for user '%d': %w", userID, err)
				}
			}
			err = db.User.UpdateOneID(userID).
				SetLastCheck(time.Now()).
				Exec(ctx)
			if err != nil {
				return fmt.Errorf("error saving LastCheck for user '%d': %w", userID, err)
			}
		}
	}
	return nil
}
