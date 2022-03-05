package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/member-gentei/member-gentei/gentei/async"
	"github.com/member-gentei/member-gentei/gentei/membership"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// asyncCmd represents the async command
var asyncCmd = &cobra.Command{
	Use:   "async",
	Short: "Works with the async task queue.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if flagDiscordClientID == "" {
			log.Fatal().Msgf("env var %s must not be empty", envNameDiscordClientID)
		}
		if flagDiscordClientSecret == "" {
			log.Fatal().Msgf("env var %s must not be empty", envNameDiscordClientSecret)
		}
		if flagYouTubeClientID == "" {
			log.Fatal().Msgf("env var %s must not be empty", envNameYouTubeClientID)
		}
		if flagYouTubeClientSecret == "" {
			log.Fatal().Msgf("env var %s must not be empty", envNameYouTubeClientSecret)
		}
		if flagPubSubSubscription == "" {
			log.Fatal().Msgf("env var %s must not be empty", envNamePubSubSubscription)
		}
		if flagPubSubTopic == "" {
			log.Fatal().Msgf("env var %s must not be empty", envNamePubSubTopic)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ctx, cancel = context.WithCancel(context.Background())
			db          = mustOpenDB(ctx)
			wg          = &sync.WaitGroup{}
		)
		ps, err := pubsub.NewClient(ctx, flagGCPProjectID)
		if err != nil {
			log.Fatal().Err(err).Msg("error calling pubsub.NewClient")
		}
		botTopic := ps.Topic(flagPubSubTopic)
		defer botTopic.Flush()
		changeHandler := async.NewPubSubMembershipChangeHandler(ctx, botTopic)
		membership.HookMembershipChanges(db, changeHandler)
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Info().Str("subscription", flagPubSubSubscription).Msg("listening to general pubsub subscription")
			err = async.ListenGeneral(ctx, db, getYouTubeConfig(), ps.Subscription(flagPubSubSubscription), botTopic, changeHandler.SetChangeReason)
			if err != nil {
				log.Fatal().Err(err).Msg("error calling async.ListenGeneral")
			}
		}()
		// graceful sigint handling
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)
		<-stop
		log.Info().Msg("shutting down")
		cancel()
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(asyncCmd)
}
