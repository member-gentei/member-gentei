package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	zlg "github.com/mark-ignacio/zerolog-gcp"
	"github.com/member-gentei/member-gentei/pkg/clients"
	"github.com/member-gentei/member-gentei/pkg/common"
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
	Use:   "refresh-data",
	Short: "Refreshes 'friendly' data like names and thumbnails.",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ctx        = context.Background()
			gcpProject = viper.GetString("gcp-project")
			botToken   = viper.GetString("token")
		)
		// set up logger
		if flagVerbose {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
		if botToken == "" {
			log.Fatal().Msg("must specify a --token")
		}
		dg, err := discordgo.New("Bot " + botToken)
		if err != nil {
			log.Fatal().Err(err).Msg("error connecting bot")
		}
		// set up logging
		if flagVerbose {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
		gcpWriter, err := zlg.NewCloudLoggingWriter(ctx, gcpProject, "refresh-data", zlg.CloudLoggingOptions{})
		if err != nil {
			log.Fatal().Err(err).Msg("could not create a CloudLoggingWriter")
		}
		defer zlg.Flush()
		log.Logger = log.Output(zerolog.MultiLevelWriter(
			zerolog.NewConsoleWriter(),
			gcpWriter,
		))
		// get clients to things
		fs, err := clients.NewRetryFirestoreClient(ctx, gcpProject)
		if err != nil {
			log.Fatal().Err(err).Msg("error creating Firestore client")
		}
		yt, err := common.GetYoutubeServerService(ctx, fs)
		if err != nil {
			log.Fatal().Err(err).Msg("error creating YouTube service")
		}
		err = common.RefreshYouTubeChannels(ctx, fs, yt)
		if err != nil {
			log.Fatal().Err(err).Msg("error refreshing YouTube channel data")
		}
		err = common.RefreshDiscordGuilds(ctx, fs, dg)
		if err != nil {
			log.Fatal().Err(err).Msg("error refreshing Discord guild data")
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
	persistent.StringVar(&cfgFile, "config", "", "config file (default is $HOME/.refresh-data.yaml)")
	persistent.BoolVarP(&flagVerbose, "verbose", "v", false, "DEBUG level logging")
	persistent.String("gcp-project", "member-gentei", "GCP project ID")
	persistent.String("token", "", "Discord bot token")
	viper.BindPFlags(persistent)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
