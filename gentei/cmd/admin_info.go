package cmd

import (
	"context"

	"github.com/member-gentei/member-gentei/gentei/ent/guildrole"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Dumps membership info about a user/talent/etc in a nice format.",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ctx    = context.Background()
			db     = mustOpenDB(ctx)
			uid, _ = cmd.Flags().GetUint64("uid")
		)
		if uid != 0 {
			u := db.User.Query().
				WithMemberships().
				Where(user.ID(uid)).
				OnlyX(ctx)
			log.Info().Uint64("userID", uid).Str("fullName", u.FullName).Msg("got user")
			// get what they should have
			return
		}
		grs := db.GuildRole.Query().
			WithGuild().
			Where(
				guildrole.HasTalentWith(youtubetalent.DisabledIsNil()),
			).
			IDsX(ctx)
		log.Info().Uints64("guildRoles", grs).Int("len", len(grs)).Msg("matches")
	},
}

func init() {
	adminCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	flags := infoCmd.Flags()
	flags.Uint64("uid", 0, "user ID")
}
