package cmd

import (
	"context"
	"strconv"
	"time"

	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	flagAdminPruneLeaveUnusedServers bool
	flagAdminPruneDryRun             bool
)

// pruneCmd represents the prune command
var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Deletes data, leaves servers, and performs other maintenance tasks to keep things tidy.",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ctx = context.Background()
			db  = mustOpenDB(ctx)
		)
		if !flagAdminPruneDryRun {
			log.Fatal().Msg("dry run mode only")
		}
		if flagAdminPruneLeaveUnusedServers {
			pruneUnusedServers(ctx, db, flagAdminPruneDryRun)
		}
	},
}

func pruneUnusedServers(ctx context.Context, db *ent.Client, dryRun bool) {
	guilds := db.Guild.Query().Where(
		guild.Not(guild.HasYoutubeTalents()),
		guild.FirstJoinedLT(time.Now().Add(-time.Hour*24*7)),
	).AllX(ctx)
	for _, dg := range guilds {
		log.Info().
			Str("guildID", strconv.FormatUint(dg.ID, 10)).
			Str("guildName", dg.Name).
			Msg("server has no configured YouTube channels, is candidate for leaving")
	}
}

func init() {
	adminCmd.AddCommand(pruneCmd)
	flags := pruneCmd.Flags()

	flags.BoolVar(&flagAdminPruneLeaveUnusedServers, "leave-unused-servers", false, "leaves servers that do not have roles configured for >7 days")
	flags.BoolVarP(&flagAdminPruneDryRun, "dry-run", "n", false, "print instead of taking action")
}
