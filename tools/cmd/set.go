package cmd

import (
	"context"

	"github.com/member-gentei/member-gentei/pkg/clients"
	"github.com/member-gentei/member-gentei/pkg/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	flagSetUID            string
	flagSetChannel        bool
	flagSetChannelSlug    string
	flagSetChannelID      string
	flagSetChannelVideoID string
	flagSetLinkGuild      bool
	flagSetLinkGuildID    string
)

type fsChannel struct {
	ChannelID    string
	ChannelTitle string
	Thumbnail    string
}

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set Discord <-> channel links, new Youtube channels, and more!",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		fs, err := clients.NewRetryFirestoreClient(ctx, flagProjectID)
		if err != nil {
			log.Fatal().Err(err).Msg("error creating Firestore client")
		}
		if flagSetChannel {
			if flagSetChannelSlug == "" || flagSetChannelID == "" {
				log.Fatal().Msg("must specify channel slug and ID")
			}
			log.Info().Str("slug", flagSetChannelSlug).Msg("setting channel")
			yt, err := common.GetYouTubeService(ctx, fs, flagSetUID)
			if err != nil {
				log.Fatal().Err(err).Msg("error creating Youtube service")
			}
			clr, err := yt.Channels.List([]string{"snippet"}).Id(flagSetChannelID).Do()
			if err != nil {
				log.Fatal().Err(err).Msg("error getting Youtube channel")
			}
			channel := fsChannel{
				ChannelID:    flagSetChannelID,
				ChannelTitle: clr.Items[0].Snippet.Title,
				Thumbnail:    clr.Items[0].Snippet.Thumbnails.High.Url,
			}
			channelDocRef := fs.Collection("channels").Doc(flagSetChannelSlug)
			_, err = channelDocRef.Set(ctx, channel)
			if err != nil {
				log.Fatal().Err(err).Msg("error setting Youtube channel")
			}
			if flagSetChannelVideoID != "" {
				log.Info().Msg("setting membership verification video")
				_, err = channelDocRef.Collection("check").Doc("check").Set(ctx, map[string]string{
					"VideoID": flagSetChannelVideoID,
				})
				if err != nil {
					log.Fatal().Err(err).Msg("error setting membership verification video")
				}
			}
		} else if flagSetLinkGuild {
			if flagSetChannelSlug == "" || flagSetLinkGuildID == "" {
				log.Fatal().Msg("must specify channel slug and guild ID")
			}
			log.Info().Str("channel", flagSetChannelSlug).Str("guild", flagSetLinkGuildID).Msg("linking Discord guild")
			_, err := fs.Collection("guilds").Doc(flagSetLinkGuildID).Set(ctx, map[string]interface{}{
				"Channel": fs.Collection("channels").Doc(flagSetChannelSlug),
				"ID":      flagSetLinkGuildID,
			})
			if err != nil {
				log.Fatal().Err(err).Msg("error linking Discord guild")
			}
			log.Warn().Msg("all users' guild memberships need to be checked if a link was added")
		}
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	flags := setCmd.Flags()
	flags.StringVar(&flagSetUID, "uid", "", "userID's tokens to use (use your own please)")
	flags.BoolVar(&flagSetChannel, "channel", false, "set channel information")
	flags.StringVar(&flagSetChannelSlug, "channel-slug", "", "channel slug/document ID")
	flags.StringVar(&flagSetChannelID, "channel-id", "", "channel ID")
	flags.StringVar(&flagSetChannelVideoID, "channel-video-id", "", "membership video ID")
	flags.BoolVar(&flagSetLinkGuild, "link-guild", false, "link Discord guild to channel")
	flags.StringVar(&flagSetLinkGuildID, "guild-id", "", "Discord guild")
}
