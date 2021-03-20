package cmd

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/member-gentei/member-gentei/pkg/clients"
	"github.com/member-gentei/member-gentei/pkg/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

var (
	flagRestoreInput string
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore memberships revoked from a previous time",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ctx     = context.Background()
			fs, err = clients.NewRetryFirestoreClient(ctx, flagProjectID)
		)
		if err != nil {
			log.Fatal().Err(err).Msg("error creating Firestore client")
		}
		// assumes that this file was generated with something of the tune of:
		// fmt.Sprintf(`
		// 	logName="projects/member-gentei/logs/member-check"
		// 	AND jsonPayload.message = "membership lapsed for channel"
		// 	AND timestamp > "%s"`, time.Now().Format(time.RFC3339))),
		fn, err := os.Open(flagRestoreInput)
		if err != nil {
			log.Fatal().Err(err).Msg("error opening input file")
		}
		dec := json.NewDecoder(fn)
		var logEntries []struct {
			JSONPayload struct {
				UserID string
				Slug   string
			}
		}
		err = dec.Decode(&logEntries)
		if err != nil {
			log.Fatal().Err(err).Msg("error decoding log entries")
		}
		var putUsersBack = map[string][]string{}
		for _, entry := range logEntries {
			userID := entry.JSONPayload.UserID
			channelSlug := entry.JSONPayload.Slug
			putUsersBack[userID] = append(putUsersBack[userID], channelSlug)
		}
		p := mpb.New()
		bar := p.AddBar(
			int64(len(putUsersBack)),
			mpb.PrependDecorators(
				decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
			),
			mpb.AppendDecorators(
				decor.AverageSpeed(0, "%f/s"),
			),
		)
		for userID, slugs := range putUsersBack {
			var (
				logger      = log.With().Str("userID", userID).Logger()
				memberships = map[string]bool{}
			)
			selects, err := fs.CollectionGroup("members").Where("DiscordID", "==", userID).Select().Documents(ctx).GetAll()
			if err != nil {
				logger.Err(err).Msg("error querying for members CollectionGroup")
				return
			}
			for _, selected := range selects {
				slug := selected.Ref.Parent.Parent.ID
				memberships[slug] = true
			}
			for _, slug := range slugs {
				memberships[slug] = true
			}
			// update user
			membershipDocRefs := make([]*firestore.DocumentRef, 0, len(memberships))
			for slug := range memberships {
				membershipDocRefs = append(membershipDocRefs, fs.Collection("channels").Doc(slug))
			}
			nowIsh := time.Now().In(time.UTC)
			_, err = fs.Collection("users").Doc(userID).Update(ctx,
				[]firestore.Update{
					{
						Path:  "Memberships",
						Value: membershipDocRefs,
					},
					{
						Path:  "LastRefreshed",
						Value: nowIsh,
					},
				},
			)
			if err != nil {
				logger.Fatal().Err(err).Msg("error updating user object memberships")
			}
			cm := common.ChannelMember{
				DiscordID: userID,
				Timestamp: nowIsh,
			}
			for _, docRef := range membershipDocRefs {
				_, err := docRef.Collection(common.ChannelMemberCollection).
					Doc(userID).
					Set(ctx, cm)
				if err != nil {
					logger.Fatal().Err(err).Str("channelSlug", docRef.ID).Msg("error setting ChannelMember doc")
				}
			}
			bar.Increment()
		}
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	flags := restoreCmd.Flags()
	flags.StringVar(&flagRestoreInput, "file", "logs.json", "output of 'gcloud logging read --format=json'")
}
