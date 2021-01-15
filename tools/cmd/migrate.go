package cmd

import (
	"context"

	"cloud.google.com/go/firestore"

	"github.com/member-gentei/member-gentei/pkg/clients"
	"github.com/member-gentei/member-gentei/pkg/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	flagMigrateDeleteDeprecated bool
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Short-lived schema migration command",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		fs, err := clients.NewRetryFirestoreClient(ctx, flagProjectID)
		if err != nil {
			log.Fatal().Err(err).Msg("error creating Firestore client")
		}
		docs, err := fs.Collection(common.DiscordGuildCollection).Documents(ctx).GetAll()
		if err != nil {
			log.Fatal().Err(err).Msg("error getting all Discord guilds")
		}
		for _, doc := range docs {
			var guild common.DiscordGuild
			if err = doc.DataTo(&guild); err != nil {
				log.Fatal().Err(err).Msg("error unmarshalling doc to DiscordGuild")
			}
			logger := log.With().Str("guildID", guild.ID).Str("name", guild.Name).Logger()
			if guild.MembershipRoles == nil {
				logger.Info().Msg("guild requires MembershipRoles migration")
				// _, err = doc.Ref.Update(ctx, []firestore.Update{
				// 	{
				// 		Path: "MembershipRoles",
				// 		Value: map[string]string{
				// 			guild.Channel.ID: guild.MemberRoleID,
				// 		},
				// 	},
				// })
				// if err != nil {
				// 	logger.Fatal().Err(err).Msg("error performing MembershipRoles migration")
				// }
			}
			if flagMigrateDeleteDeprecated {
				_, err = doc.Ref.Update(ctx, []firestore.Update{
					{
						Path:  "Channel",
						Value: firestore.Delete,
					},
					{
						Path:  "MemberRoleID",
						Value: firestore.Delete,
					},
				})
				if err != nil {
					logger.Fatal().Err(err).Msg("error removing deprecated fields")
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().BoolVar(&flagMigrateDeleteDeprecated, "delete-deprecated", false, "delete deprecated fields")
}
