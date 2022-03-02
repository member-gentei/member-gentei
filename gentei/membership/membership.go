package membership

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/member-gentei/member-gentei/gentei/apis"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/member-gentei/member-gentei/gentei/ent/guildrole"
	"github.com/member-gentei/member-gentei/gentei/ent/predicate"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/member-gentei/member-gentei/gentei/ent/usermembership"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi"
)

type CheckForUserOptions struct {
	// Specify ChannelIDs to restrict checks to these channels.
	ChannelIDs []string
}

type CheckResultSet struct {
	// Gained contains memberships newly gained.
	Gained []CheckResult
	// Retained contains memberships that have been kept and re-validated.
	Retained []CheckResult
	// Lost contains memberships newly lost.
	Lost []CheckResult
	// Not contains memberships that this user does not have but did not newly lose.
	Not []CheckResult
}

type CheckResult struct {
	ChannelID string
	Time      time.Time
}

func CheckForUser(
	ctx context.Context, db *ent.Client,
	youtubeConfig *oauth2.Config,
	userID uint64,
	options *CheckForUserOptions,
) (results *CheckResultSet, err error) {
	if options == nil {
		options = &CheckForUserOptions{}
	}
	svc, err := apis.GetYouTubeService(ctx, db, userID, youtubeConfig)
	if err != nil {
		return
	}
	// load talents
	var (
		talents []*ent.YouTubeTalent
	)
	if options.ChannelIDs != nil {
		talents, err = db.YouTubeTalent.Query().
			Where(
				youtubetalent.IDIn(options.ChannelIDs...),
				youtubetalent.Or(
					youtubetalent.Disabled(time.Time{}),
					youtubetalent.DisabledIsNil(),
				),
			).All(ctx)
		if err != nil {
			err = fmt.Errorf("error getting specified channels: %w", err)
			return
		}
	} else {
		talents, err = db.YouTubeTalent.Query().
			Where(
				youtubetalent.HasGuildsWith(
					guild.Or(
						guild.HasMembersWith(user.ID(userID)),
						guild.HasAdminsWith(user.ID(userID)),
					),
				),
			).All(ctx)
		if err != nil {
			err = fmt.Errorf("error querying for eligible channels: %w", err)
			return
		}
	}
	var (
		checkTimestamps              = map[string]time.Time{}
		verifiedMembershipChannelIDs []string
		nonMemberChannelIDs          []string
	)
	for _, talent := range talents {
		logger := log.With().
			Str("userID", strconv.FormatUint(userID, 10)).
			Str("channelID", talent.ID).
			Logger()
		logger.Debug().Msg("checking membership for channel")
		// if we don't have a membership video ID for this talent, try to get one using a possibly eligible user every ~24 hours
		if talent.MembershipVideoID == "" {
			if time.Since(talent.LastMembershipVideoIDMiss).Hours() < 24 {
				continue
			}
			logger.Debug().Msg("fetching random members-only video ID")
			var videoID string
			videoID, err = apis.SelectRandomMembersOnlyVideoID(ctx, svc, talent.ID)
			if errors.Is(err, apis.ErrNoMembersOnlyVideos) {
				err = db.YouTubeTalent.UpdateOne(talent).
					SetLastMembershipVideoIDMiss(time.Now()).
					Exec(ctx)
				if err != nil {
					return
				}
				continue
			} else if err != nil {
				return
			}
			talent.MembershipVideoID = videoID
			err = db.YouTubeTalent.UpdateOne(talent).
				SetMembershipVideoID(videoID).
				Exec(ctx)
			if err != nil {
				return
			}
		}
		// perform membership check
		_, err = svc.CommentThreads.
			List([]string{"id"}).
			VideoId(talent.MembershipVideoID).Do()
		logger = logger.With().Str("videoID", talent.MembershipVideoID).Logger()
		logger.Info().Msg("CommentThreads.List")
		if err != nil {
			var gErr *googleapi.Error
			if errors.As(err, &gErr) {
				logger.Debug().Interface("gErr", gErr).Msg("membership check encountered googleapi.Error")
				if gErr.Code == 403 {
					if apis.IsCommentsDisabledErr(gErr) {
						// if comments are disabled on this video, we need to select a new video.
						err = gErr
						logger.Err(err).Msg("comments disabled on membership check video")
						return
					}
					// not a member
					checkTimestamps[talent.ID] = time.Now()
					nonMemberChannelIDs = append(nonMemberChannelIDs, talent.ID)
					continue
				}
			}
			if !strings.HasSuffix(err.Error(), "commentsDisabled") {
				logger.Err(err).Msg("actual error fetching comments for membership check video")
				return
			}
		}
		checkTimestamps[talent.ID] = time.Now()
		verifiedMembershipChannelIDs = append(verifiedMembershipChannelIDs, talent.ID)
	}
	// merge in results
	results = &CheckResultSet{}
	// 1. get changes
	// 1a. get lost
	wasLost := map[string]bool{}
	lostIDs, err := db.YouTubeTalent.Query().Where(
		youtubetalent.IDIn(nonMemberChannelIDs...),
		youtubetalent.Not(youtubetalent.HasMembershipsWith(
			usermembership.HasUserWith(user.ID(userID)),
		)),
	).IDs(ctx)
	if err != nil {
		err = fmt.Errorf("error fetching lost membership IDs: %w", err)
		return
	}
	for _, cid := range lostIDs {
		results.Lost = append(results.Lost, CheckResult{
			ChannelID: cid,
			Time:      checkTimestamps[cid],
		})
		wasLost[cid] = true
	}
	// 1b. get gained
	gainedIDs, err := db.YouTubeTalent.Query().Where(
		youtubetalent.IDIn(verifiedMembershipChannelIDs...),
		youtubetalent.Not(youtubetalent.HasMembershipsWith(
			usermembership.HasUserWith(user.ID(userID)),
		)),
	).IDs(ctx)
	if err != nil {
		err = fmt.Errorf("error fetching gained membership IDs: %w", err)
		return
	}
	for _, cid := range gainedIDs {
		results.Gained = append(results.Gained, CheckResult{
			ChannelID: cid,
			Time:      checkTimestamps[cid],
		})
	}
	// 2. next, get non-changes
	// 2a. get retained
	retainedIDs, err := db.YouTubeTalent.Query().Where(
		youtubetalent.IDIn(verifiedMembershipChannelIDs...),
		youtubetalent.HasMembershipsWith(
			usermembership.HasUserWith(user.ID(userID)),
		),
	).IDs(ctx)
	if err != nil {
		err = fmt.Errorf("error fetching retained membership IDs: %w", err)
		return
	}
	for _, cid := range retainedIDs {
		results.Retained = append(results.Retained, CheckResult{
			ChannelID: cid,
			Time:      checkTimestamps[cid],
		})
	}
	// 2b. get not (that were checked)
	for _, cid := range nonMemberChannelIDs {
		if wasLost[cid] {
			continue
		}
		results.Not = append(results.Not, CheckResult{
			ChannelID: cid,
			Time:      checkTimestamps[cid],
		})
	}
	return
}

