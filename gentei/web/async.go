package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/member-gentei/member-gentei/gentei/apis"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
	"github.com/rs/zerolog/log"
)

type ErrNoMembershipPlaylist struct {
	ChannelID string
	error
}

func (e ErrNoMembershipPlaylist) Error() string {
	return fmt.Sprintf("channel does not have a membership playlist: %s", e.ChannelID)
}

// UpsertYouTubeChannelID upserts YouTube channels by ID.
//
// Returns an error if the channel does not have a membership playlist.
func UpsertYouTubeChannelID(ctx context.Context, db *ent.Client, channelID string) error {
	logger := log.With().Str("ytChannelID", channelID).Logger()
	// check that the RSS feed for members-only videos exists
	// TODO: this URL will probably break at any moment because it is useless for actual legitimate membership video updates
	membershipPlaylist := fmt.Sprintf("https://www.youtube.com/feeds/videos.xml?playlist_id=UUMO%s", channelID[2:])
	pls, err := http.Head(membershipPlaylist)
	if err != nil {
		logger.Err(err).Msg("error querying for membership playlist")
		return err
	}
	if pls.StatusCode == http.StatusNotFound {
		return ErrNoMembershipPlaylist{ChannelID: channelID}
	}
	logger.Debug().Msg("attempting to get OpenGraph data for channel")
	ogp, err := apis.GetYouTubeChannelOG(channelID)
	if err != nil {
		logger.Err(err).Msg("error getting OpenGraph data for YouTube channel")
		return err
	}
	var thumbnailURL string
	for _, imageData := range ogp.Image {
		if imageData.Height == 900 {
			thumbnailURL = imageData.URL
			break
		}
		thumbnailURL = imageData.URL
	}
	err = db.YouTubeTalent.Create().
		SetID(channelID).
		SetChannelName(ogp.Title).
		SetThumbnailURL(thumbnailURL).
		SetLastUpdated(time.Now()).
		OnConflictColumns(youtubetalent.FieldID).
		UpdateNewValues().
		Exec(ctx)
	if err != nil {
		logger.Err(err).Msg("error saving channel")
		return err
	}
	logger.Info().Msg("saved OpenGraph data for YouTube channel")
	return nil
}
