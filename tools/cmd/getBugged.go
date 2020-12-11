package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/member-gentei/member-gentei/pkg/common"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/logging/logadmin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type oopsLog struct {
	Time   string
	UserID string
}

// getBuggedCmd represents the getBugged command
var getBuggedCmd = &cobra.Command{
	Use:   "getBugged",
	Short: "Get bugged users for the 2020-12-10",
	Run: func(cmd *cobra.Command, args []string) {
		affectedBeforeTime, err := time.Parse(time.RFC3339, "2020-12-09T00:00:00Z")
		if err != nil {
			panic(err)
		}
		// this happened to users before the 12/8 membership check
		ctx := context.Background()
		client, err := logadmin.NewClient(ctx, flagProjectID)
		if err != nil {
			log.Fatal().Err(err).Msg("error creating Logging client")
		}
		fs, err := firestore.NewClient(ctx, flagProjectID)
		if err != nil {
			log.Fatal().Err(err).Msg("error creating Firestore client")
		}
		var (
			// after     = "2020-12-08T09:44:40.398Z"
			after     = "2020-12-05T09:44:40.398Z"
			before    = "2020-12-08T20:44:40.398Z"
			logFilter = strings.TrimSpace(fmt.Sprintf(`
			timestamp > "%s"
			timestamp < "%s"
			logName="projects/member-gentei/logs/member-check"
			jsonPayload.message="Discord token invalid, deleting user"
		`, after, before))
		)
		var maybeUserIDs []string
		userIDMap := map[string]struct{}{}
		log.Info().Str("filter", logFilter).Msg("querying for logs...")
		iter := client.Entries(ctx, logadmin.Filter(logFilter))
		for {
			entry, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatal().Err(err).Msg("error iterating over log entries")
			}
			var payload struct {
				Time   string
				UserID string
			}
			err = jsonMarshalUnmarshal(entry.Payload, &payload)
			if _, dupe := userIDMap[payload.UserID]; dupe {
				continue
			}
			log.Debug().Interface("payload", payload).Send()
			userIDMap[payload.UserID] = struct{}{}
			maybeUserIDs = append(maybeUserIDs, payload.UserID)
		}
		log.Info().Int("count", len(maybeUserIDs)).Msg("possibly affected users")
		// check to see if this user is affected...
		var userIDs []string
		for i, userID := range maybeUserIDs {
			privateCollection := fs.Collection(common.UsersCollection).Doc(userID).
				Collection(common.PrivateCollection)
			doc, err := privateCollection.Doc("discord").Get(ctx)
			if c := status.Code(err); c == codes.NotFound {
				// affected: (still) deleted user
				userIDs = append(userIDs, maybeUserIDs[i])
				continue
			} else if err != nil {
				log.Fatal().Err(err).Str("userID", userID).Msg("error getting user")
			}
			var token oauth2.Token
			err = doc.DataTo(&token)
			if err != nil {
				log.Fatal().Err(err).Str("userID", userID).Msg("error unmarshalling user Discord token")
			}
			if token.Expiry.Before(affectedBeforeTime) {
				// affected: refresh token failed to renew
				userIDs = append(userIDs, maybeUserIDs[i])
				continue
			}
		}
		log.Info().Int("count", len(userIDs)).Msg("printing affected users to stdout as JSON list")
		json.NewEncoder(os.Stdout).Encode(userIDs)
	},
}

func jsonMarshalUnmarshal(thing interface{}, into interface{}) error {
	data, err := json.Marshal(thing)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, into)
	return err
}

func init() {
	rootCmd.AddCommand(getBuggedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getBuggedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getBuggedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
