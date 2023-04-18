package membership

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/member-gentei/member-gentei/gentei/apis"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/member-gentei/member-gentei/gentei/ent/guildrole"
	"github.com/member-gentei/member-gentei/gentei/ent/predicate"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/member-gentei/member-gentei/gentei/ent/usermembership"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/youtube/v3"
)

type CheckForUserOptions struct {
	// Specify ChannelIDs to restrict checks to these channels.
	ChannelIDs []string
	// CheckDisabledChannels forces a check on channels that have been disabled *if ChannelIDs is not empty*
	CheckDisabledChannels bool
}

type CheckResultSet struct {
	// Gained contains memberships newly gained.
	Gained []CheckResult `json:",omitempty"`
	// Retained contains memberships that have been kept and re-validated.
	Retained []CheckResult `json:",omitempty"`
	// Lost contains memberships newly lost.
	Lost []CheckResult `json:",omitempty"`
	// Not contains memberships that this user does not have but did not newly lose.
	Not []CheckResult `json:",omitempty"`

	// Disabled contains channel IDs that were skipped due to disabled membership checks.
	DisabledChannels []string

	// YouTubeTokenInvalid notes that check failures are due to token expiry or revocation.
	YouTubeTokenInvalid bool
}

func (c *CheckResultSet) HasResults() bool {
	return len(c.Gained)+len(c.Retained)+len(c.Lost)+len(c.Not) > 0
}

