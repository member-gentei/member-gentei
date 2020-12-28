package dc

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/member-gentei/member-gentei/pkg/clients"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// DisconnectYTAccount completes disconnection of a YouTube account upon deletion of a
// corresponding /users/{doc}/private/youtube document.
func DisconnectYTAccount(ctx context.Context, event FirestoreEvent) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Interface("recovered", r).Msg("recovered from a panic")
			err, _ = r.(error)
		}
	}()
	resourcePath := strings.Split(event.OldValue.Name, "/documents/")[1]
	userDocRef := fs.Doc(resourcePath).Parent.Parent
	_, err = userDocRef.Update(ctx, []firestore.Update{
		firestore.Update{
			Path:  "Memberships",
			Value: []*firestore.DocumentRef{},
		},
		firestore.Update{
			Path:  "YoutubeChannelID",
			Value: "",
		},
	})
	if status.Code(err) == codes.NotFound {
		log.Info().Str("userID", userDocRef.ID).Err(err).Msg("user doc deleted, no need to update")
		err = nil
	}
	return
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
