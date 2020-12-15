package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"

	zlg "github.com/mark-ignacio/zerolog-gcp"
	"github.com/member-gentei/member-gentei/bot/discord"
	"github.com/member-gentei/member-gentei/bot/discord/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	flagVerbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bot",
	Short: "A brief description of your application",
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		if flagVerbose {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
		log.Logger = log.Output(zerolog.NewConsoleWriter())
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		var (
			token      = viper.GetString("token")
			apiKey     = viper.GetString("api-key")
			apiServer  = viper.GetString("api-server")
			gcpProject = viper.GetString("gcp-project")
		)
		if token == "" {
			log.Fatal().Msg("must specify a Discord token")
		}
		if apiServer == "" {
			log.Fatal().Msg("must specify an API server")
		}
		if apiKey == "" {
			log.Fatal().Msg("must specify an API key")
		}
		if gcpProject == "" {
			log.Fatal().Msg("must specify a GCP project ID")
		}
		gcpWriter, err := zlg.NewCloudLoggingWriter(ctx, gcpProject, "discord-bot", zlg.CloudLoggingOptions{})
		if err != nil {
			log.Panic().Err(err).Msg("could not create a CloudLoggingWriter")
		}
		log.Logger = log.Output(zerolog.MultiLevelWriter(
			zerolog.NewConsoleWriter(),
			gcpWriter,
		))
		authHeader := fmt.Sprintf("Bearer %s", apiKey)
		apiClient, err := api.NewClientWithResponses(
			viper.GetString("api-server"),
			api.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
				req.Header.Set("Authorization", authHeader)
				return nil
			}),
		)
		if err != nil {
			log.Fatal().Err(err).Msg("error loading API client")
		}
		fs, err := firestore.NewClient(ctx, gcpProject)
		if err != nil {
			log.Fatal().Err(err).Msg("error loading Firestore client")
		}
		if err := discord.Start(ctx, token, apiClient, fs); err != nil {
			log.Fatal().Err(err).Msg("error running Discord bot")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	persistent := rootCmd.PersistentFlags()
	persistent.StringVar(&cfgFile, "config", "", "config file (default is .bot.yml)")
	persistent.BoolVarP(&flagVerbose, "verbose", "v", false, "DEBUG level logging")
	persistent.String("token", "", "Discord bot token")
	persistent.String("api-server", "https://us-central1-member-gentei.cloudfunctions.net/API", "API URL")
	viper.BindPFlags(persistent)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName(".bot")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
