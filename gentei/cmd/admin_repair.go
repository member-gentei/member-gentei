package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/member-gentei/member-gentei/gentei/apis"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/member-gentei/member-gentei/gentei/ent/guildrole"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/member-gentei/member-gentei/gentei/ent/usermembership"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
	"github.com/member-gentei/member-gentei/gentei/web"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	flagRepairAll            bool
	flagRepairChannelID      string
	flagRepairFallbackUserID uint64
)

// repairCmd represents the repair command
var repairCmd = &cobra.Command{
	Use:       "repair",
	Short:     "Maintenance commands for various fun events",
	ValidArgs: []string{"unwind-multirole"},
	Args:      cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ctx = context.Background()
			db  = mustOpenDB(ctx)
		)
		if len(args) == 1 {
			switch args[0] {
			case "unwind-multirole":
				repairUnwindMultirole(ctx, db)
			}
			return
		}
		if flagRepairChannelID != "" {
			if err := repairChannelID(ctx, db, flagRepairChannelID); err != nil {
				log.Fatal().Err(err).Msg("error repairing channel")
			}
			return
		}
		if flagRepairAll {
			disabledChannelIDs := db.YouTubeTalent.Query().Where(youtubetalent.DisabledNotNil()).IDsX(ctx)
			for _, channelID := range disabledChannelIDs {
				if err := repairChannelID(ctx, db, channelID); err != nil {
					log.Err(err).Str("channelID", channelID).Msg("error repairing channel")
				}
			}
			session, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv(envNameDiscordBotToken)))
			if err != nil {
				log.Fatal().Err(err).Msg("error creating discordgo client")
			}
			err = repairGuilds(ctx, db, session)
			if err != nil {
				log.Fatal().Err(err).Msg("error repairing Guild set")
			}
			if err = refreshYouTubeChannels(ctx, db); err != nil {
				log.Fatal().Err(err).Msg("error refreshing YouTube channel info")
			}
			return
		}
		log.Fatal().Msg("please specify --channel-id or --all")
	},
}

func repairChannelID(ctx context.Context, db *ent.Client, channelID string) error {
	var (
		logger = log.With().Str("channelID", channelID).Logger()
		talent = db.YouTubeTalent.GetX(ctx, channelID)
	)
	// get a YouTube service from a member
	userID := db.User.Query().
		Where(
			user.YoutubeTokenNotNil(),
			user.HasMembershipsWith(
				usermembership.HasYoutubeTalentWith(youtubetalent.ID(channelID)),
			),
		).
		FirstIDX(ctx)
	if userID == 0 {
		logger.Warn().
			Str("userID", strconv.FormatUint(flagRepairFallbackUserID, 10)).
			Msg("no previously-verified members, using fallback user ID")
		userID = flagRepairFallbackUserID
	}
	svc, err := apis.GetYouTubeService(ctx, db, userID, getYouTubeConfig())
	if err != nil {
		return fmt.Errorf("error getting YouTube service for user %d: %w", userID, err)
	}
	var getNewMemberVideo bool
	if talent.MembershipVideoID == "" {
		logger.Info().Msg("membership video ID is blank, getting one...")
		getNewMemberVideo = true
	} else if !talent.Disabled.IsZero() {
		logger.Info().Msg("membership checks disabled, getting new video ID...")
		getNewMemberVideo = true
	}
	if getNewMemberVideo {
		videoID, err := apis.SelectRandomMembersOnlyVideoID(ctx, logger, svc, channelID)
		if err != nil {
			return fmt.Errorf("error selecting membership video ID: %w", err)
		}
		talent = db.YouTubeTalent.UpdateOne(talent).
			SetMembershipVideoID(videoID).
			ClearDisabled().
			SaveX(ctx)
		logger.Info().
			Str("channelID", talent.ID).
			Str("channelName", talent.ChannelName).
			Str("videoID", videoID).Msg("set new MembershipVideoID")
	}
	return nil
}

