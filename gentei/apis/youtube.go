package apis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	ErrInvalidGrant          = errors.New("invalid_grant - token expired or revoked")
	ErrYouTubeSignupRequired = errors.New("youtubeSignupRequired - Google account does not have YouTube account")
	ErrNoMembersOnlyVideos   = errors.New("YouTube channel has membership enabled, but no members-only videos")
)

func GetYouTubeService(ctx context.Context, db *ent.Client, userID uint64, config *oauth2.Config) (*youtube.Service, error) {
	notify, err := refreshingTokenSourceNotify(ctx, db, userID, config)
	if err != nil {
		return nil, err
	}
	client := retryablehttp.NewClient()
	client.HTTPClient = notify.Client(ctx)
	client.Logger = &zlgLeveledLoggerWrapper{logger: log.Logger}
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
				"Account has been deleted",
				"Token has been expired or revoked.",
				"Bad Request":
				if r != nil && r.Body != nil {
					r.Body.Close()
				}
				if errResponse.Error == "invalid_grant" {
					log.Warn().Err(err).Msg("invalid_grant error detected, returning more general case")
					return false, ErrInvalidGrant
				}
				return false, err
			}
		}
		if gErr, ok := err.(*googleapi.Error); ok {
			if gErr.Code == http.StatusBadRequest {
				// "While this can be a transient error" with YouTube is likely always a transient error
				if strings.HasSuffix(gErr.Message, "processingFailure") {
					return true, nil
				}
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
	var (
		membersOnlyVideoID    string
		membersOnlyPlaylistID = fmt.Sprintf("UUMO%s", channelID[2:])
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
				if videoID == "" {
					// youtube api sometimes be like this
					continue
				}
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
							vidLogger.Info().Err(ctlErr).Msg("comments disabled on video")
							continue
						}
						if gErr.Code == 403 {
							// this is fine! we just don't have permissions on this video.
							membersOnlyVideoID = videoID
							logger.Info().Str("videoID", videoID).Msg("selected members-only video ID")
							return errStopPagination
						}
					}
					vidLogger.Err(ctlErr).Msg("error checking members-only video validity")
					return ctlErr
				}
				membersOnlyVideoID = videoID
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
	return membersOnlyVideoID, err
}

func IsCommentsDisabledErr(err *googleapi.Error) bool {
	return GoogleErrHasReason(err, 403, "commentsDisabled")
}

func IsYouTubeSignupRequiredErr(err *googleapi.Error) bool {
	return GoogleErrHasReason(err, 401, "youtubeSignupRequired")
}

func GoogleErrHasReason(err *googleapi.Error, code int, reason string) bool {
	if err.Code != code {
		return false
	}
	for _, item := range err.Errors {
		if item.Reason == reason {
			return true
		}
	}
	return false
}

func IsUnusableYouTubeTokenErr(err error) bool {
	return errors.Is(err, ErrInvalidGrant) || errors.Is(err, ErrYouTubeSignupRequired)
}

type zlgLeveledLoggerWrapper struct {
	logger zerolog.Logger

	retryablehttp.LeveledLogger
}

func (l *zlgLeveledLoggerWrapper) logEvent(event *zerolog.Event, msg string, keysAndValues ...interface{}) {
	var (
		key string
	)
	for i := range keysAndValues {
		if i%2 == 0 {
			key = keysAndValues[i].(string)
			continue
		}
		switch v := keysAndValues[i].(type) {
		case func() (io.ReadCloser, error):
			// keys to skip
			continue
		case io.ReadCloser:
			content, err := ioutil.ReadAll(v)
			if err != nil {
				log.Warn().Str("key", key).Msg("zlgLeveledLoggerWrapper could not read io.ReadCloser")
				continue
			}
			v.Close()
			event = event.Str(key, string(content))
		default:
			event = event.Interface(key, keysAndValues[i])
		}
	}
	event.Bool("retryableHttp", true).Msg(msg)
}

func (l *zlgLeveledLoggerWrapper) Error(msg string, keysAndValues ...interface{}) {
	l.logEvent(l.logger.Error(), msg, keysAndValues...)
}

func (l *zlgLeveledLoggerWrapper) Info(msg string, keysAndValues ...interface{}) {
	l.logEvent(l.logger.Info(), msg, keysAndValues...)
}

func (l *zlgLeveledLoggerWrapper) Debug(msg string, keysAndValues ...interface{}) {
	l.logEvent(l.logger.Debug(), msg, keysAndValues...)
}

func (l *zlgLeveledLoggerWrapper) Warn(msg string, keysAndValues ...interface{}) {
	l.logEvent(l.logger.Warn(), msg, keysAndValues...)
}
