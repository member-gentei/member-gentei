package commands

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/member-gentei/member-gentei/gentei/ent/guildrole"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/member-gentei/member-gentei/gentei/ent/usermembership"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
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
	for i := range guildRoles {
		var (
			statusMessage string
			guildRole     = *guildRoles[i]
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

// GetDisabledChannelEmbeds creates a
func GetDisabledChannelEmbeds(ctx context.Context, db *ent.Client, channelIDs []string) ([]*discordgo.MessageEmbed, error) {
	talents, err := db.YouTubeTalent.Query().
		Where(youtubetalent.IDIn(channelIDs...)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting disabled channels by ID: %w", err)
	}
	embeds := make([]*discordgo.MessageEmbed, len(talents))
	for i := range talents {
		talent := *talents[i]
		embeds[i] = &discordgo.MessageEmbed{
			Type:  discordgo.EmbedTypeRich,
			Title: talent.ChannelName,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: talent.ThumbnailURL,
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Error",
					Value: "Membership checks are temporarily disabled for this YouTube channel. Please try again later!",
				},
			},
		}
	}
	return embeds, nil
}
