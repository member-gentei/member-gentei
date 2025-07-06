package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type backupFile struct {
	GuildsWithRoles      []*ent.Guild
	TalentsWithRoles     []*ent.YouTubeTalent
	UsersWithMemberships []*ent.User

	// Users <-> Guilds are populated on-demand, so we don't need to store that ever
}

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Take application-level backups",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ctx           = context.Background()
			db            = mustOpenDB(ctx)
			backup        backupFile
			flagOutput, _ = cmd.Flags().GetString("output")
		)
		// get all guilds
		backup.GuildsWithRoles = db.Guild.Query().
			WithRoles().
			AllX(ctx)
		backup.TalentsWithRoles = db.YouTubeTalent.Query().
			WithRoles().
			AllX(ctx)
		backup.UsersWithMemberships = db.User.Query().
			WithMemberships(func(umq *ent.UserMembershipQuery) {
				umq.WithRoles().WithYoutubeTalent()
			}).
			AllX(ctx)
		// open the file and write it
		of, err := os.OpenFile(flagOutput, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o640)
		if err != nil {
			log.Fatal().Err(err).Msg("error opening output file for backup")
		}
		defer of.Close()
		log.Info().
			Int("guilds", len(backup.GuildsWithRoles)).
			Int("talents", len(backup.TalentsWithRoles)).
			Int("users", len(backup.UsersWithMemberships)).
			Str("path", flagOutput).
			Msg("writing backup file")
		if err := json.NewEncoder(of).Encode(backup); err != nil {
			log.Fatal().Err(err).Msg("error encoding backup file")
		}
	},
}

func init() {
	adminCmd.AddCommand(backupCmd)
	flags := backupCmd.Flags()
	flags.StringP("output", "o", "backup.json", "backup output file")
}
