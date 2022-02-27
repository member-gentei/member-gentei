package commands

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

// CreateAuditLogEmbed creates the embed used for audit log messages.
func CreateAuditLogEmbed(
	userID uint64, userAvatarURL string,
	reason string,
	add bool,
) *discordgo.MessageEmbed {
	// color coded!
	var (
		action string
		color  int
	)
	if add {
		color = 0x00bd00
		action = "Add role"
	} else {
		color = 0xbd0000
		action = "Remove role"
	}
	embed := &discordgo.MessageEmbed{
		Color: color,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Name",
				Value:  fmt.Sprintf("<@%d>", userID),
				Inline: true,
			},
			{
				Name:   "ID",
				Value:  strconv.FormatUint(userID, 10),
				Inline: true,
			},
			{
				Name:   "Action",
				Value:  action,
				Inline: true,
			},
			{
				Name: "Timestamp",
				// https://discord.com/developers/docs/reference#message-formatting-timestamp-styles
				Value:  fmt.Sprintf("<t:%d:f>", time.Now().Unix()),
				Inline: true,
			},
			{
				Name:   "Reason",
				Value:  reason,
				Inline: true,
			},
		},
	}
	if userAvatarURL != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: userAvatarURL,
		}
	}
	return embed
}