func (c CheckResultSet) IsMember(channelID string) bool {
	for _, gained := range c.Gained {
		if channelID == gained.ChannelID {
			return true
		}
	}
	for _, retained := range c.Retained {
		if channelID == retained.ChannelID {
			return true
		}
	}
	return false
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
		predicates := []predicate.YouTubeTalent{
			youtubetalent.IDIn(options.ChannelIDs...),
		}
		talents, err = db.YouTubeTalent.Query().
			Where(predicates...).
			All(ctx)
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
		disabledChannelIDs           []string
		tokenInvalid                 bool
	)
	for _, talent := range talents {
		logger := log.With().
			Str("userID", strconv.FormatUint(userID, 10)).
			Str("channelID", talent.ID).
			Logger()
		if !talent.Disabled.IsZero() {
			if !options.CheckDisabledChannels {
				logger.Info().Msg("membership checks for channel disabled, skipping")
				disabledChannelIDs = append(disabledChannelIDs, talent.ID)
				continue
			}
			logger.Info().Msg("checking membership on disabled channel")
		}
		logger.Debug().Msg("checking membership for channel")
		// if we don't have a membership video ID for this talent, try to get one using a possibly eligible user every ~24 hours
		if talent.MembershipVideoID == "" {
			if time.Since(talent.LastMembershipVideoIDMiss).Hours() < 24 {
				continue
			}
			logger.Debug().Msg("fetching random members-only video ID")
			var videoID string
			videoID, err = apis.SelectRandomMembersOnlyVideoID(ctx, logger, svc, talent.ID)
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
		var isMember bool
		if !tokenInvalid {
			isMember, err = checkSingleMembership(ctx, logger, db, svc, talent.ID, talent.MembershipVideoID)
			if apis.IsUnusableYouTubeTokenErr(err) {
				logger.Warn().Err(err).Msg("YouTube token invalid for user, all memberships will be lost")
				tokenInvalid = true
			} else if err != nil {
				return nil, fmt.Errorf("error checking membership: %w", err)
			}
		}
		checkTimestamps[talent.ID] = time.Now()
		if isMember {
			verifiedMembershipChannelIDs = append(verifiedMembershipChannelIDs, talent.ID)
		} else {
			nonMemberChannelIDs = append(nonMemberChannelIDs, talent.ID)
		}
	}
	// merge in results
	results = &CheckResultSet{
		DisabledChannels:    disabledChannelIDs,
		YouTubeTokenInvalid: tokenInvalid,
	}
	// 1. get changes
	// 1a. get lost
	wasLost := map[string]bool{}
	lostIDs, err := db.YouTubeTalent.Query().Where(
		// "you're not a member but the bot enforced a role with you in it at some point"
		youtubetalent.IDIn(nonMemberChannelIDs...),
		youtubetalent.HasMembershipsWith(
			usermembership.HasUserWith(user.ID(userID)),
		),
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
		if wasLost[cid] {
			continue
		}
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

// SaveMemberships maintains UserMembership objects and YouTube association info, but not its GuildRole edges. It also sets User.LastCheck to the current time after effecting changes.
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
				Str("userID", strconv.FormatUint(userID, 10)).
				Str("talentID", c.ChannelID).
				Logger()
		)
		count, err = db.UserMembership.Update().
			Where(
				usermembership.HasYoutubeTalentWith(
					youtubetalent.ID(c.ChannelID),
					youtubetalent.DisabledIsNil(),
				),
				usermembership.HasUserWith(user.ID(userID)),
			).
			SetLastVerified(c.Time).
			SetFailCount(0).
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
				Str("userID", strconv.FormatUint(userID, 10)).
				Str("talentID", c.ChannelID).
				Logger()
		)
		count, err = db.UserMembership.Update().
			Where(
				usermembership.HasYoutubeTalentWith(
					youtubetalent.ID(c.ChannelID),
					youtubetalent.DisabledIsNil(),
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
		if count > 0 {
			logger.Info().
				Time("firstFailed", c.Time).
				Msg("lost membership")
		} else {
			logger.Debug().Msg("not a member; was 'lost'")
		}
	}
	// clear not
	for _, c := range results.Not {
		var (
			count  int
			logger = log.With().
				Str("userID", strconv.FormatUint(userID, 10)).
				Str("talentID", c.ChannelID).
				Logger()
		)
		count, err = db.UserMembership.Update().
			Where(
				usermembership.HasYoutubeTalentWith(
					youtubetalent.ID(c.ChannelID),
					youtubetalent.DisabledIsNil(),
				),
				usermembership.HasUserWith(user.ID(userID)),
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
	updateOne := db.User.UpdateOneID(userID).
		SetLastCheck(time.Now())
	if results.YouTubeTokenInvalid {
		updateOne = updateOne.
			ClearYoutubeID().
			ClearYoutubeToken()
	}
	err = updateOne.Exec(ctx)
	return
}

// checkSingleMembership performs a membership check and handles membership video reassignment
func checkSingleMembership(
	ctx context.Context,
	logger zerolog.Logger,
	db *ent.Client,
	svc *youtube.Service,
	channelID string,
	membershipVideoID string,
) (isMember bool, err error) {
	_, err = svc.
		CommentThreads.List([]string{"id"}).
		VideoId(membershipVideoID).Do()
	logger = logger.With().Str("videoID", membershipVideoID).Logger()
	logger.Info().Msg("CommentThreads.List")
	if err != nil {
		var gErr *googleapi.Error
		if errors.As(err, &gErr) {
			logger.Debug().Interface("gErr", gErr).Msg("membership check encountered googleapi.Error")
			if apis.IsCommentsDisabledErr(gErr) || gErr.Code == 404 {
				// if comments are disabled on this video, we need to select a new video.
				logger.Warn().Err(err).
					Msg("missing video or comments disabled on membership check video, getting a new one")
				newVideoID, selectErr := apis.SelectRandomMembersOnlyVideoID(ctx, logger, svc, channelID)
				if selectErr != nil || newVideoID == "" {
					logger.Err(err).Msg("error getting new membership check video, disabling checks for channel")
					disableErr := db.YouTubeTalent.UpdateOneID(channelID).
						SetDisabled(time.Now()).
						Exec(ctx)
					if disableErr != nil {
						logger.Err(err).Msg("error disabling checks on channel")
					}
					if selectErr != nil {
						err = selectErr
					} else {
						err = apis.ErrNoMembersOnlyVideos
					}
					return
				}
				// do it all over again!
				err = db.YouTubeTalent.UpdateOneID(channelID).
					SetMembershipVideoID(newVideoID).
					Exec(ctx)
				if err != nil {
					err = fmt.Errorf("error setting new membership check video: %w", err)
					return
				}
				logger.Info().Str("videoID", newVideoID).Msg("checking with new video")
				return checkSingleMembership(ctx, logger, db, svc, channelID, membershipVideoID)
			}
			if apis.IsYouTubeSignupRequiredErr(gErr) {
				return false, apis.ErrYouTubeSignupRequired
			}
			if gErr.Code == 403 {
				// not a member
				return false, nil
			}
		}
		// the caller should decide what to do on errors
		return
	}
	isMember = true
	return
}

func createMissingUserMemberships(ctx context.Context, db *ent.Client, userID uint64, c CheckResult) error {
	var (
		logger = log.With().
			Str("userID", strconv.FormatUint(userID, 10)).
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
	if len(roleIDsToGrant) == 0 {
		return nil
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
