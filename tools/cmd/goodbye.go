package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/member-gentei/member-gentei/pkg/clients"
	"github.com/member-gentei/member-gentei/pkg/common"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// goodbyeCmd represents the goodbye command
var goodbyeCmd = &cobra.Command{
	Use:   "goodbye",
	Short: "Revoke all refresh tokens!",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ctx = context.Background()
		)
		fs, err := clients.NewRetryFirestoreClient(ctx, flagProjectID)
		if err != nil {
			log.Fatal().Err(err).Msg("error creating retry client")
		}
		docs, err := fs.CollectionGroup(common.PrivateCollection).Documents(ctx).GetAll()
		if err != nil {
			log.Fatal().Err(err).Msg("error getting all Private docs")
		}
		for _, doc := range docs {
			var (
				logger = log.With().Str("userID", doc.Ref.Parent.Parent.ID).Logger()
				token  oauth2.Token
			)
			err = doc.DataTo(&token)
			if err != nil {
				logger.Fatal().Err(err).Str("id", doc.Ref.ID).Msg("error unmarshalling token")
			}
			switch doc.Ref.ID {
			case "youtube":
				err = revokeYoutubeToken(token.RefreshToken, logger)
				if err != nil {
					log.Err(err).Msg("error revoking YouTube token, persisting doc")
				} else {
					_, err = doc.Ref.Delete(ctx)
					if err != nil {
						log.Err(err).Msg("error deleting YouTube token doc")
					}
					logger.Info().Msg("revoked YouTube token")
				}

			case "discord":
				err = revokeDiscordToken(token.RefreshToken, logger)
				if err != nil {
					log.Err(err).Msg("error revoking Discord token, persisting doc")
				} else {
					_, err = doc.Ref.Delete(ctx)
					if err != nil {
						log.Err(err).Msg("error deleting Discord token doc")
					}
					logger.Info().Msg("revoked Discord token")
				}
			default:
				panic(doc.Ref.ID)
			}
		}
	},
}

func revokeYoutubeToken(refreshToken string, logger zerolog.Logger) error {
	r, err := http.Post(
		fmt.Sprintf("https://oauth2.googleapis.com/revoke?token=%s", refreshToken),
		"application/x-www-form-urlencoded",
		nil,
	)
	if err != nil {
		logger.Warn().Err(err).Msg("error revoking YouTube token")
		return err
	}
	if r.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(r.Body)
		logger.Warn().Bytes("body", body).Int("statusCode", r.StatusCode).Msg("non-200 response while revoking YouTube token")
	}
	return nil
}

func revokeDiscordToken(refreshToken string, logger zerolog.Logger) error {
	values := url.Values{}
	values.Set("token", refreshToken)
	values.Set("token_type_hint", "refresh_token")
	r, err := http.PostForm("https://discord.com/api/oauth2/token/revoke", values)
	if err != nil {
		logger.Warn().Err(err).Msg("error revoking Discord token")
		return err
	}
	if r.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(r.Body)
		logger.Warn().Bytes("body", body).Int("statusCode", r.StatusCode).Msg("non-200 response while revoking YouTube token")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(goodbyeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// goodbyeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// goodbyeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
