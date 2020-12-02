package dfu

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	app     *firebase.App
	appAuth *auth.Client
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

// DeleteFirebaseUser deletes a Firebase user when its corresponding Firestore user object deleted.
func DeleteFirebaseUser(ctx context.Context, event FirestoreEvent) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Interface("recovered", r).Msg("recovered from a panic")
			err, _ = r.(error)
		}
	}()
	userID := strings.Split(event.OldValue.Name, "/users/")[1]
	log.Info().Str("userID", userID).Msg("deleting Firebase user")
	err = appAuth.DeleteUser(ctx, userID)
	return
}

func init() {
	var err error
	zerolog.LevelFieldName = "severity"
	ctx := context.Background()
	app, err = firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("error initializing app")
	}
	appAuth, err = app.Auth(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("error initializing Firestore auth")
	}
}
