package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/member-gentei/member-gentei/gentei/ent/guildrole"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/member-gentei/member-gentei/gentei/ent/usermembership"
	"github.com/member-gentei/member-gentei/gentei/membership"
)

// CreateMembershipInfoEmbeds creates a []*discordgo.MessageEmbed that describes the user's current role grants on a server.
func CreateMembershipInfoEmbeds(ctx context.Context, db *ent.Client, userID, guildID uint64) ([]*discordgo.MessageEmbed, error) {
	guildRoles, err := db.GuildRole.Query().
		WithTalent().
		Where(
			guildrole.HasGuildWith(guild.ID(guildID)),
		).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("error querying for all roles for guild: %w", err)
	}
	applicableGuildRoleIDs := map[uint64]bool{}
	applicableGuildRoleIDSlice, err := db.GuildRole.Query().
		Where(
			guildrole.HasGuildWith(guild.ID(guildID)),
			guildrole.HasUserMembershipsWith(
				usermembership.HasUserWith(user.ID(userID)),
			),
		).
		IDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("error querying for applicable roles for guild: %w", err)
	}
	for _, guildRoleID := range applicableGuildRoleIDSlice {
		applicableGuildRoleIDs[guildRoleID] = true
	}
	embeds := make([]*discordgo.MessageEmbed, len(guildRoles))
	for i, guildRole := range guildRoles {
		var (
			statusMessage string
			talent        = guildRole.Edges.Talent
		)
		if applicableGuildRoleIDs[guildRole.ID] {
			statusMessage = "✅ Membership verified"
		} else {
			statusMessage = "⛔ Not a member"
		}
		embeds[i] = &discordgo.MessageEmbed{
			Type:  discordgo.EmbedTypeRich,
			Title: talent.ChannelName,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: talent.ThumbnailURL,
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Membership Role",
					Value:  fmt.Sprintf("<@&%d>", guildRole.ID),
					Inline: true,
				},
				{
					Name:   "Status",
					Value:  statusMessage,
					Inline: true,
				},
			},
		}
	}
	return embeds, nil
}

// CreateLGRWebhookEdit creates a webhook edit with embeds of lost/gained/retained membership roles.
func CreateCheckResultWebhookEdit(
	results *membership.CheckResultSet,
	roleMap map[string]string,
	talentMap map[string]*ent.YouTubeTalent,
	asOf time.Time,
) *discordgo.WebhookEdit {
	var embeds []*discordgo.MessageEmbed
	for _, lost := range results.Lost {
		embeds = append(embeds, checkResultAsEmbed(lost, "⚠️ Verification failed", talentMap))
	}
	for _, lost := range results.Gained {
		embeds = append(embeds, checkResultAsEmbed(lost, "✅ Membership verified", talentMap))
	}
	for _, lost := range results.Retained {
		embeds = append(embeds, checkResultAsEmbed(lost, "✅ Membership verified", talentMap))
	}
	for _, lost := range results.Not {
		embeds = append(embeds, checkResultAsEmbed(lost, "⛔ Not a member", talentMap))
	}
	edit := &discordgo.WebhookEdit{
		Content: "Your membership check results are below.",
		Embeds:  embeds,
	}
	return edit
}

func checkResultAsEmbed(result membership.CheckResult, statusValue string, talentMap map[string]*ent.YouTubeTalent) *discordgo.MessageEmbed {
	var (
		talent = talentMap[result.ChannelID]
	)
	return &discordgo.MessageEmbed{
		Type:  discordgo.EmbedTypeRich,
		Title: talent.ChannelName,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: talent.ThumbnailURL,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Membership Role",
				Value:  fmt.Sprintf("<@&%s>", result.ChannelID),
				Inline: true,
			},
			{
				Name:   "Status",
				Value:  statusValue,
				Inline: true,
			},
			{
				Name:   "Check Time",
				Value:  fmt.Sprintf("<t:%d:R>", result.Time.Unix()),
				Inline: true,
			},
		},
	}
}
