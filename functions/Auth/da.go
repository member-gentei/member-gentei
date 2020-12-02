// Package auth handles Discord + YouTube authentication.
package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"

	"github.com/member-gentei/member-gentei/pkg/common"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/errorutils"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

const (
	discordClientIDEnvName     = "DISCORD_CLIENT_ID"
	discordClientSecretEnvName = "DISCORD_CLIENT_SECRET"
	discordRedirectURIEnvName  = "DISCORD_REDIRECT_URI"
	discordCollectionEnvName   = "DISCORD_COLLECTION"
	youtubeClientIDEnvName     = "YOUTUBE_CLIENT_ID"
	youtubeClientSecretEnvName = "YOUTUBE_CLIENT_SECRET"
	youtubeRedirectURIEnvName  = "YOUTUBE_REDIRECT_URI"

	discordAuthURL     = "https://discord.com/api/oauth2/authorize"
	discordMeURL       = "https://discord.com/api/users/@me"
	discordMeGuildsURL = "https://discord.com/api/users/@me/guilds"
	discordTokenURL    = "https://discord.com/api/oauth2/token"
	discordRevokeURL   = "https://discord.com/api/oauth2/token/revoke"
)

var (
	app     *firebase.App
	appAuth *auth.Client
	fs      *firestore.Client

	discordOAuthConfig *oauth2.Config
	discordCollection  string
	youtubeOAuthConfig *oauth2.Config
)

// Auth does the third leg of the OAuuth dance via an XHR.
func Auth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "https://member-gentei.tindabox.net")
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "OPTIONS, POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	switch svc := r.URL.Query().Get("service"); svc {
	case "discord":
		discordAuth(w, r)
	case "youtube":
		youtubeAuth(w, r)
	default:
		fmt.Fprintf(w, "unrecognized service: '%s'", svc)
		w.WriteHeader(http.StatusBadRequest)
	}

}

func discordAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseForm()
	if err != nil {
		log.Err(err).Msg("error parsing form")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	oauthCode := r.Form.Get("code")
	// oauthState := query.Get("state")
	if oauthCode == "" {
		w.Write([]byte("OAuth code not found"))
		w.WriteHeader(http.StatusForbidden)
		return
	}
	token, err := discordOAuthConfig.Exchange(ctx, oauthCode, oauth2.AccessTypeOffline)
	// get identity
	if err != nil {
		log.Err(err).Msg("error exchanging OAuth token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	client := discordOAuthConfig.Client(ctx, token)
	response, err := client.Get(discordMeURL)
	if err != nil {
		log.Err(err).Msg("error getting Discord identity")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Err(err).Msg("error reading response body for Discord identity")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if response.StatusCode >= 400 {
		log.Error().Int("status_code", response.StatusCode).
			Bytes("body", body).
			Msg("HTTP error getting Discord identity")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var meResponse struct {
		ID            string
		Username      string
		Discriminator string
	}
	err = json.Unmarshal(body, &meResponse)
	if err != nil {
		log.Err(err).Msg("error unmarshalling response body as JSON for Discord identity")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// write discord identity stuff
	discordDoc := fs.Collection(discordCollection).Doc(meResponse.ID)
	discordIdentity := common.DiscordIdentity{
		UserID:        meResponse.ID,
		Username:      meResponse.Username,
		Discriminator: meResponse.Discriminator,
	}
	// if the user currently exists, retain memberships, candidate channels, and YouTubeChannelID
	existingDoc, err := discordDoc.Get(ctx)
	if err != nil {
		if c := status.Code(err); c != codes.Unknown && c != codes.NotFound {
			log.Err(err).Msg("GRPC error saving Discord identity to Firestore")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = nil
		channelRefs, err := getCandidateChannels(ctx, fs, client)
		if err != nil {
			log.Err(err).Msg("error getting candidate channels")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		discordIdentity.CandidateChannels = channelRefs
	} else {
		var existingDiscordIdentity common.DiscordIdentity
		err = existingDoc.DataTo(&existingDiscordIdentity)
		if err != nil {
			log.Err(err).Msg("error unmarshalling existing Discord identity")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		discordIdentity.Memberships = existingDiscordIdentity.Memberships
		discordIdentity.CandidateChannels = existingDiscordIdentity.CandidateChannels
		discordIdentity.YoutubeChannelID = existingDiscordIdentity.YoutubeChannelID
	}
	_, err = discordDoc.Set(ctx, discordIdentity)
	if err != nil {
		log.Err(err).Msg("error saving Discord identity to Firestore")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = discordDoc.Collection("private").Doc("discord").Set(ctx, token)
	if err != nil {
		log.Err(err).Msg("error saving Discord token to Firestore")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jwt, err := appAuth.CustomToken(ctx, meResponse.ID)
	if err != nil {
		log.Err(err).Msg("error creating custom token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	payload, err := json.Marshal(map[string]string{"jwt": jwt})
	if err != nil {
		log.Err(err).Msg("error marshaling custom token to JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// get or create user
	_, err = appAuth.GetUser(ctx, meResponse.ID)
	if err != nil {
		if errorutils.IsNotFound(err) {
			toCreate := &auth.UserToCreate{}
			toCreate = toCreate.UID(meResponse.ID)
			_, err = appAuth.CreateUser(ctx, toCreate)
			if err != nil {
				log.Err(err).Msg("error creating user")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			log.Err(err).Msg("error checking for user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		toUpdate := &auth.UserToUpdate{}
		toUpdate = toUpdate.DisplayName(fmt.Sprintf("%s#%s", meResponse.Username, meResponse.Discriminator))
		_, err := appAuth.UpdateUser(ctx, meResponse.ID, toUpdate)
		if err != nil {
			log.Err(err).Msg("error updating existing username")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.Write(payload)
}

func youtubeAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseForm()
	if err != nil {
		log.Err(err).Msg("error parsing form")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	currentUser, err := appAuth.VerifyIDTokenAndCheckRevoked(ctx, r.Form.Get("jwt"))
	if err != nil {
		log.Err(err).Msg("could not verify gentei JWT")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	code := r.Form.Get("code")
	token, err := youtubeOAuthConfig.Exchange(ctx, code, oauth2.AccessTypeOffline)
	if err != nil {
		log.Err(err).Str("userID", currentUser.UID).Msg("could not exchange YouTube OAuth2 token")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// verify that we're connecting the same channel
	var identity common.DiscordIdentity
	userDocRef := fs.Collection(discordCollection).Doc(currentUser.UID)
	userDoc, err := userDocRef.Get(ctx)
	if err != nil {
		log.Err(err).Msg("error getting current user data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = userDoc.DataTo(&identity)
	if err != nil {
		log.Err(err).Msg("error unmarshalling current user data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ytClient, err := youtube.New(youtubeOAuthConfig.Client(ctx, token))
	if err != nil {
		log.Err(err).Msg("error creating Youtube client")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	clr, err := ytClient.Channels.List([]string{"id"}).Mine(true).Do()
	if err != nil {
		log.Err(err).Msg("error listing own channel")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var channelID string
	if len(clr.Items) == 0 {
		log.Error().Str("uid", currentUser.UID).Msg("no channels for user, cannot check duplicate enrollment")
	} else {
		channelID = clr.Items[0].Id
		// check that nobody else is using this channelID
		snaps, err := fs.Collection(common.UsersCollection).
			Where("YoutubeChannelID", "==", channelID).
			Select().Documents(ctx).GetAll()
		if err != nil {
			log.Err(err).Msg("error querying for existing YouTube channel ID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// assumes that we don't already have multiple users sharing one YouTube channel
		if len(snaps) > 0 && snaps[0].Ref.ID != identity.UserID {
			existingIDs := make([]string, len(snaps))
			for i := range snaps {
				snapID := snaps[i].Ref.ID
				if snapID == identity.UserID {
					continue
				}
				existingIDs[i] = snaps[i].Ref.ID
			}
			log.Info().Str("authUser", identity.UserID).
				Strs("existing", existingIDs).
				Msg("denying channel ID associated with other user(s)")
			err = writeJSONError(w, "YouTube channel is already associated with another user")
			if err != nil {
				log.Err(err).Msg("error writing error message")
			}
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
	// great! store it
	_, err = userDocRef.Collection("private").Doc("youtube").Set(ctx, token)
	if err != nil {
		log.Err(err).Msg("error storing Youtube credentials")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	identity.YoutubeChannelID = channelID
	_, err = userDocRef.Set(ctx, identity)
	if err != nil {
		log.Err(err).Msg("error setting YoutubeConnection as connected")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getCandidateChannels(
	ctx context.Context,
	fs *firestore.Client,
	httpClient *http.Client,
) (channelRefs []*firestore.DocumentRef, err error) {
	response, err := httpClient.Get(discordMeGuildsURL)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		log.Warn().Int("code", response.StatusCode).Bytes("body", body).
			Msg("non-200 status getting discord guilds for user")
		return
	}
	var guildMemberships []struct {
		ID string
		// Permissions string
	}
	err = json.Unmarshal(body, &guildMemberships)
	if err != nil {
		return
	}
	if len(guildMemberships) == 0 {
		return
	}
	ids := make([]string, len(guildMemberships))
	for i := range guildMemberships {
		ids[i] = guildMemberships[i].ID
	}
	// batches of 10, so up to 10 calls per user (yeesh)
	var guildSnaps []*firestore.DocumentSnapshot
	for i := 0; i < len(ids); i += 10 {
		j := i + 10
		if j > len(ids) {
			j = len(ids)
		}
		var snapBatch []*firestore.DocumentSnapshot
		idBatch := ids[i:j]
		if len(idBatch) == 0 {
			continue
		}
		snapBatch, err = fs.Collection(common.DiscordGuildCollection).Where("ID", "in", idBatch).Documents(ctx).GetAll()
		if err != nil {
			return
		}
		guildSnaps = append(guildSnaps, snapBatch...)
	}
	candidateMap := make(map[string]*firestore.DocumentRef)
	for _, snap := range guildSnaps {
		var partialGuild struct {
			Channel *firestore.DocumentRef
		}
		err = snap.DataTo(&partialGuild)
		if err != nil {
			return
		}
		candidateMap[partialGuild.Channel.Path] = partialGuild.Channel
	}
	channelRefs = make([]*firestore.DocumentRef, 0, len(candidateMap))
	for _, candidate := range candidateMap {
		channelRefs = append(channelRefs, candidate)
	}
	// sort by docID
	sort.Slice(channelRefs, func(i, j int) bool {
		return sort.StringsAreSorted([]string{channelRefs[i].ID, channelRefs[j].ID})
	})
	return
}

func init() {
	var err error
	zerolog.LevelFieldName = "severity"
	discordClientID := mustLoadEnv(discordClientIDEnvName)
	discordClientSecret := mustLoadEnv(discordClientSecretEnvName)
	discordCollection = mustLoadEnv(discordCollectionEnvName)
	discordOAuthConfig = &oauth2.Config{
		ClientID:     discordClientID,
		ClientSecret: discordClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://discord.com/api/oauth2/authorize",
			TokenURL:  "https://discord.com/api/oauth2/token",
			AuthStyle: oauth2.AuthStyleInHeader,
		},
		RedirectURL: mustLoadEnv(discordRedirectURIEnvName),
		Scopes:      []string{"identify", "guilds"},
	}
	youtubeOAuthConfig = &oauth2.Config{
		ClientID:     mustLoadEnv(youtubeClientIDEnvName),
		ClientSecret: mustLoadEnv(youtubeClientSecretEnvName),
		Endpoint:     google.Endpoint,
		RedirectURL:  mustLoadEnv(youtubeRedirectURIEnvName),
		Scopes:       []string{"https://www.googleapis.com/auth/youtube.force-ssl"},
	}
	ctx := context.Background()
	app, err = firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("error initializing app")
	}
	fs, err = app.Firestore(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("error initializing Firestore")
	}
	appAuth, err = app.Auth(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("error initializing Firestore auth")
	}
}

func writeJSONError(writer io.Writer, message string) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "")
	return encoder.Encode(map[string]string{"error": message})
}

func mustLoadEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatal().Msgf("environment variable '%s' must not be empty", name)
	}
	return value
}