// SaveMemberships maintains UserMembership objects, but not its GuildRole edges.
func SaveMemberships(
	ctx context.Context,
	db *ent.Client,
	userID uint64,
	results *CheckResultSet,
) (err error) {
	// upsert retained
	for _, c := range results.Retained {
		var (
			count  int
			logger = log.With().
				Uint64("userID", userID).
				Str("talentID", c.ChannelID).
				Logger()
		)
		count, err = db.UserMembership.Update().
			Where(
				usermembership.HasYoutubeTalentWith(
					youtubetalent.ID(c.ChannelID),
				),
				usermembership.HasUserWith(user.ID(userID)),
			).
			SetLastVerified(c.Time).
			Save(ctx)
		if err != nil {
			err = fmt.Errorf("error updating last verified time for %s: %w", c.ChannelID, err)
			return
		}
		if count > 0 {
			logger.Info().Int("count", count).Time("lastVerified", c.Time).
				Msg("updated retained membership")
		}
		err = createMissingUserMemberships(ctx, db, userID, c)
		if err != nil {
			return err
		}
	}
	// create gained
	for _, c := range results.Gained {
		err = createMissingUserMemberships(ctx, db, userID, c)
		if err != nil {
			return err
		}
	}
	// update lost
	for _, c := range results.Lost {
		var (
			count  int
			logger = log.With().
				Uint64("userID", userID).
				Str("talentID", c.ChannelID).
				Logger()
		)
		count, err = db.UserMembership.Update().
			Where(
				usermembership.HasYoutubeTalentWith(
					youtubetalent.ID(c.ChannelID),
				),
				usermembership.HasUserWith(user.ID(userID)),
			).
			SetFirstFailed(c.Time).
			AddFailCount(1).
			Save(ctx)
		if err != nil {
			err = fmt.Errorf("error setting first failed time for %s: %w", c.ChannelID, err)
			return
		}
		logger.Info().Int("count", count).Time("firstFailed", c.Time).
			Msg("lost membership")
	}
	// clear not
	for _, c := range results.Not {
		var (
			count  int
			logger = log.With().
				Uint64("userID", userID).
				Str("talentID", c.ChannelID).
				Logger()
		)
		count, err = db.UserMembership.Update().
			Where(
				usermembership.HasYoutubeTalentWith(
					youtubetalent.ID(c.ChannelID),
				),
			).
			AddFailCount(1).
			Save(ctx)
		if err != nil {
			err = fmt.Errorf("error incrementing check failures for %s: %w", c.ChannelID, err)
			return
		}
		logger.Debug().Int("count", count).Time("firstFailed", c.Time).
			Msg("incremented non-membership")
	}
	return
}

func createMissingUserMemberships(ctx context.Context, db *ent.Client, userID uint64, c CheckResult) error {
	var (
		logger = log.With().
			Uint64("userID", userID).
			Str("talentID", c.ChannelID).
			Logger()
		missingRolePredicates = []predicate.GuildRole{
			guildrole.HasGuildWith(
				guild.Or(
					guild.HasAdminsWith(user.ID(userID)),
					guild.HasMembersWith(user.ID(userID)),
				),
			),
			guildrole.HasTalentWith(youtubetalent.ID(c.ChannelID)),
			// "not created yet"
			guildrole.Not(
				guildrole.HasUserMembershipsWith(usermembership.HasUserWith(user.ID(userID))),
			),
		}
		roleIDsToGrant []uint64
	)
	roleIDsToGrant, err := db.GuildRole.Query().
		Where(missingRolePredicates...).
		IDs(ctx)
	if err != nil {
		return fmt.Errorf("error querying for eligible roles to %s: %w", c.ChannelID, err)
	}
	logger.Info().Uints64("roleIDs", roleIDsToGrant).Msg("granting newly gained Discord roles to user")
	err = db.UserMembership.Create().
		SetFailCount(0).
		SetLastVerified(c.Time).
		SetUserID(userID).
		SetYoutubeTalentID(c.ChannelID).
		AddRoleIDs(roleIDsToGrant...).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error creating UserMembership for %s: %w", c.ChannelID, err)
	}
	return nil
}
