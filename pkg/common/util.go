package common

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"google.golang.org/api/option"

	"libs.altipla.consulting/tokensource"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

var (
	youTubeOAuthConfig *oauth2.Config
	discordOAuthConfig = &oauth2.Config{}
	youTubeAPIKey      string
)

// GetYouTubeService initializes a YouTube service for a user.
func GetYouTubeService(ctx context.Context, fs *firestore.Client, userID string) (svc *youtube.Service, err error) {
	if youTubeOAuthConfig == nil {
		err = loadYoutubeConfig(ctx, fs)
		if err != nil {
			return
		}
	}
	// get and load the user token
	var (
		token  oauth2.Token
		docRef = fs.Collection(UsersCollection).Doc(userID).
			Collection(PrivateCollection).Doc("youtube")
	)
	doc, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}
	err = doc.DataTo(&token)
	if err != nil {
		log.Err(err).Msg("error unmarshalling YouTube token")
		return nil, err
	}
	notifyHook := tokensource.NewNotifyHook(
		ctx, youTubeOAuthConfig, &token,
		func(newToken *oauth2.Token) error {
			// save the new token
			log.Debug().Str("userID", userID).Msg("saving newly refreshed YouTube token for user")
			_, err := docRef.Set(ctx, newToken)
			return err
		},
	)
	return youtube.New(notifyHook.Client(ctx))
}

// GetDiscordHTTPClient creates a Discord HTTP client for a user.
func GetDiscordHTTPClient(ctx context.Context, fs *firestore.Client, userID string) (client *http.Client, err error) {
	if discordOAuthConfig.ClientID == "" {
		err = loadDiscordConfig(ctx, fs)
		if err != nil {
			return
		}
	}
	var (
		token  oauth2.Token
		docRef = fs.Collection(UsersCollection).Doc(userID).
			Collection(PrivateCollection).Doc("discord")
	)
	doc, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}
	err = doc.DataTo(&token)
	if err != nil {
		log.Err(err).Msg("error unmarshalling Discord token")
		return nil, err
	}
	notifyHook := tokensource.NewNotifyHook(
		ctx, discordOAuthConfig, &token,
		func(newToken *oauth2.Token) error {
			// save the new token
			log.Debug().Str("userID", userID).Msg("saving newly refreshed Discord token for user")
			_, err := docRef.Set(ctx, newToken)
			return err
		},
	)
	client = notifyHook.Client(ctx)
	return
}

// GetYoutubeServerService initializes a YouTube service using the project API key.
func GetYoutubeServerService(ctx context.Context, fs *firestore.Client) (svc *youtube.Service, err error) {
	if youTubeAPIKey == "" {
		snap, err := fs.Collection("config").Doc("youtube-server").Get(ctx)
		if err != nil {
			log.Err(err).Msg("error getting YouTube API key")
			return nil, err
		}
		var apiKey struct {
			Data string
		}
		err = snap.DataTo(&apiKey)
		if err != nil {
			log.Err(err).Msg("error unmarshalling YouTube API key")
			return nil, err
		}
		youTubeAPIKey = apiKey.Data
	}
	return youtube.NewService(ctx, option.WithAPIKey(youTubeAPIKey))
}

func loadYoutubeConfig(ctx context.Context, fs *firestore.Client) error {
	var clientJSON struct {
		Data string
	}
	snap, err := fs.Collection("config").Doc("youtube").Get(ctx)
	if err != nil {
		log.Err(err).Msg("error loading YouTube client config")
		return err
	}
	err = snap.DataTo(&clientJSON)
	if err != nil {
		log.Err(err).Msg("error unmarshalling YouTube client config")
		return err
	}
	youTubeOAuthConfig, err = google.ConfigFromJSON([]byte(clientJSON.Data), youtube.YoutubeForceSslScope)
	if err != nil {
		log.Err(err).Msg("error creating YouTube oauth2.Config")
		return err
	}
	return nil
}

func loadDiscordConfig(ctx context.Context, fs *firestore.Client) error {
	snap, err := fs.Collection("config").Doc("discord").Get(ctx)
	if err != nil {
		return err
	}
	err = snap.DataTo(discordOAuthConfig)
	if err != nil {
		log.Err(err).Msg("error unmarshalling Discord OAuth config")
		return err
	}
	if discordOAuthConfig.ClientID == "" {
		return fmt.Errorf("Discord OAuth config does not exist")
	}
	return nil
}

func getMemberVideoIDs(ctx context.Context, fs *firestore.Client) (slug2Video map[string]string, err error) {
	snaps, err := fs.CollectionGroup(ChannelCheckCollection).Documents(ctx).GetAll()
	if err != nil {
		return
	}
	slug2Video = make(map[string]string, len(snaps))
	for _, snap := range snaps {
		var check ChannelCheck
		err = snap.DataTo(&check)
		if err != nil {
			log.Err(err).Str("path", snap.Ref.Path).Msg("error unmarshalling ChannelCheck")
			return
		}
		slug2Video[snap.Ref.Parent.Parent.ID] = check.VideoID
	}
	return
}

// revokeYoutubeToken can take a regular or refresh token and emits warnings instead of erroring out.
func revokeYoutubeToken(refreshToken string, logger zerolog.Logger) {
	r, err := http.Post(
		fmt.Sprintf("https://oauth2.googleapis.com/revoke?token=%s", refreshToken),
		"application/x-www-form-urlencoded",
		nil,
	)
	if err != nil {
		logger.Warn().Err(err).Msg("error revoking YouTube token")
	}
	if r.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(r.Body)
		logger.Warn().Bytes("body", body).Int("statusCode", r.StatusCode).Msg("non-200 response while revoking YouTube token")
	}
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
