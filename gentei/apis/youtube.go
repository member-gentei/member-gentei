package apis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"libs.altipla.consulting/tokensource"
)

var (
	ErrNoYouTubeTokenForUser = errors.New("user does not have a YouTube token")
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
	logger := log.With().Uint64("userID", userID).Logger()
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
