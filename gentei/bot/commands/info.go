package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/member-gentei/member-gentei/gentei/ent"
)

func GetGuildInfoEmbeds(dg *ent.Guild) []*discordgo.MessageEmbed {
	var (
		embeds []*discordgo.MessageEmbed
	)
	for _, role := range dg.Edges.Roles {
		talent := role.Edges.Talent
		embed := &discordgo.MessageEmbed{
			Type:  discordgo.EmbedTypeRich,
			Title: talent.ChannelName,
			URL:   fmt.Sprintf("https://www.youtube.com/channel/%s", talent.ID),
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: talent.ThumbnailURL,
			},
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "Membership Role",
			Value: fmt.Sprintf("<@&%d>", role.ID),
		})
		embeds = append(embeds, embed)
	}
	return embeds
}
