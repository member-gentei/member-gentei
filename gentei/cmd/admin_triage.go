package cmd

import (
	"context"
	"fmt"

	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/rs/zerolog/log"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// triageCmd represents the triage command
var triageCmd = &cobra.Command{
	Use:   "triage",
	Short: "Displays data about a server + configuration for help with debugging things.",
	Run: func(cmd *cobra.Command, args []string) {
		guildID, _ := cmd.Flags().GetUint64("guild")
		if guildID == 0 {
			log.Fatal().Msg("--guild required")
		}
		var (
			ctx = context.Background()
			db  = mustOpenDB(ctx)
		)
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
}
