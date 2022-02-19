package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/member-gentei/member-gentei/gentei/ent"
)

func GetGuildInfoEmbeds(dg *ent.Guild, adminView bool) []*discordgo.MessageEmbed {
	var (
		embeds          []*discordgo.MessageEmbed
		rolesByTalentID = map[string]*ent.GuildRole{}
	)
	for _, role := range dg.Edges.Roles {
		rolesByTalentID[role.Edges.Talent.ID] = role
	}
	for _, talent := range dg.Edges.YoutubeTalents {
		embed := &discordgo.MessageEmbed{
			Type:  discordgo.EmbedTypeRich,
			Title: talent.ChannelName,
			URL:   fmt.Sprintf("https://www.youtube.com/channel/%s", talent.ID),
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: talent.ThumbnailURL,
			},
		}
		var membershipRoleValue string
		if dg.Settings != nil {
			role, found := rolesByTalentID[talent.ID]
			if found {
				membershipRoleValue = fmt.Sprintf("<@&%d>", role.ID)
			} else {
				membershipRoleValue = "⛔ Not yet configured"
			}
		} else {
			membershipRoleValue = "⛔ Not yet configured"
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "Membership Role",
			Value: membershipRoleValue,
		})
		embeds = append(embeds, embed)
	}
	return embeds
}
