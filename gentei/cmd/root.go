package cmd

import (
	"context"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/migrate"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	envNameDiscordClientID     = "DISCORD_CLIENT_ID"
	envNameDiscordClientSecret = "DISCORD_CLIENT_SECRET"
)

var (
	flagDBEngine string
	flagOpenDB   string
	flagVerbose  bool

	flagDiscordClientID     string
	flagDiscordClientSecret string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gentei",
	Short: "Everything in member-gentei that isn't the frontend.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if flagVerbose {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		flagOpenDB = viper.GetString("DB")
		flagDBEngine = viper.GetString("engine")
		flagDiscordClientID = os.Getenv(envNameDiscordClientID)
		flagDiscordClientSecret = os.Getenv(envNameDiscordClientSecret)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
func mustOpenDB(ctx context.Context) *ent.Client {
	logger := log.With().Str("engine", flagDBEngine).Logger()
	db, err := ent.Open(flagDBEngine, flagOpenDB)
	if err != nil {
		logger.Fatal().Err(err).Msg("error opening SQL database")
	}
	if err := db.Schema.Create(
		ctx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		logger.Fatal().Err(err).Msg("failed to create schema resources")
	}
	return db
}

func init() {
	persistent := rootCmd.PersistentFlags()
	persistent.BoolVarP(&flagVerbose, "verbose", "v", false, "debug/verbose logging")
	persistent.String("engine", "sqlite3", "one of: sqlite3, pgx")
	persistent.String("db", "file:ent.sqlite3?cache=shared&_fk=1", "sql connection string")
	viper.BindPFlags(persistent)
	viper.EnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}
