package cmd

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/member-gentei/member-gentei/pkg/clients"
	"github.com/member-gentei/member-gentei/pkg/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	flagSetChannel         bool
	flagSetChannelSlug     string
	flagSetChannelID       string
	flagSetChannelVideoID  string
	flagSetLinkGuild       bool
	flagSetLinkGuildID     string
	flagSetLinkGuildRoleID string
)

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
			yt, err := common.GetYoutubeServerService(ctx, fs)
			if err != nil {
				log.Fatal().Err(err).Msg("error creating Youtube service")
			}
			clr, err := yt.Channels.List([]string{"snippet"}).Id(flagSetChannelID).Do()
			if err != nil {
				log.Fatal().Err(err).Msg("error getting Youtube channel")
			}
			channel := common.Channel{
				ChannelID:    flagSetChannelID,
				ChannelTitle: clr.Items[0].Snippet.Title,
				Thumbnail:    clr.Items[0].Snippet.Thumbnails.High.Url,
			}
			channelDocRef := fs.Collection(common.ChannelCollection).Doc(flagSetChannelSlug)
			_, err = channelDocRef.Set(ctx, channel)
			if err != nil {
				log.Fatal().Err(err).Msg("error setting Youtube channel")
			}
			if flagSetChannelVideoID != "" {
				log.Info().Msg("setting membership verification video")
				_, err = channelDocRef.Collection(common.ChannelCheckCollection).Doc("check").Set(ctx, map[string]string{
					"VideoID": flagSetChannelVideoID,
				})
				if err != nil {
					log.Fatal().Err(err).Msg("error setting membership verification video")
				}
			}
		} else if flagSetLinkGuild {
			if flagSetChannelSlug == "" || flagSetLinkGuildID == "" || flagSetLinkGuildRoleID == "" {
				log.Fatal().Msg("must specify channel slug, guild ID, and guild role ID")
			}
			log.Info().Str("channel", flagSetChannelSlug).Str("guild", flagSetLinkGuildID).Msg("linking Discord guild")
			// fetch existing
			var guild common.DiscordGuild
			guildRef := fs.Collection(common.DiscordGuildCollection).Doc(flagSetLinkGuildID)
			doc, err := guildRef.Get(ctx)
			if status.Code(err) == codes.NotFound {
				guild.MembershipRoles = map[string]string{}
			} else if err != nil {
				log.Fatal().Err(err).Msg("error getting existing Discord guild")
			}
			err = doc.DataTo(&guild)
			if err != nil {
				log.Fatal().Err(err).Msg("error unmarshalling Discord guild")
			}
			// set
			guild.ID = flagSetLinkGuildID
			guild.MembershipRoles[flagSetChannelSlug] = flagSetLinkGuildRoleID
			_, err = fs.Collection(common.DiscordGuildCollection).Doc(flagSetLinkGuildID).Set(ctx, guild)
			if err != nil {
				log.Fatal().Err(err).Msg("error linking Discord guild")
			}
			log.Warn().Msg("all users' guild memberships need to be checked if a link was added")
		}
	},
}

func completeChannelSlug(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var (
		ctx       = context.Background()
		directive = cobra.ShellCompDirectiveNoFileComp
	)
	fs, err := clients.NewRetryFirestoreClient(ctx, flagProjectID)
	if err != nil {
		return nil, directive
	}
	allDocs, err := fs.Collection(common.ChannelCollection).Select().Documents(ctx).GetAll()
	if err != nil {
		return nil, directive
	}
	candidates := make([]string, 0, len(allDocs))
	for _, doc := range allDocs {
		if strings.HasPrefix(doc.Ref.ID, toComplete) {
			candidates = append(candidates, doc.Ref.ID)
		}
	}
	return candidates, directive
}

func init() {
	rootCmd.AddCommand(setCmd)

	flags := setCmd.Flags()
	flags.BoolVar(&flagSetChannel, "channel", false, "set channel information")
	flags.StringVar(&flagSetChannelSlug, "channel-slug", "", "channel slug/document ID")
	flags.StringVar(&flagSetChannelID, "channel-id", "", "channel ID")
	flags.StringVar(&flagSetChannelVideoID, "channel-video-id", "", "membership video ID")
	flags.BoolVar(&flagSetLinkGuild, "link-guild", false, "link Discord guild to channel")
	flags.StringVar(&flagSetLinkGuildID, "guild-id", "", "Discord guild")
	flags.StringVar(&flagSetLinkGuildRoleID, "guild-role-id", "", "Discord guild")
	setCmd.RegisterFlagCompletionFunc("channel-slug", completeChannelSlug)
}
