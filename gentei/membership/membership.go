package membership

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/member-gentei/member-gentei/gentei/apis"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/member-gentei/member-gentei/gentei/ent/schema"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi"
)

type CheckForUserOptions struct {
	// Specify ChannelIDs to restrict checks to these channels.
	ChannelIDs []string
}

func CheckForUser(
	ctx context.Context, db *ent.Client,
	youtubeConfig *oauth2.Config,
	userID uint64,
	options *CheckForUserOptions,
) (lost []string, gained []string, retained []string, err error) {
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
			Where(youtubetalent.IDIn(options.ChannelIDs...)).All(ctx)
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
		checkTime             = time.Now()
		verifiedMemberships   = map[string]time.Time{}
		verifiedMembershipIDs []string
	)
	for _, talent := range talents {
		logger := log.With().Str("channelID", talent.ID).Logger()
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
		if err != nil {
			var gErr *googleapi.Error
			if errors.As(err, &gErr) {
				if gErr.Code == 403 {
					// not a member
					continue
				}
			}
			if !strings.HasSuffix(err.Error(), "commentsDisabled") {
				log.Err(err).Msg("actual error fetching comments for membership check video")
				return
			}
		}
		verifiedMemberships[talent.ID] = time.Now()
		verifiedMembershipIDs = append(verifiedMembershipIDs, talent.ID)
	}
	// merge in results
	u, err := db.User.Get(ctx, userID)
	if err != nil {
		return
	}
	// set Past on newly lost memberships
	for cid, meta := range u.MembershipMetadata {
		if verifiedMemberships[cid].IsZero() && !meta.Past {
			meta.Past = true
			u.MembershipMetadata[cid] = meta
			lost = append(lost, cid)
		}
	}
	// set newly gained memberships
	if u.MembershipMetadata == nil {
		u.MembershipMetadata = map[string]schema.MembershipMetadata{}
	}
	for cid, verifiedTime := range verifiedMemberships {
		if u.MembershipMetadata[cid].LastVerified.IsZero() {
			u.MembershipMetadata[cid] = schema.MembershipMetadata{
				LastVerified: verifiedTime,
			}
			gained = append(gained, cid)
		} else {
			meta := u.MembershipMetadata[cid]
			meta.LastVerified = verifiedTime
			u.MembershipMetadata[cid] = meta
			retained = append(retained, cid)
		}
	}
	err = db.User.UpdateOneID(userID).
		SetMembershipMetadata(u.MembershipMetadata).
		SetLastCheck(checkTime).
		ClearYoutubeMemberships().
		AddYoutubeMembershipIDs(verifiedMembershipIDs...).
		Exec(ctx)
	return
}
