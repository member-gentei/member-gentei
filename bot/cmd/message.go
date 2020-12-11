package cmd

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/member-gentei/member-gentei/bot/discord"

	zlg "github.com/mark-ignacio/zerolog-gcp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagMessageName            string
	flagMessageUID             string
	flagMessageUIDFile         string
	flagMessageMultiplePeriod  uint
	flagMessageUserRegOptional bool
)

// messageCmd represents the message command
var messageCmd = &cobra.Command{
	Use:   "message",
	Short: "Send templated messages for ad-hoc reasons, testing, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		var (
			token      = viper.GetString("token")
			gcpProject = viper.GetString("gcp-project")
		)
		if token == "" {
			log.Fatal().Msg("must specify a Discord token")
		}
		if gcpProject == "" {
			log.Fatal().Msg("must specify a GCP project ID")
		}
		gcpWriter, err := zlg.NewCloudLoggingWriter(ctx, gcpProject, "discord-bot-message", zlg.CloudLoggingOptions{})
		if err != nil {
			log.Panic().Err(err).Msg("could not create a CloudLoggingWriter")
		}
		log.Logger = log.Output(zerolog.MultiLevelWriter(
			zerolog.NewConsoleWriter(),
			gcpWriter,
		))
		fs, err := firestore.NewClient(ctx, gcpProject)
		if err != nil {
			log.Fatal().Err(err).Msg("error loading Firestore client")
		}
		msgBot, err := discord.NewMessagingBot(ctx, token, fs)
		if err != nil {
			log.Fatal().Err(err).Msg("error creating MessagingBot")
		}
		if flagMessageUID != "" {
			log.Info().Str("name", flagMessageName).Str("uid", flagMessageUID).Msg("messaging")
			err = msgBot.Message(flagMessageName, flagMessageUID, !flagMessageUserRegOptional)
			if err != nil {
				log.Fatal().Err(err).Msg("error sending message")
			}
		} else {
			var userIDs []string
			data, err := ioutil.ReadFile(flagMessageUIDFile)
			if err != nil {
				log.Fatal().Err(err).Msg("error opening UIDs file")
			}
			err = json.Unmarshal(data, &userIDs)
			if err != nil {
				log.Fatal().Err(err).Msg("error unmarshalling UIDs file")
			}
			ticker := time.NewTicker(time.Second * time.Duration(flagMessageMultiplePeriod))
			log.Info().Str("name", flagMessageName).Msg("messaging users")
			for _, uid := range userIDs {
				logger := log.With().Str("uid", uid).Logger()
				err = msgBot.Message(flagMessageName, flagMessageUID, !flagMessageUserRegOptional)
				if err != nil {
					logger.Fatal().Err(err).Msg("error sending message")
				} else {
					logger.Info().Msg("DM'd user")
				}
				<-ticker.C
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(messageCmd)
	flags := messageCmd.Flags()
	flags.StringVarP(&flagMessageName, "message", "m", "", "name of the message to send")
	flags.StringVar(&flagMessageUID, "uid", "", "Discord user ID for single messages")
	flags.StringVar(&flagMessageUIDFile, "uid-file", "", "UID file for sending many DMs")
	flags.UintVar(&flagMessageMultiplePeriod, "period", 5, "period to wait between DMs (in seconds)")
	flags.BoolVar(&flagMessageUserRegOptional, "unregistered", false, "allow sending to Gentei-unregistered users")
}
