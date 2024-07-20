package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// prefillCmd represents the prefill command
var prefillCmd = &cobra.Command{
	Use:     "prefill",
	Short:   "Prefill Discord Guild objects so that nothing explodes on launch.",
	PreRunE: botCmd.PreRunE,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ctx = context.Background()
			db  = mustOpenDB(ctx)
		)
		session, err := discordgo.New(fmt.Sprintf("Bot %s", flagBotToken))
		if err != nil {
			log.Fatal().Err(err).Msg("error creating session")
		}
		userGuilds, err := session.UserGuilds(0, "", "", false, nil)
		if err != nil {
			log.Fatal().Err(err).Msg("error listing own guilds")
		}
		log.Info().Int("count", len(userGuilds)).Msg("joined guild count")
		for _, userDg := range userGuilds {
			dg, err := session.Guild(userDg.ID)
			if err != nil {
				log.Fatal().Err(err).Msg("error fetching guild details")
			}
			guildID, err := strconv.ParseUint(dg.ID, 10, 64)
			if err != nil {
				log.Fatal().Err(err).Msg("error parsing guildID as uint64")
			}
			if !db.Guild.Query().
				Where(guild.ID(guildID)).
				ExistX(ctx) {
				ownerID, err := strconv.ParseUint(dg.OwnerID, 10, 64)
				if err != nil {
					panic(err)
				}
				entDG := db.Guild.Create().
					SetID(guildID).
					SetName(dg.Name).
					SetIconHash(dg.Icon).
					SetAdminSnowflakes([]uint64{ownerID}).
					SaveX(ctx)
				log.Info().Interface("guild", entDG).Msg("inserted guild")
			}
		}
	},
}

func init() {
	adminCmd.AddCommand(prefillCmd)
}
