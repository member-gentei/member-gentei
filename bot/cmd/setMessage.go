package cmd

import (
	"context"
	"io/ioutil"

	"github.com/member-gentei/member-gentei/pkg/common"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagSetMessageName string
	flagSetMessageBody string
)

// setMessageCmd represents the setMessage command
var setMessageCmd = &cobra.Command{
	Use:   "setMessage",
	Short: "Sets a message (because the Firestore UI doesn't support newlines)",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			gcpProject = viper.GetString("gcp-project")
			ctx        = context.Background()
		)
		log.Logger = log.Output(zerolog.NewConsoleWriter())
		if gcpProject == "" {
			log.Fatal().Msg("must specify a GCP project ID")
		}
		if flagSetMessageName == "" {
			log.Fatal().Msg("must specify a message name")
		}
		if flagSetMessageBody == "" {
			log.Fatal().Msg("must specify a message body file")
		}
		body, err := ioutil.ReadFile(flagSetMessageBody)
		if err != nil {
			log.Fatal().Err(err).Msg("error reading message body file")
		}
		fs, err := firestore.NewClient(ctx, gcpProject)
		if err != nil {
			log.Fatal().Err(err).Msg("error loading Firestore client")
		}
		_, err = fs.Collection(common.DMTemplateCollection).Doc(flagSetMessageName).
			Set(ctx, common.DMTemplate{
				Name: flagSetMessageName,
				Body: string(body),
			})
		if err != nil {
			log.Fatal().Err(err).Msg("error setting DMTemplate")
		}
	},
}

func init() {
	rootCmd.AddCommand(setMessageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setMessageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	flags := setMessageCmd.Flags()
	flags.StringVarP(&flagSetMessageName, "message", "m", "", "name of the message to set")
	flags.StringVarP(&flagSetMessageBody, "file", "f", "", "message body file")
}
