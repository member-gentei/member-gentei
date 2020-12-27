package cmd

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagLeaveGuildID string
)

// leaveGuildCmd represents the leaveGuild command
var leaveGuildCmd = &cobra.Command{
	Use:   "leave-guild",
	Short: "Makes the bot leave a Discord guild.",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			token = viper.GetString("token")
		)
		if token == "" {
			log.Fatal().Msg("must specify a Discord token")
		}
		if flagLeaveGuildID == "" {
			log.Fatal().Msg("must specify a Discord guild ID")
		}
		dg, err := discordgo.New("Bot " + token)
		if err != nil {
			log.Fatal().Err(err).Msg("error creating Discord client")
		}
		err = dg.GuildLeave(flagLeaveGuildID)
		if err != nil {
			log.Fatal().Err(err).Msg("error leaving Discord guild")
		}
		log.Info().Str("guildID", flagLeaveGuildID).Msg("left Discord guild")
	},
}

func init() {
	rootCmd.AddCommand(leaveGuildCmd)
	flags := leaveGuildCmd.Flags()
	flags.StringVar(&flagLeaveGuildID, "id", "", "Discord guild ID")
}
