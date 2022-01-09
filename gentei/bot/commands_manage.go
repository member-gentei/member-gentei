package bot

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/schema"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
	"github.com/rs/zerolog"
)

func (b *DiscordBot) handleManageMap(
	ctx context.Context,
	logger zerolog.Logger,
	i *discordgo.InteractionCreate,
	cmd *discordgo.ApplicationCommandInteractionDataOption,
) (*discordgo.WebhookEdit, error) {
	var (
		options   = subcommandOptionsMap(cmd)
		channelID = options["youtube-channel-id"].StringValue()
		roleValue = options["role"].RoleValue(nil, "")
	)
	role := i.ApplicationCommandData().Resolved.Roles[roleValue.ID]
	// ensure that the YouTubeTalent exists
	talent, err := b.db.YouTubeTalent.Query().
		Where(youtubetalent.ID(channelID)).
		First(ctx)
	if ent.IsNotFound(err) {
		return &discordgo.WebhookEdit{
			Content: fmt.Sprintf("This is not the ID of a channel we have on file. Please check for a typo in the channel ID or add the channel first through https://gentei.tindabox.net/app/enroll?server=%s", i.GuildID),
		}, nil
	} else if err != nil {
		return nil, err
	}
	// overwrite role mapping
	guildID, _, err := getMessageAttributionIDs(i)
	if err != nil {
		return nil, err
	}
	dg, err := b.db.Guild.Get(ctx, guildID)
	if err != nil {
		return nil, err
	}
	guildSettings := dg.Settings
	if guildSettings == nil {
		guildSettings = &schema.GuildSettings{
			RoleMapping: map[string]schema.GuildSettingsRoleMapping{},
		}
	}
	if guildSettings.RoleMapping == nil {
		dg.Settings.RoleMapping = map[string]schema.GuildSettingsRoleMapping{}
	} else {
		// ...unless nothing changed.
		existingRole, exists := guildSettings.RoleMapping[channelID]
		if exists && existingRole.ID == role.ID {
			return &discordgo.WebhookEdit{
				Content: fmt.Sprintf("%s is already the role mapped to this YouTube channel.", role.Mention()),
			}, nil
		}
	}
	guildSettings.RoleMapping[channelID] = schema.GuildSettingsRoleMapping{
		ID:   role.ID,
		Name: role.Name,
	}
	err = b.db.Guild.UpdateOneID(guildID).
		SetSettings(guildSettings).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &discordgo.WebhookEdit{
		Content: fmt.Sprintf("Role for membership to **%s** is now %s.\nThis role mapping will take effect during the next daily membership check for any users already registered with Gentei.", talent.ChannelName, role.Mention()),
	}, nil
}
