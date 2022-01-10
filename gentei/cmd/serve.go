package cmd

import (
	"context"
	"os"

	"github.com/member-gentei/member-gentei/gentei/web"
	discordoauth "github.com/ravener/discord-oauth2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	envNameServeJWTKey              = "JWT_KEY"
	envNameServeYouTubeClientID     = "YOUTUBE_CLIENT_ID"
	envNameServeYouTubeClientSecret = "YOUTUBE_CLIENT_SECRET"
)

var (
	serveJWTKey                  []byte
	flagServeAddress             string
	flagServeDiscordRedirectURL  string
	flagServeYouTubeClientID     string
	flagServeYouTubeClientSecret string
	flagServeYouTubeRedirectURL  string
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the API server.",
	PreRun: func(cmd *cobra.Command, args []string) {
		serveJWTKey = []byte(os.Getenv(envNameServeJWTKey))
		flagServeYouTubeClientID = os.Getenv(envNameServeYouTubeClientID)
		flagServeYouTubeClientSecret = os.Getenv(envNameServeYouTubeClientSecret)
		if len(serveJWTKey) == 0 {
			log.Fatal().Msgf("env var %s must not be empty", envNameServeJWTKey)
		}
		if flagDiscordClientID == "" {
			log.Fatal().Msgf("env var %s must not be empty", envNameDiscordClientID)
		}
		if flagDiscordClientSecret == "" {
			log.Fatal().Msgf("env var %s must not be empty", envNameDiscordClientSecret)
		}
		if flagServeYouTubeClientID == "" {
			log.Fatal().Msgf("env var %s must not be empty", envNameServeYouTubeClientID)
		}
		if flagServeYouTubeClientSecret == "" {
			log.Fatal().Msgf("env var %s must not be empty", envNameServeYouTubeClientSecret)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ctx           = context.Background()
			db            = mustOpenDB(ctx)
			discordConfig = &oauth2.Config{
				ClientID:     flagDiscordClientID,
				ClientSecret: flagDiscordClientSecret,
				Endpoint:     discordoauth.Endpoint,
				RedirectURL:  flagServeDiscordRedirectURL,
			}
			youTubeConfig = &oauth2.Config{
				ClientID:     flagServeYouTubeClientID,
				ClientSecret: flagServeYouTubeClientSecret,
				Endpoint:     google.Endpoint,
				RedirectURL:  flagServeYouTubeRedirectURL,
			}
		)
		err := web.ServeAPI(db, discordConfig, youTubeConfig, serveJWTKey, flagServeAddress, flagVerbose)
		if err != nil {
			log.Fatal().Err(err).Msg("server terminated")
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	flags := serveCmd.Flags()
	flags.StringVar(&flagServeAddress, "address", "localhost:5000", "API listening address")
	flags.StringVar(&flagServeDiscordRedirectURL, "discord-redirect-url", "http://localhost:3000/login/discord", "")
	flags.StringVar(&flagServeYouTubeRedirectURL, "youtube-redirect-url", "http://localhost:3000/login/youtube", "")
}
