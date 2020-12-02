package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi"
	"github.com/member-gentei/member-gentei/pkg/common"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	youtubeClientIDEnvName     = "YOUTUBE_CLIENT_ID"
	youtubeClientSecretEnvName = "YOUTUBE_CLIENT_SECRET"
	youtubeRedirectURIEnvName  = "YOUTUBE_REDIRECT_URI"
)

var (
	youtubeOAuthConfig *oauth2.Config
	apiHandler         http.Handler
	fs                 *firestore.Client
	swagger            *openapi3.Swagger

	// to be written by tests
	ytClientOptions []option.ClientOption
)

type apiImpl struct {
	ServerInterface
}

type getMembersResponse struct {
	Users []ChannelMember `json:"users"`
	After string          `json:"after,omitempty"`
}

func (a *apiImpl) GetMembers(
	w http.ResponseWriter, r *http.Request,
	channelSlug ChannelSlugPathParam, params GetMembersParams,
) {
	if !keyHasPermission(r, string(channelSlug)) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	ctx := r.Context()
	membersRef := fs.Collection("channels").Doc(string(channelSlug)).Collection("members")
	var (
		query firestore.Query
		limit = 100
	)
	if params.Limit != nil {
		rawLimit := int(*params.Limit)
		if rawLimit > 0 {
			limit = rawLimit
		}
	}
	// we secretly query limit+1 and return the second-to-last item as the "after"
	query = membersRef.Limit(limit+1).OrderBy("DiscordID", firestore.Asc)
	if params.After != nil {
		query = query.StartAfter(*params.After)
	}
	if params.Snowflakes != nil {
		query = query.Where("DiscordID", "in", *params.Snowflakes)
	}
	// TODO: ACL check here
	snaps, err := query.Select().Documents(ctx).GetAll()
	if err != nil {
		log.Err(err).Msg("error getting members")
		fmt.Fprint(w, "error getting members")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	apiResponse := getMembersResponse{
		Users: make([]ChannelMember, minInt(len(snaps), limit)),
	}
	// if there are more
	if len(snaps) == limit+1 {
		apiResponse.After = snaps[limit-1].Ref.ID
	}
	for i := 0; i < len(apiResponse.Users); i++ {
		apiResponse.Users[i] = ChannelMember{Id: &snaps[i].Ref.ID}
	}
	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	if err = enc.Encode(apiResponse); err != nil {
		log.Err(err).Msg("error encoding JSON payload")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func minInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

const (
	accountNotConnected = "not connected"
	accountNotMember    = "not member"
)

func (a *apiImpl) CheckMembership(
	w http.ResponseWriter, r *http.Request,
	channelSlug ChannelSlugPathParam,
) {
	if !keyHasPermission(r, string(channelSlug)) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	var (
		token oauth2.Token
		ctx   = r.Context()
		enc   = json.NewEncoder(w)
	)
	var jsonBody CheckMembershipJSONBody
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Err(err).Msg("error reading request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logger := log.With().Str("user", jsonBody.Snowflake).Logger()
	// retrieve user token
	doc, err := fs.Collection("users").Doc(jsonBody.Snowflake).Collection("private").Doc("youtube").Get(ctx)
	if status.Code(err) == codes.NotFound {
		err = enc.Encode(map[string]interface{}{
			"member": false,
			"reason": accountNotConnected,
		})
		if err != nil {
			logger.Err(err).Msg("error encoding account-not-connected error message")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	} else if err != nil {
		logger.Err(err).Msg("error retrieving Youtube account token for user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = doc.DataTo(&token)
	if err != nil {
		logger.Err(err).Msg("error unmarshalling Youtube account token for user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	yt, err := youtube.NewService(ctx,
		append(
			ytClientOptions,
			option.WithTokenSource(youtubeOAuthConfig.TokenSource(ctx, &token)),
		)...,
	)
	if err != nil {
		logger.Err(err).Msg("error creating Youtube client")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// retrieve target videoID
	channelDocRef := fs.Collection("channels").Doc(string(channelSlug))
	doc, err = channelDocRef.Collection("check").Doc("check").Get(ctx)
	if status.Code(err) == codes.NotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Err(err).Msg("error retrieving Youtube channel membership check video")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var check common.ChannelCheck
	err = doc.DataTo(&check)
	if err != nil {
		logger.Err(err).Msg("error unmarshalling Youtube channel membership check video")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var isMember bool
	_, err = yt.CommentThreads.List([]string{"id"}).VideoId(check.VideoID).Do()
	if err != nil {
		if !strings.HasSuffix(err.Error(), "commentsDisabled") {
			logger.Err(err).Msg("actual error fetching comments for membership check video")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = nil
	} else {
		isMember = true
	}
	// update things and respond
	w.Header().Set("Content-Type", "application/json")
	if isMember {
		_, err = channelDocRef.Collection("members").Doc(jsonBody.Snowflake).
			Set(ctx, map[string]interface{}{
				"DiscordID": jsonBody.Snowflake,
				"Timestamp": time.Now(),
			})
		if err != nil {
			logger.Err(err).Msg("error setting membership")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = enc.Encode(map[string]bool{
			"member": true,
		})
	} else {
		err = enc.Encode(map[string]interface{}{
			"member": false,
			"reason": accountNotMember,
		})
	}
	if err != nil {
		logger.Err(err).Msg("error encoding response")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// API _
func API(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "API/")
	r.URL.RawPath = strings.TrimPrefix(r.URL.RawPath, "API/")
	apiHandler.ServeHTTP(w, r)
}

func init() {
	var (
		ctx  = context.Background()
		impl apiImpl
		err  error
	)
	zerolog.LevelFieldName = "severity"
	router := chi.NewRouter()
	router.Use(NewAuthHandler)
	apiHandler = HandlerFromMux(&impl, router)
	youtubeOAuthConfig = &oauth2.Config{
		ClientID:     mustLoadEnv(youtubeClientIDEnvName),
		ClientSecret: mustLoadEnv(youtubeClientSecretEnvName),
		Endpoint:     google.Endpoint,
		RedirectURL:  mustLoadEnv(youtubeRedirectURIEnvName),
		Scopes:       []string{"https://www.googleapis.com/auth/youtube.force-ssl"},
	}
	swagger, err = GetSwagger()
	if err != nil {
		log.Fatal().Err(err).Msg("error initializing Swagger spec")
	}
	fs, err = firestore.NewClient(ctx, mustLoadEnv("GCP_PROJECT"))
	if err != nil {
		log.Fatal().Err(err).Msg("error initializing Firestore")
	}
}

func mustLoadEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatal().Msgf("environment variable '%s' must not be empty", name)
	}
	return value
}
