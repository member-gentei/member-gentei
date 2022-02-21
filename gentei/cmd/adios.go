package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cbroglie/mustache"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	flagAdiosDryRun bool
)

var (
	adiosMessageTemplate *mustache.Template
)

// adiosCmd represents the adios command
var adiosCmd = &cobra.Command{
	Use:   "adios",
	Short: "says goodbye, leaves",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		flagBotToken = os.Getenv(envNameDiscordBotToken)
		if flagBotToken == "" {
			return fmt.Errorf("must specify env var '%s'", envNameDiscordBotToken)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		session, err := discordgo.New(fmt.Sprintf("Bot %s", flagBotToken))
		if err != nil {
			log.Fatal().Err(err).Msg("error calling discordgo.New")
		}
		// get all of our guilds
		var (
			afterID        string
			serverOwnerMap = map[string][]*discordgo.Guild{}
		)
		for {
			userGuilds, err := session.UserGuilds(100, "", afterID)
			if err != nil {
				log.Fatal().Err(err).Msg("error getting UserGuilds page")
			}
			for _, userGuild := range userGuilds {
				guild, err := session.Guild(userGuild.ID)
				if err != nil {
					log.Fatal().Err(err).Msg("error getting Guild info")
				}
				serverOwnerMap[guild.OwnerID] = append(serverOwnerMap[guild.OwnerID], guild)
			}
			if len(userGuilds) != 100 {
				break
			} else {
				afterID = userGuilds[len(userGuilds)-1].ID
			}
		}
		for ownerID, guilds := range serverOwnerMap {
			guildNames := make([]string, len(guilds))
			for i := range guilds {
				guildNames[i] = guilds[i].Name
			}
			logger := log.With().Str("ownerID", ownerID).Strs("guildNames", guildNames).Logger()
			userChannel, err := session.UserChannelCreate(ownerID)
			if err != nil {
				logger.Fatal().Err(err).Msg("error creating UserChannel")
			}
			content, err := adiosMessageTemplate.Render(map[string][]string{
				"serverNames": guildNames,
			})
			content = strings.TrimSpace(content)
			if err != nil {
				logger.Fatal().Err(err).Msg("error rendering message template")
			}
			if flagAdiosDryRun {
				logger.Info().Str("content", content).Msg("DRY RUN: would have sent message and left servers")
			} else {
				_, err = session.ChannelMessageSend(userChannel.ID, content)
				if err != nil {
					logger.Fatal().Err(err).Msg("error sending parting message")
				}
				for _, guild := range guilds {
					err = session.GuildLeave(guild.ID)
					if err != nil {
						logger.Fatal().Err(err).Str("guildIDStr", guild.ID).Msg("error leaving guild")
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(adiosCmd)

	flags := adiosCmd.Flags()
	flags.BoolVarP(&flagAdiosDryRun, "dry-run", "n", false, "do not actually message anyone")

	var err error
	adiosMessageTemplate, err = mustache.ParseString(strings.TrimSpace(`
Hello! version 2 of Gentei, the fully automatic YouTube membership role assignment bot, is coming soon! To prepare for its soft launch today, this bot will be leaving all servers that it was previously invited to. 

The YouTube API issue that disabled Gentei for months last year has been fixed, so we've taken this opportunity to rewrite the bot to use /gentei slash commands and rewrite the site to make self-enrollment and self-management possible! Please visit https://gentei.tindabox.net if you're interested in participating in the soft launch - otherwise, please keep your eyes peeled a week or two from now when v2 regains feature parity with v1.

This message is being sent to you because this bot has determined that you are the owner of the following servers:
{{#serverNames}}
* {{ . }}
{{/serverNames}}
`))
	if err != nil {
		panic(err)
	}
}
