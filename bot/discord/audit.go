package discord

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func (d *discordBot) emitMemberAuditLog(
	auditLogChannelID string, action roleAction,
	userID, avatarURL, reason string,
) {
	// color coded!
	var color int
	switch action {
	case roleAdd:
		color = 0x00bd00
	case roleRevoke:
		color = 0xbd0000
	}
	_, err := d.dgSession.ChannelMessageSendEmbed(auditLogChannelID, &discordgo.MessageEmbed{
		Color: color,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: avatarURL,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Name",
				Value:  fmt.Sprintf("<@%s>", userID),
				Inline: true,
			},
			{
				Name:   "ID",
				Value:  userID,
				Inline: true,
			},
			{
				Name:   "Action",
				Value:  action.String(),
				Inline: true,
			},
			{
				Name: "Timestamp",
				// RFC 3339 with some readability changes
				Value:  time.Now().In(time.UTC).Format("2006/01/02 15:04:05Z"),
				Inline: true,
			},
			{
				Name:   "Reason",
				Value:  reason,
				Inline: true,
			},
		},
	})
	if err != nil {
		log.Err(err).Str("action", action.String()).
			Str("auditLogChannelID", auditLogChannelID).
			Str("userID", userID).
			Msg("error emitting Discord role audit log")
	}
}
