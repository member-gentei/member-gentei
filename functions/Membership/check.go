// Package membership _
package membership

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/member-gentei/member-gentei/pkg/clients"
	"github.com/member-gentei/member-gentei/pkg/common"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
)

// FirestoreEvent is the payload of a Firestore event.
type FirestoreEvent struct {
	OldValue   FirestoreValue `json:"oldValue"`
	Value      FirestoreValue `json:"value"`
	UpdateMask struct {
		FieldPaths []string `json:"fieldPaths"`
	} `json:"updateMask"`
}

// FirestoreValue holds Firestore fields.
type FirestoreValue struct {
	CreateTime time.Time `json:"createTime"`
	// Fields is the data for this value. The type depends on the format of your
	// database. Log the interface{} value and inspect the result to see a JSON
	// representation of your database fields.
	Fields     json.RawMessage `json:"fields"`
	Name       string          `json:"name"`
	UpdateTime time.Time       `json:"updateTime"`
}

var (
	fs *firestore.Client
	yt *youtube.Service
)

// CheckMembershipWrite checks memberships when Youtube tokens are provisioned.
func CheckMembershipWrite(ctx context.Context, event FirestoreEvent) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Interface("recovered", r).Msg("recovered from a panic")
			err, _ = r.(error)
		}
	}()
	var resourcePath string
	if event.Value.Fields != nil {
		resourcePath = strings.Split(event.Value.Name, "/documents/")[1]
	} else {
		resourcePath = strings.Split(event.OldValue.Name, "/documents/")[1]
	}
	userDocRef := fs.Doc(resourcePath).Parent.Parent
	userID := userDocRef.ID
	logger := log.With().Str("userID", userID).Logger()
	if event.Value.Fields == nil {
		logger.Debug().Msg("ignoring delete")
		return
	}
	if event.OldValue.Fields != nil {
		oldToken, uErr := protoUnmarshalToken(event.OldValue.Fields)
		if uErr != nil {
			logger.Err(uErr).Msg("error unmarshalling OldValue.Fields")
			return uErr
		}
		newToken, uErr := protoUnmarshalToken(event.Value.Fields)
		if uErr != nil {
			logger.Err(uErr).Msg("error unmarshalling OldValue.Fields")
			return uErr
		}
		// ignore refresh tokens that did not change
		if oldToken.RefreshToken == newToken.RefreshToken {
			logger.Debug().Msg("ignoring write, refresh token did not change")
			return
		}
		// ignore if this is the same user - this is likely a reauth and/or just a new refresh token
		svc, err := common.GetYouTubeService(ctx, fs, userID)
		if err != nil {
			logger.Err(err).Msg("error creating YouTube service for new user token")
			return err
		}
		clr, err := svc.Channels.List([]string{"id"}).Mine(true).Do()
		if err != nil {
			logger.Err(err).Msg("error getting YouTube channel ID")
			return err
		}
		if len(clr.Items) == 0 {
			err = fmt.Errorf("unable to get channel")
			logger.Err(err).Msg("new token cannot get own channel ID")
			return err
		}
		var user common.DiscordIdentity
		userDoc, err := userDocRef.Get(ctx)
		if err != nil {
			logger.Err(err).Msg("error getting Discord user doc")
			return err
		}
		if err = userDoc.DataTo(&user); err != nil {
			logger.Err(err).Msg("error unmarshalling Discord user")
			return err
		}
		if user.YoutubeChannelID == clr.Items[0].Id {
			return nil
		}
	}
	_, err = common.EnforceMemberships(ctx, fs, &common.EnforceMembershipsOptions{
		Apply:   true,
		UserIDs: []string{userID},
	})
	return
}

func protoUnmarshalToken(fields json.RawMessage) (*oauth2.Token, error) {
	var protoToken struct {
		AccessToken struct {
			StringValue string
		}
		TokenType struct {
			StringValue string
		}
		RefreshToken struct {
			StringValue string
		}
		Expiry struct {
			TimestampValue string
		}
	}
	err := json.Unmarshal(fields, &protoToken)
	if err != nil {
		return nil, err
	}
	ts, err := time.Parse(time.RFC3339, protoToken.Expiry.TimestampValue)
	if err != nil {
		return nil, err
	}
	return &oauth2.Token{
		AccessToken:  protoToken.AccessToken.StringValue,
		TokenType:    protoToken.TokenType.StringValue,
		RefreshToken: protoToken.RefreshToken.StringValue,
		Expiry:       ts,
	}, nil
}

func init() {
	var (
		ctx = context.Background()
		err error
	)
	zerolog.LevelFieldName = "severity"
	fs, err = clients.NewRetryFirestoreClient(ctx, mustLoadEnv("GCP_PROJECT"))
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
