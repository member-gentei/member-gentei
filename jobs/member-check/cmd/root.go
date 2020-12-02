package cmd

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/member-gentei/member-gentei/pkg/common"
	zlg "github.com/mark-ignacio/zerolog-gcp"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	flagVerbose bool
)
var rootCmd = &cobra.Command{
	Use:   "gentei-member-check",
	Short: "Checks memberships for Gentei users",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ctx        = context.Background()
			gcpProject = viper.GetString("gcp-project")
		)
		// set up logger
		if flagVerbose {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
		gcpWriter, err := zlg.NewCloudLoggingWriter(ctx, gcpProject, "member-check", zlg.CloudLoggingOptions{})
		if err != nil {
			log.Fatal().Err(err).Msg("could not create a CloudLoggingWriter")
		}
		defer zlg.Flush()
		log.Logger = log.Output(zerolog.MultiLevelWriter(
			zerolog.NewConsoleWriter(),
			gcpWriter,
		))
		// start up Firestore
		fs, err := firestore.NewClient(ctx, gcpProject)
		if err != nil {
			log.Fatal().Err(err).Msg("error creating Firestore client")
		}
		// perform the check!
		err = common.EnforceMemberships(ctx, fs, &common.EnforceMembershipsOptions{
			ReloadDiscordGuilds: true,
			RemoveInvalidTokens: true,
			Apply:               true,
		})
		if err != nil {
			log.Fatal().Err(err).Msg("error performing enforcement check")
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
	persistent.StringVar(&cfgFile, "config", "", "config file (default is $HOME/.member-check.yaml)")
	persistent.BoolVarP(&flagVerbose, "verbose", "v", false, "DEBUG level logging")
	persistent.String("gcp-project", "member-gentei", "GCP project ID")
	viper.BindPFlags(persistent)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".member-check" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".member-check")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
