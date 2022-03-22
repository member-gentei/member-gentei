package apis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	irand "math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"libs.altipla.consulting/tokensource"
)

var (
	ErrNoYouTubeTokenForUser = errors.New("user does not have a YouTube token")
	ErrNoMembersOnlyVideos   = errors.New("YouTube channel has membership enabled, but no members-only videos")
)

func GetYouTubeService(ctx context.Context, db *ent.Client, userID uint64, config *oauth2.Config) (*youtube.Service, error) {
	notify, err := refreshingTokenSourceNotify(ctx, db, userID, config)
	if err != nil {
		return nil, err
	}
	client := retryablehttp.NewClient()
	client.HTTPClient = notify.Client(ctx)
	client.CheckRetry = youTubeAPIRetryPolicy
	return youtube.NewService(ctx, option.WithHTTPClient(client.StandardClient()))
}

// youTubeAPIRetryPolicy is either some known 400 errors or retryablehttp.DefaultRetryPolicy.
func youTubeAPIRetryPolicy(ctx context.Context, r *http.Response, err error) (bool, error) {
	if err != nil {
		if rErr, ok := scavengeRetrieveError(err); ok {
			var errResponse struct {
				Error            string
				ErrorDescription string `json:"error_description"`
			}
			json.Unmarshal(rErr.Body, &errResponse)
			switch errResponse.ErrorDescription {
			case
				"Token has been expired or revoked.",
				"Bad Request":
				r.Body.Close()
				return false, err
			}
		}
	}
	return retryablehttp.DefaultRetryPolicy(ctx, r, err)
}

// refreshingTokenSourceNotify keeps tokens fresh in the ent database.
func refreshingTokenSourceNotify(ctx context.Context, db *ent.Client, userID uint64, config *oauth2.Config) (*tokensource.Notify, error) {
	u, err := db.User.Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error performing initial user load: %w", err)
	}
	if u.YoutubeToken == nil {
		return nil, ErrNoYouTubeTokenForUser
	}
	logger := log.With().Str("userID", strconv.FormatUint(userID, 10)).Logger()
	notify := tokensource.NewNotifyHook(ctx, config, u.YoutubeToken, func(token *oauth2.Token) error {
		logger.Debug().Msg("YouTube token for user refreshed")
		return db.User.UpdateOneID(userID).
			SetYoutubeToken(token).Exec(ctx)
	})
	return notify, nil
}

func scavengeRetrieveError(err error) (*oauth2.RetrieveError, bool) {
	if rErr, ok := err.(*oauth2.RetrieveError); ok {
		return rErr, ok
	}
	errString := err.Error()
	log.Debug().Str("errString", errString).Msg("oauth2.RetrieveError?")
	if strings.Contains(errString, "oauth2: cannot fetch token: ") {
		rIdx := strings.Index(errString, "\nResponse: ")
		stringBody := errString[rIdx+len("\nResponse: "):]
		return &oauth2.RetrieveError{
			Body: []byte(stringBody),
		}, true
	}
	return nil, false
}

var (
	errStopPagination = errors.New("pls stop")
)

// SelectRandomMembersOnlyVideoID chooses a random members-only video that has comments enabled.
func SelectRandomMembersOnlyVideoID(
	ctx context.Context,
	logger zerolog.Logger,
	svc *youtube.Service,
	channelID string,
) (string, error) {
	membersOnlyPlaylistID := fmt.Sprintf("UUMO%s", channelID[2:])
	var (
		membersOnlyVideoID string
	)
	err := svc.PlaylistItems.List([]string{"snippet"}).
		PlaylistId(membersOnlyPlaylistID).
		MaxResults(50).
		Pages(ctx, func(pilr *youtube.PlaylistItemListResponse) error {
			if len(pilr.Items) == 0 {
				return ErrNoMembersOnlyVideos
			}
			// shuffle the current page
			irand.Shuffle(len(pilr.Items), func(i, j int) {
				pilr.Items[i], pilr.Items[j] = pilr.Items[j], pilr.Items[i]
			})
			for _, item := range pilr.Items {
				videoID := item.Snippet.ResourceId.VideoId
				// perform membership check
				_, ctlErr := svc.CommentThreads.
					List([]string{"id"}).
					VideoId(videoID).Do()
				vidLogger := logger.With().Str("videoID", videoID).Bool("selectVideoID", true).Logger()
				vidLogger.Info().Msg("CommentThreads.List")
				if ctlErr != nil {
					var gErr *googleapi.Error
					if errors.As(ctlErr, &gErr) {
						if IsCommentsDisabledErr(gErr) {
							vidLogger.Info().Msg("comments disabled on video")
							continue
						}
						if gErr.Code == 403 {
							// this is fine! we just don't have permissions on this video.
							membersOnlyVideoID = videoID
							return errStopPagination
						}
					}
					vidLogger.Err(ctlErr).Msg("error checking members-only video validity")
					return ctlErr
				}

			}
			return nil
		})
	if errors.Is(err, errStopPagination) {
		if membersOnlyVideoID == "" {
			err = ErrNoMembersOnlyVideos
		} else {
			err = nil
		}
	} else if err != nil {
		var gErr *googleapi.Error
		if errors.As(err, &gErr) && gErr.Code == 404 {
			return "", ErrNoMembersOnlyVideos
		}
		return "", err
	}
	return membersOnlyVideoID, nil
}

func IsCommentsDisabledErr(err *googleapi.Error) bool {
	if err.Code == 403 {
		for _, item := range err.Errors {
			if item.Reason == "commentsDisabled" {
				return true
			}
		}
	}
	return false
}
