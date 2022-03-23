package cmd

import (
	"context"
	"strconv"

	"github.com/member-gentei/member-gentei/gentei/apis"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/member-gentei/member-gentei/gentei/ent/usermembership"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	flagRepairChannelID      string
	flagRepairFallbackUserID uint64
)

// repairCmd represents the repair command
var repairCmd = &cobra.Command{
	Use:   "repair",
	Short: "Maintenance commands for various fun events",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ctx = context.Background()
			db  = mustOpenDB(ctx)
		)
		if flagRepairChannelID != "" {
			var (
				logger = log.With().Str("channelID", flagRepairChannelID).Logger()
				talent = db.YouTubeTalent.GetX(ctx, flagRepairChannelID)
			)
			// get a YouTube service from a member
			userID := db.User.Query().
				Where(
					user.YoutubeTokenNotNil(),
					user.HasMembershipsWith(
						usermembership.HasYoutubeTalentWith(youtubetalent.ID(flagRepairChannelID)),
					),
				).
				FirstIDX(ctx)
			if userID == 0 {
				logger.Warn().
					Str("userID", strconv.FormatUint(flagRepairFallbackUserID, 10)).
					Msg("no previously-verified members, using fallback user ID")
				userID = flagRepairFallbackUserID
			}
			svc, err := apis.GetYouTubeService(ctx, db, userID, getYouTubeConfig())
			if err != nil {
				logger.Fatal().Uint64("userID", userID).Err(err).Msg("error getting YouTube service for user")
			}
			var getNewMemberVideo bool
			if talent.MembershipVideoID == "" {
				logger.Info().Msg("membership video ID is blank, getting one...")
				getNewMemberVideo = true
			} else if !talent.Disabled.IsZero() {
				logger.Info().Msg("membership checks disabled, getting new video ID...")
				getNewMemberVideo = true
			}
			if getNewMemberVideo {
				videoID, err := apis.SelectRandomMembersOnlyVideoID(ctx, logger, svc, flagRepairChannelID)
				if err != nil {
					logger.Fatal().Err(err).Msg("error selecting membership video ID")
				}
				talent = db.YouTubeTalent.UpdateOne(talent).
					SetMembershipVideoID(videoID).
					ClearDisabled().
					SaveX(ctx)
				logger.Info().Str("videoID", videoID).Msg("set new MembershipVideoID")
			}
		}
	},
}

func init() {
	adminCmd.AddCommand(repairCmd)
	flags := repairCmd.Flags()
	flags.StringVarP(&flagRepairChannelID, "channel-id", "c", "", "YouTube channel ID to repair")
	// the default is my Discord user ID...
	flags.Uint64VarP(&flagRepairFallbackUserID, "fallback-user-id", "f", 196078350496825345, "Fallback user ID for membership video filling")
}
