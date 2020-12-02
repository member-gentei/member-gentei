package common

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

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
	var token oauth2.Token
	doc, err := fs.Collection(UsersCollection).Doc(userID).
		Collection(PrivateCollection).Doc("youtube").Get(ctx)
	if err != nil {
		return nil, err
	}
	err = doc.DataTo(&token)
	if err != nil {
		log.Err(err).Msg("error unmarshalling YouTube token")
		return nil, err
	}
	return youtube.New(youTubeOAuthConfig.Client(ctx, &token))
}

// GetDiscordHTTPClient creates a Discord HTTP client for a user.
func GetDiscordHTTPClient(ctx context.Context, fs *firestore.Client, userID string) (client *http.Client, err error) {
	if discordOAuthConfig.ClientID == "" {
		err = loadDiscordConfig(ctx, fs)
		if err != nil {
			return
		}
	}
	var token oauth2.Token
	doc, err := fs.Collection(UsersCollection).Doc(userID).
		Collection(PrivateCollection).Doc("discord").Get(ctx)
	if err != nil {
		return nil, err
	}
	err = doc.DataTo(&token)
	if err != nil {
		log.Err(err).Msg("error unmarshalling Discord token")
		return nil, err
	}
	client = discordOAuthConfig.Client(ctx, &token)
	return
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
