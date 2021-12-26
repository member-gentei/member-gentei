package cmd

import (
	"context"
	"errors"
	"time"

	"github.com/member-gentei/member-gentei/gentei/apis"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	flagUpsertYouTubeChannelID string
)

// upsertCmd represents the upsert command
var upsertCmd = &cobra.Command{
	Use:   "upsert",
	Short: "Upserts various objects.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if flagUpsertYouTubeChannelID == "" {
			return errors.New("must specify a flag")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ctx = context.Background()
			db  = mustOpenDB(ctx)
		)
		if flagUpsertYouTubeChannelID != "" {
			ogp, err := apis.GetYouTubeChannelOG(flagUpsertYouTubeChannelID)
			if err != nil {
				log.Fatal().Err(err).Msg("error getting OpenGraph data for YouTube channel")
			}
			log.Debug().Interface("ogp", ogp).Msg("OpenGraph data for channel")
			var thumbnailURL string
			for _, imageData := range ogp.Image {
				if imageData.Height == 900 {
					thumbnailURL = imageData.URL
					break
				}
				thumbnailURL = imageData.URL
			}
			err = db.YouTubeTalent.Create().
				SetID(flagUpsertYouTubeChannelID).
				SetChannelName(ogp.Title).
				SetThumbnailURL(thumbnailURL).
				SetLastUpdated(time.Now()).
				OnConflictColumns(youtubetalent.FieldID).
				UpdateNewValues().
				Exec(ctx)
			if err != nil {
				log.Fatal().Err(err).Msg("error upserting YouTubeTalent")
			}
			talent := db.YouTubeTalent.GetX(ctx, flagUpsertYouTubeChannelID)
			log.Info().Interface("talent", talent).Msg("upserted talent")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(upsertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upsertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	flags := upsertCmd.Flags()
	flags.StringVar(&flagUpsertYouTubeChannelID, "youtube", "", "YouTube channel ID")
}
