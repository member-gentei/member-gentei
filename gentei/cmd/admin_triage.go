package cmd

import (
	"context"
	"fmt"

	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
	"github.com/rs/zerolog/log"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// triageCmd represents the triage command
var triageCmd = &cobra.Command{
	Use:   "triage",
	Short: "Displays data about a server + configuration for help with debugging things.",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			guildID, _   = cmd.Flags().GetUint64("guild")
			userID, _    = cmd.Flags().GetUint64("user")
			youtubeID, _ = cmd.Flags().GetString("youtube-channel")
			ctx          = context.Background()
			db           = mustOpenDB(ctx)
		)
		if guildID != 0 {
			dg := db.Guild.Query().
				Where(guild.ID(guildID)).
				WithRoles(
					func(grq *ent.GuildRoleQuery) {
						grq.WithTalent()
					},
				).
				OnlyX(ctx)
			tw := table.NewWriter()
			tw.AppendHeader(table.Row{
				"Guild ID",
				"Guild Name",
				"Role ID",
				"Role Name",
				"Talent URL",
				"Talent Name",
			})
			for i, role := range dg.Edges.Roles {
				var row table.Row
				if i == 0 {
					row = append(row, dg.ID, dg.Name)
				} else {
					row = append(row, "", "")
				}
				row = append(row, role.ID, role.Name, role.Edges.Talent.ID, role.Edges.Talent.ChannelName)
				tw.AppendRow(row)
			}
			fmt.Println(tw.Render())
		} else if userID != 0 {
			u := db.User.Query().
				Where(user.ID(userID)).
				WithMemberships(func(umq *ent.UserMembershipQuery) {
					umq.WithYoutubeTalent(func(yttq *ent.YouTubeTalentQuery) {
						yttq.Select(youtubetalent.FieldID, youtubetalent.FieldChannelName)
					})
				}).
				OnlyX(ctx)
			log.Info().
				Uint64("userID", userID).
				Str("user", u.FullName).
				Time("lastCheck", u.LastCheck).
				Msg("fetched user membership information")
			tw := table.NewWriter()
			tw.AppendHeader(table.Row{
				"Channel ID",
				"Channel name",
				"Last verified",
				"First failed",
				"Fail count",
			})
			for _, membership := range u.Edges.Memberships {
				tw.AppendRow(table.Row{
					membership.Edges.YoutubeTalent.ID,
					membership.Edges.YoutubeTalent.ChannelName,
					membership.LastVerified,
					membership.FirstFailed,
					membership.FailCount,
				})
			}
			fmt.Println(tw.Render())
		} else if youtubeID != "" {
			y := db.YouTubeTalent.GetX(ctx, youtubeID)
			tw := table.NewWriter()
			tw.AppendHeader(table.Row{
				"Channel ID",
				"Channel Name",
				"Membership Video ID",
				"Last Updated",
				"Last Membership Video ID Miss",
				"Disabled",
				"Disabled Permanently",
			})
			tw.AppendRow(table.Row{
				y.ID,
				y.ChannelName,
				y.MembershipVideoID,
				y.LastUpdated,
				y.LastMembershipVideoIDMiss,
				y.Disabled,
				y.DisabledPermanently,
			})
			fmt.Println(tw.Render())
		} else {
			log.Fatal().Msg("--guild or --user or --youtube-channel required")
		}
	},
}

func init() {
	adminCmd.AddCommand(triageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// triageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	flags := triageCmd.Flags()
	flags.Uint64P("guild", "s", 0, "Discord guild (+server) ID")
	flags.Uint64P("user", "u", 0, "Discord user ID")
	flags.StringP("youtube-channel", "y", "", "YouTube channel ID")
}