func repairGuilds(ctx context.Context, db *ent.Client, session *discordgo.Session) error {
	var (
		guildIDSlice   = db.Guild.Query().IDsX(ctx)
		storedGuildIDs = make(map[string]bool, len(guildIDSlice))
	)
	for _, gid := range guildIDSlice {
		storedGuildIDs[strconv.FormatUint(gid, 10)] = true
	}
	var afterID string
	for {
		guilds, err := session.UserGuilds(100, "", afterID)
		if err != nil {
			return err
		}
		for _, dg := range guilds {
			delete(storedGuildIDs, dg.ID)
		}
		if len(guilds) < 100 {
			break
		}
		afterID = guilds[len(guilds)-1].ID
	}
	if len(storedGuildIDs) > 0 {
		log.Info().Int("count", len(storedGuildIDs)).Msg("removing guilds from database")
	}
	for gid := range storedGuildIDs {
		guildID, err := strconv.ParseUint(gid, 10, 64)
		if err != nil {
			return err
		}
		err = db.Guild.DeleteOneID(guildID).Exec(ctx)
		if err != nil {
			return err
		}
		log.Info().Str("guildID", gid).Msg("deleted Guild")
	}
	return nil
}

func refreshYouTubeChannels(ctx context.Context, db *ent.Client) error {
	toUpdate, err := db.YouTubeTalent.Query().
		Where(youtubetalent.LastUpdatedLT(time.Now().Add(-24 * time.Hour))).
		IDs(ctx)
	if err != nil {
		return fmt.Errorf("error getting stale YouTube channels: %w", err)
	}
	for _, channelID := range toUpdate {
		err = web.UpsertYouTubeChannelID(ctx, db, channelID)
		if err != nil {
			return fmt.Errorf("error upserting channel info: %w", err)
		}
	}
	return nil
}

func repairUnwindMultirole(ctx context.Context, db *ent.Client) {
	flagBotToken = os.Getenv(envNameDiscordBotToken)
	if flagBotToken == "" {
		log.Fatal().Msgf("must specify env var '%s'", envNameDiscordBotToken)
	}
	log.Info().Msg("unwinding multirole issue")
	var vs []struct {
		Guild  uint64 `json:"guild_roles"`
		Talent string `json:"you_tube_talent_roles"`
		Count  int
	}
	db.Debug().GuildRole.Query().
		GroupBy(guildrole.GuildColumn, guildrole.TalentColumn).
		Aggregate(ent.Count()).
		ScanX(ctx, &vs)
	session := must(discordgo.New(fmt.Sprintf("Bot %s", flagBotToken)))
	for _, v := range vs {
		if v.Count == 1 {
			continue
		}
		log.Info().Interface("v", v).Msg("rectifying")
		grids := db.GuildRole.Query().
			Where(
				guildrole.HasGuildWith(guild.ID(v.Guild)),
				guildrole.HasTalentWith(youtubetalent.ID(v.Talent)),
			).
			Order(ent.Asc(guildrole.FieldLastUpdated)).
			IDsX(ctx)
		// delete roles that no longer exist
		roles, err := session.GuildRoles(strconv.FormatUint(v.Guild, 10))
		if err != nil {
			log.Fatal().Err(err).Msg("error listing Guild roles")
		}
		serverRoleIDs := make(map[uint64]bool, len(roles))
		for _, role := range roles {
			roleID := must(strconv.ParseUint(role.ID, 10, 64))
			serverRoleIDs[roleID] = true
		}
		var deleteCount int
		for _, grid := range grids {
			if !serverRoleIDs[grid] {
				log.Info().Uint64("grid", grid).Msg("deleting now missing GuildRole")
				db.GuildRole.DeleteOneID(grid).ExecX(ctx)
				deleteCount++
			}
		}
		if deleteCount == len(grids)-1 {
			log.Info().Msg("rectified by sheer luck")
		}
	}
}

func init() {
	adminCmd.AddCommand(repairCmd)
	flags := repairCmd.Flags()
	flags.StringVarP(&flagRepairChannelID, "channel-id", "c", "", "YouTube channel ID to repair")
	flags.BoolVar(&flagRepairAll, "all", false, "perform all repair actions")
	// the default is my Discord user ID...
	flags.Uint64VarP(&flagRepairFallbackUserID, "fallback-user-id", "f", 196078350496825345, "Fallback user ID for membership video filling")
}
