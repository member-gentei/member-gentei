package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

var apiKeys = map[string]map[string]bool{}

// Logical errors
const (
	apiKeysEnvName    = "API_CONFIG"
	tokenInvalidError = "API token invalid"
	// ErrNoChannelPermission = errors.New("no permissions to view channel")
)

// NewAuthHandler creates the API authentication middleware.
func NewAuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeaders := r.Header["Authorization"]
		if len(authHeaders) != 1 {
			fmt.Fprint(w, tokenInvalidError)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		authHeader := strings.TrimSpace(authHeaders[0])
		if len(authHeader) != len("Bearer 3fdab78b-6b89-4e66-ac8b-ec89f556f27b") {
			fmt.Fprint(w, tokenInvalidError)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if _, exists := apiKeys[mustGetToken(r)]; !exists {
			fmt.Fprint(w, tokenInvalidError)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func mustGetToken(r *http.Request) string {
	return strings.TrimSpace(r.Header["Authorization"][0])[len("Bearer "):]
}

func keyHasPermission(r *http.Request, channelSlug string) bool {
	slugs, exists := apiKeys[mustGetToken(r)]
	if !exists {
		return false
	}
	_, exists = slugs[channelSlug]
	if !exists {
		_, exists = slugs["all"]
	}
	return exists
}

func init() {
	data, err := base64.StdEncoding.DecodeString(mustLoadEnv(apiKeysEnvName))
	if err != nil {
		log.Warn().Err(err).Msgf("error reading '%s', all requests will be denied", apiKeysEnvName)
		return
	}
	err = json.Unmarshal(data, &apiKeys)
	if err != nil {
		log.Warn().Err(err).Msgf("error decoding '%s', all requests will be denied", apiKeysEnvName)
		return
	}
	log.Info().Interface("apiKeys", apiKeys).Msg("loaded keys")
}
