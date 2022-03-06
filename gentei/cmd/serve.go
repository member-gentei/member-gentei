package cmd

import (
	"context"
	"os"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/member-gentei/member-gentei/gentei/web"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	envNameServeJWTKey = "JWT_KEY"
)

var (
	serveJWTKey      []byte
	flagServeAddress string
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the API server.",
	PreRun: func(cmd *cobra.Command, args []string) {
		serveJWTKey = []byte(os.Getenv(envNameServeJWTKey))
		if len(serveJWTKey) == 0 {
			log.Fatal().Msgf("env var %s must not be empty", envNameServeJWTKey)
		}
		if !strings.HasPrefix(flagServeAddress, "localhost") {
			if flagPubSubTopic == "" {
				log.Fatal().Msgf("env var %s must not be empty in prod", envNamePubSubTopic)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ctx           = context.Background()
			db            = mustOpenDB(ctx)
			discordConfig = getDiscordConfig()
			youTubeConfig = getYouTubeConfig()
			topic         *pubsub.Topic
		)
		ps, err := pubsub.NewClient(ctx, flagGCPProjectID)
		if err != nil {
			log.Fatal().Err(err).Msg("error calling pubsub.NewClient")
		}
		if flagPubSubTopic != "" {
			topic = ps.Topic(flagPubSubTopic)
			topic.PublishSettings.CountThreshold = 1
		}
		err = web.ServeAPI(db, discordConfig, youTubeConfig, topic, serveJWTKey, flagServeAddress, flagVerbose)
		if err != nil {
			log.Fatal().Err(err).Msg("server terminated")
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	flags := serveCmd.Flags()
	flags.StringVar(&flagServeAddress, "address", "localhost:5000", "API listening address")
}
