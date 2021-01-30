package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"cloud.google.com/go/pubsub"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
	"google.golang.org/protobuf/types/known/timestamppb"

	zlg "github.com/mark-ignacio/zerolog-gcp"
	"github.com/member-gentei/member-gentei/bot/discord"
	"github.com/member-gentei/member-gentei/bot/discord/api"
	"github.com/member-gentei/member-gentei/pkg/clients"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile         string
	flagVerbose     bool
	flagNoHeartbeat bool
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
			token           = viper.GetString("token")
			apiKey          = viper.GetString("api-key")
			apiServer       = viper.GetString("api-server")
			gcpProject      = viper.GetString("gcp-project")
			membershipSubID = viper.GetString("membership-sub-id")
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
		if membershipSubID == "" {
			log.Fatal().Msg("must specify a pub/sub subscription ID")
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
		fs, err := clients.NewRetryFirestoreClient(ctx, gcpProject)
		if err != nil {
			log.Fatal().Err(err).Msg("error loading Firestore client")
		}
		psClient, err := pubsub.NewClient(ctx, gcpProject)
		if err != nil {
			log.Fatal().Err(err).Msg("error loading Pub/Sub client")
		}
		psSubscription := psClient.Subscription(membershipSubID)
		opts := &discord.StartOptions{
			Token:                        token,
			APIClient:                    apiClient,
			FirestoreClient:              fs,
			MembershipReloadSubscription: psSubscription,
			HeartbeatCallback:            makeHeartbeatCallback(ctx, gcpProject),
			Heartbeat:                    !flagNoHeartbeat,
		}
		if err := discord.Start(ctx, opts); err != nil {
			log.Fatal().Err(err).Msg("error running Discord bot")
		}
	},
}

func makeHeartbeatCallback(ctx context.Context, projectID string) func() {
	client, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating monitoring metric client")
	}
	hostname, _ := os.Hostname()
	resource := &monitoredrespb.MonitoredResource{
		Type:   "global",
		Labels: map[string]string{"project_id": projectID},
	}
	timeSeriesMetric := &metricpb.Metric{
		Type:   "custom.googleapis.com/gentei/heartbeat",
		Labels: map[string]string{"host": hostname},
	}
	timeSeriesValue := &monitoringpb.TypedValue{
		Value: &monitoringpb.TypedValue_Int64Value{Int64Value: 1},
	}
	return func() {
		err := client.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
			Name: monitoring.MetricProjectPath(projectID),
			TimeSeries: []*monitoringpb.TimeSeries{
				{
					Metric:   timeSeriesMetric,
					Resource: resource,
					Points: []*monitoringpb.Point{
						{
							Interval: &monitoringpb.TimeInterval{EndTime: timestamppb.Now()},
							Value:    timeSeriesValue,
						},
					},
				},
			},
		})
		if err != nil {
			log.Err(err).Msg("error sending heartbeat gauge metric")
		}
	}
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
	persistent.String("membership-sub-id", "", "Pub/Sub subscription ID for membership list reloads")
	persistent.BoolVar(&flagNoHeartbeat, "no-heartbeat", false, "do not emit heartbeat metrics to GCP Monitoring")
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
