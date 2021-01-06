package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	zlg "github.com/mark-ignacio/zerolog-gcp"
	"github.com/member-gentei/member-gentei/pkg/clients"
	"github.com/member-gentei/member-gentei/pkg/common"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile            string
	flagVerbose        bool
	flagDryRun         bool
	flagNoCloudLogging bool
	flagUID            string
	flagPubsubTopic    string
	flagStartAfter     string
)
var rootCmd = &cobra.Command{
	Use:   "gentei-member-check",
	Short: "Checks memberships for Gentei users",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ctx        = context.Background()
			gcpProject = viper.GetString("gcp-project")
			numWorkers = viper.GetUint("num-workers")
			psTopic    *pubsub.Topic
		)
		// set up logger
		if flagVerbose {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
		if flagNoCloudLogging {
			log.Logger = log.Output(zerolog.NewConsoleWriter())
			log.Info().Msg("Google Cloud Logging is disabled")
		} else {
			gcpWriter, err := zlg.NewCloudLoggingWriter(ctx, gcpProject, "member-check", zlg.CloudLoggingOptions{})
			if err != nil {
				log.Fatal().Err(err).Msg("could not create a CloudLoggingWriter")
			}
			defer zlg.Flush()
			log.Logger = log.Output(zerolog.MultiLevelWriter(
				zerolog.NewConsoleWriter(),
				gcpWriter,
			))
		}
		if flagPubsubTopic != "" {
			psClient, err := pubsub.NewClient(ctx, gcpProject)
			if err != nil {
				log.Fatal().Err(err).Msg("could not create Pub/Sub Client")
			}
			psTopic = psClient.Topic(flagPubsubTopic)
		}
		// start up Firestore
		fs, err := clients.NewRetryFirestoreClient(ctx, gcpProject)
		if err != nil {
			log.Fatal().Err(err).Msg("error creating Firestore client")
		}
		var uids []string
		if flagUID != "" {
			uids = []string{flagUID}
		}
		// perform the check!
		startTime := time.Now()
		results, err := common.EnforceMemberships(ctx, fs, &common.EnforceMembershipsOptions{
			ReloadDiscordGuilds:       true,
			RemoveInvalidDiscordToken: true,
			RemoveInvalidYouTubeToken: true,
			Apply:                     !flagDryRun,
			UserIDs:                   uids,
			StartAfter:                flagStartAfter,
			NumWorkers:                numWorkers,
		})
		if err != nil {
			log.Fatal().Err(err).Msg("error performing enforcement check")
		}
		endTime := time.Now()
		runtime := uint64(endTime.Sub(startTime).Seconds())
		if psTopic != nil {
			logger := log.With().Str("topic", flagPubsubTopic).Logger()
			result := psTopic.Publish(ctx, &pubsub.Message{
				Data: []byte(endTime.In(time.UTC).Format(time.RFC3339)),
			})
			if _, err = result.Get(ctx); err != nil {
				logger.Err(err).Msg("error publishing Pub/Sub message")
			} else {
				logger.Debug().Msg("published message to Pub/Sub topic")
			}
		}
		// don't log per-uid metrics to GCP - print them out!
		if uids == nil {
			log.Info().Interface("checkResults", results).Uint64("runtime", runtime).Msg("periodic check complete")
		} else {
			fmt.Printf("check results (%ds): %+v\n", runtime, results)
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
	persistent.BoolVarP(&flagDryRun, "dry-run", "n", false, "dry run mode")
	persistent.BoolVar(&flagNoCloudLogging, "no-cloud-logging", false, "do not output results to Google Cloud Logging")
	persistent.String("gcp-project", "member-gentei", "GCP project ID")
	persistent.StringVar(&flagUID, "uid", "", "specific user ID")
	persistent.StringVar(&flagPubsubTopic, "pubsub-topic", "", "pubsub topic to notify on completion")
	persistent.Uint("num-workers", 2, "number of worker threads")
	viper.BindPFlags(persistent)
	rootCmd.Flags().StringVar(&flagStartAfter, "start-after", "", "StartAfter argument")
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
