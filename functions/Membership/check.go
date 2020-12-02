// Package membership _
package membership

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/mark-ignacio/member-gentei/pkg/common"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
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
)

// CheckMembershipWrite checks memberships when Youtube tokens are provisioned.
func CheckMembershipWrite(ctx context.Context, event FirestoreEvent) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Interface("recovered", r).Msg("recovered from a panic")
			err, _ = r.(error)
		}
	}()
	if event.Value.Fields == nil {
		log.Debug().Msg("ignoring delete")
		return
	}
	resourcePath := strings.Split(event.Value.Name, "/documents/")[1]
	log.Info().Str("resourcePath", resourcePath).Msg("handling resource")
	userDocRef := fs.Doc(resourcePath).Parent.Parent
	err = common.EnforceMemberships(ctx, fs, &common.EnforceMembershipsOptions{
		Apply:   true,
		UserIDs: []string{userDocRef.ID},
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
	zerolog.LevelFieldName = "severity"
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("error initializing app")
	}
	fs, err = app.Firestore(ctx)
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
