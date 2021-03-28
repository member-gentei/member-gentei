package clients

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

// YouTubeAPIRetryPolicy is either some known 400 errors or retryablehttp.DefaultRetryPolicy.
func YouTubeAPIRetryPolicy(ctx context.Context, r *http.Response, err error) (bool, error) {
	if err != nil && r != nil {
		if r.StatusCode == http.StatusBadRequest {
			if rErr, ok := scavengeRetrieveError(r, err); ok {
				var errResponse struct {
					Error            string
					ErrorDescription string `json:"error_description"`
				}
				json.Unmarshal(rErr.Body, &errResponse)
				switch errResponse.ErrorDescription {
				case
					`Token has been expired or revoked.`,
					`Bad Request`:
					return false, err
				}
			}
		}
	}
	return retryablehttp.DefaultRetryPolicy(ctx, r, err)
}

func scavengeRetrieveError(response *http.Response, err error) (*oauth2.RetrieveError, bool) {
	if rErr, ok := err.(*oauth2.RetrieveError); ok {
		return rErr, ok
	}
	errString := err.Error()
	log.Debug().Str("errString", errString).Msg("oauth2.RetrieveError?")
	if strings.Contains(errString, "oauth2: cannot fetch token: ") {
		rIdx := strings.Index(errString, "\nResponse: ")
		stringBody := errString[rIdx+len("\nResponse: "):]
		return &oauth2.RetrieveError{
			Response: response,
			Body:     []byte(stringBody),
		}, true
	}
	return nil, false
}
