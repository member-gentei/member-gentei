package cmd

import (
	"context"
	"encoding/json"
	"os"
	"slices"

	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore application-level backups",
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ctx    = context.Background()
			db     = mustOpenDB(ctx)
			backup backupFile
		)
		f, err := os.Open(args[0])
		if err != nil {
			log.Fatal().Str("name", args[0]).Err(err).Msg("error opening backup file")
		}
		defer f.Close()
		if err := json.NewDecoder(f).Decode(&backup); err != nil {
			log.Fatal().Str("name", args[0]).Err(err).Msg("error unmarshalling backup file")
		}
		// do all of it in one transaction (lol)
		// First, restore guilds + talents + users.
		// Second, restore GuildRoles + UserMemberships
		tx, err := db.Tx(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("error creating transaction")
		}
		gcBuilders := make([]*ent.GuildCreate, len(backup.GuildsWithRoles))
		for i, g := range backup.GuildsWithRoles {
			gc := tx.Guild.Create().
				SetID(g.ID).
				SetName(g.Name).
				SetSettings(g.Settings)
			// nillable stuff
			if g.IconHash != "" {
				gc.SetIconHash(g.IconHash)
			}
			if g.AuditChannel != 0 {
				gc.SetAuditChannel(g.AuditChannel)
			}
			if g.Language != "" {
				gc.SetLanguage(g.Language)
			}
			gcBuilders[i] = gc
		}
		if err := tx.Guild.CreateBulk(gcBuilders...).Exec(ctx); err != nil {
			tx.Rollback()
			log.Fatal().Err(err).Msg("error restoring guilds")
		}
		gcBuilders = nil
		tcBuilders := make([]*ent.YouTubeTalentCreate, len(backup.TalentsWithRoles))
		for i, t := range backup.TalentsWithRoles {
			tc := tx.YouTubeTalent.Create().
				SetID(t.ID).
				SetChannelName(t.ChannelName).
				SetThumbnailURL(t.ThumbnailURL)
			// nillable
			if t.MembershipVideoID != "" {
				tc.SetMembershipVideoID(t.MembershipVideoID)
			}
			if !t.LastMembershipVideoIDMiss.IsZero() {
				tc.SetLastMembershipVideoIDMiss(t.LastMembershipVideoIDMiss)
			}
			if !t.LastUpdated.IsZero() {
				tc.SetLastUpdated(t.LastUpdated)
			}
			if !t.Disabled.IsZero() {
				tc.SetDisabled(t.Disabled)
			}
			if t.DisabledPermanently {
				tc.SetDisabledPermanently(true)
			}
			tcBuilders[i] = tc
		}
		if err := tx.YouTubeTalent.CreateBulk(tcBuilders...).Exec(ctx); err != nil {
			tx.Rollback()
			log.Fatal().Err(err).Msg("error restoring YouTubeTalents")
		}
		tcBuilders = nil
		ucBuilders := make([]*ent.UserCreate, len(backup.UsersWithMemberships))
		for i, u := range backup.UsersWithMemberships {
			uc := tx.User.Create().
				SetID(u.ID).
				SetFullName(u.FullName).
				SetAvatarHash(u.AvatarHash).
				SetDiscordToken(u.DiscordToken).
				SetYoutubeToken(u.YoutubeToken)
			// nillable
			if !u.LastCheck.IsZero() {
				uc.SetLastCheck(u.LastCheck)
			}
			uc.SetNillableYoutubeID(u.YoutubeID)
			ucBuilders[i] = uc
		}
		// 1000 at a time because SQL variable limit lol
		for chunk := range slices.Chunk(ucBuilders, 1000) {
			if err := tx.User.CreateBulk(chunk...).Exec(ctx); err != nil {
				tx.Rollback()
				log.Fatal().Err(err).Msg("error restoring User chunk")
			}
		}
		ucBuilders = nil
		// for GuildRoles, we need to find the corresponding talent IDs
		guildRoleIDToTalentID := make(map[uint64]string, len(backup.TalentsWithRoles))
		for _, t := range backup.TalentsWithRoles {
			for _, role := range t.Edges.Roles {
				guildRoleIDToTalentID[role.ID] = t.ID
			}
		}
		grBuilders := make([]*ent.GuildRoleCreate, 0, len(backup.GuildsWithRoles))
		for _, g := range backup.GuildsWithRoles {
			for _, gr := range g.Edges.Roles {
				talentID, ok := guildRoleIDToTalentID[gr.ID]
				if !ok {
					log.Fatal().Msg("could not find talent ID for GuildRole")
				}
				grc := tx.GuildRole.Create().
					SetID(gr.ID).
					SetName(gr.Name).
					SetLastUpdated(gr.LastUpdated).
					SetGuildID(g.ID).
					SetTalentID(talentID)
				grBuilders = append(grBuilders, grc)
			}
		}
		if err := tx.GuildRole.CreateBulk(grBuilders...).Exec(ctx); err != nil {
			tx.Rollback()
			log.Fatal().Err(err).Msg("error restoring GuildRoles")
		}
		grBuilders = nil
		umBuilders := make([]*ent.UserMembershipCreate, 0, len(backup.UsersWithMemberships)*3)
		for _, u := range backup.UsersWithMemberships {
			for _, um := range u.Edges.Memberships {
				umc := tx.UserMembership.Create().
					SetUserID(u.ID).
					SetYoutubeTalentID(um.Edges.YoutubeTalent.ID).
					SetLastVerified(um.LastVerified).
					SetFailCount(um.FailCount)
				if !um.FirstFailed.IsZero() {
					umc.SetFirstFailed(um.FirstFailed)
				}
				for _, r := range um.Edges.Roles {
					umc.AddRoleIDs(r.ID)
				}
				umBuilders = append(umBuilders, umc)
			}
		}
		if len(umBuilders) == 0 {
			log.Fatal().Msg("no UserMembership builders...?")
		}
		// 1000 at a time because SQL variable limit lol
		for chunk := range slices.Chunk(umBuilders, 1000) {
			if err := tx.UserMembership.CreateBulk(chunk...).Exec(ctx); err != nil {
				tx.Rollback()
				log.Fatal().Err(err).Msg("error restoring UserMemberships chunk")
			}
		}
		if err := tx.Commit(); err != nil {
			log.Fatal().Err(err).Msg("error committing tx")
		}
		log.Info().Msg("restored backup")
	},
}

func init() {
	adminCmd.AddCommand(restoreCmd)
}
