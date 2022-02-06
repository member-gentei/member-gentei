package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/membership"
)

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
