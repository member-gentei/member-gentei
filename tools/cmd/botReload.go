package cmd

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	flagBotReloadTopic string
)

// botReloadCmd represents the botReload command
var botReloadCmd = &cobra.Command{
	Use:   "botReload",
	Short: "Tells the bot to reload messages by publishing a message.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		if flagBotReloadTopic == "" {
			log.Fatal().Msg("must specify --pubsub-topic")
		}
		psClient, err := pubsub.NewClient(ctx, flagProjectID)
		if err != nil {
			log.Fatal().Err(err).Msg("could not create Pub/Sub Client")
		}
		psTopic := psClient.Topic(flagBotReloadTopic)
		messageID, err := psTopic.Publish(ctx, &pubsub.Message{
			Data: []byte(time.Now().In(time.UTC).Format(time.RFC3339)),
		}).Get(ctx)
		if err != nil {
			log.Err(err).Msg("error publishing Pub/Sub message")
		}
		log.Info().Msgf("message ID: %s", messageID)
	},
}

func init() {
	rootCmd.AddCommand(botReloadCmd)

	flags := botReloadCmd.Flags()
	flags.StringVar(&flagBotReloadTopic, "pubsub-topic", "", "Pub/Sub topic")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// botReloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// botReloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
