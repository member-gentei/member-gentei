package bot

import (
	"context"
	"fmt"
	"strconv"

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
		// do not overwrite if
		// ...this role is mapped to a different talent
		for existingCid, mapping := range guildSettings.RoleMapping {
			if mapping.ID == role.ID {
				existingTalent, err := b.db.YouTubeTalent.Get(ctx, existingCid)
				if err != nil {
					return nil, err
				}
				return &discordgo.WebhookEdit{
					Content: fmt.Sprintf(
						"%s is already mapped to **%s** (%s). If you meant to remap this role to a new YouTube talent, please unmap this role first.",
						role.Mention(),
						existingTalent.ChannelName,
						existingTalent.ID,
					),
				}, nil
			}
		}
		// ...or nothing changed
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
	// verify that we can add and remove ourselves from the role
	roleID, err := strconv.ParseUint(role.ID, 10, 64)
	if err != nil {
		return nil, err
	}
	botUserID, err := strconv.ParseUint(b.session.State.User.ID, 10, 64)
	if err != nil {
		return nil, err
	}
	// add and remove self from role to make sure it works
	err = b.applyRole(ctx, guildID, roleID, botUserID, true)
	if err != nil {
		return &discordgo.WebhookEdit{
			Content: fmt.Sprintf(
				"We failed to confirm permissions to add users to %s. Please re-run this command after granting the bot permission to add + remove users from that role.",
				role.Mention(),
			),
		}, err
	}
	err = b.applyRole(ctx, guildID, roleID, botUserID, false)
	if err != nil {
		return &discordgo.WebhookEdit{
			Content: fmt.Sprintf(
				"We failed to confirm permissions to remove users from %s. Please re-run this command after granting the bot permission to add + remove users from that role.",
				role.Mention(),
			),
		}, err
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

func (b *DiscordBot) handleManageUnmap(
	ctx context.Context,
	logger zerolog.Logger,
	i *discordgo.InteractionCreate,
	cmd *discordgo.ApplicationCommandInteractionDataOption,
) (*discordgo.WebhookEdit, error) {
	var (
		options = subcommandOptionsMap(cmd)
	)
	guildID, _, err := getMessageAttributionIDs(i)
	if err != nil {
		return nil, err
	}
	// by YouTube channel ID
	if cidVal, exists := options["youtube-channel-id"]; exists {
		youtubeID := cidVal.StringValue()
		dg, err := b.db.Guild.Get(ctx, guildID)
		if err != nil {
			return nil, err
		}
		if dg.Settings == nil || dg.Settings.RoleMapping == nil || dg.Settings.RoleMapping[youtubeID].ID == "" {
			return &discordgo.WebhookEdit{
				Content: fmt.Sprintf("There is no role mapping in this Discord server to unmap for the specified YouTube channel (`%s`).", youtubeID),
			}, nil
		}
		settings := dg.Settings
		roleID := settings.RoleMapping[youtubeID].ID
		delete(settings.RoleMapping, youtubeID)
		err = b.db.Guild.UpdateOneID(guildID).
			SetSettings(settings).
			Exec(ctx)
		if err != nil {
			return nil, err
		}
		return &discordgo.WebhookEdit{
			Content: fmt.Sprintf("<@&%s> has been unmapped.", roleID),
		}, nil
	}
	if roleOption, exists := options["role"]; exists {
		roleValue := roleOption.RoleValue(nil, "")
		role := i.ApplicationCommandData().Resolved.Roles[roleValue.ID]
		dg, err := b.db.Guild.Get(ctx, guildID)
		if err != nil {
			return nil, err
		}
		var (
			settings = dg.Settings
			talentID string
		)
		if settings.RoleMapping == nil {
			settings.RoleMapping = map[string]schema.GuildSettingsRoleMapping{}
		}
		for key, mapping := range settings.RoleMapping {
			if mapping.ID == role.ID {
				talentID = key
				break
			}
		}
		if talentID == "" {
			return &discordgo.WebhookEdit{
				Content: fmt.Sprintf("%s is not currently mapped to a YouTube talent.", role.Mention()),
			}, nil
		}
		delete(settings.RoleMapping, talentID)
		err = b.db.Guild.UpdateOneID(guildID).
			SetSettings(settings).
			Exec(ctx)
		if err != nil {
			return nil, err
		}
		return &discordgo.WebhookEdit{
			Content: fmt.Sprintf("%s has been unmapped.", role.Mention()),
		}, nil
	}
	return &discordgo.WebhookEdit{
		Content: "Please specify either the `youtube-channel-id` or `role` option.",
	}, nil
}
