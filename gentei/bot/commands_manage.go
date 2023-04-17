package bot

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/member-gentei/member-gentei/gentei/bot/commands"
	"github.com/member-gentei/member-gentei/gentei/bot/templates"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/member-gentei/member-gentei/gentei/ent/guildrole"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
	"github.com/rs/zerolog"
)

func (b *DiscordBot) handleManageAuditSet(
	ctx context.Context,
	logger zerolog.Logger,
	i *discordgo.InteractionCreate,
	cmd *discordgo.ApplicationCommandInteractionDataOption,
) (*discordgo.WebhookEdit, error) {
	var (
		options              = subcommandOptionsMap(cmd)
		targetChannel        = options["channel"].ChannelValue(b.session)
		targetChannelID, err = strconv.ParseUint(targetChannel.ID, 10, 64)
	)
	if err != nil {
		return nil, err
	}
	guildID, _, err := getMessageAttributionIDs(i)
	if err != nil {
		return nil, err
	}
	// check if anything changed
	dg, err := b.db.Guild.Get(ctx, guildID)
	if err != nil {
		return nil, err
	}
	if dg.AuditChannel == targetChannelID {
		return &discordgo.WebhookEdit{
			Content: ptr(fmt.Sprintf("%s is already the configured audit log channel for this Discord server.", targetChannel.Mention())),
		}, nil
	}
	// test that we can actually send messages to this channel
	_, err = b.session.ChannelMessageSendEmbed(targetChannel.ID, commands.AuditTestMessageEmbed)
	if err != nil {
		var restErr *discordgo.RESTError
		if errors.As(err, &restErr) {
			if restErr.Message != nil {
				switch restErr.Message.Code {
				case discordgo.ErrCodeMissingAccess, discordgo.ErrCodeMissingPermissions:
					return &discordgo.WebhookEdit{
						Content: ptr("We don't have permission to send embed messages to that test channel. Please re-run this command after granting the bot or its role the `Send Messages` and `Embed Links` permissions on that channel."),
					}, nil
				}
			}
		}
		return nil, fmt.Errorf("error sending test message: %w", err)
	}
	err = b.db.Guild.UpdateOneID(dg.ID).
		SetAuditChannel(targetChannelID).
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error setting Guild.AuditChannel: %w", err)
	}
	b.auditChannelCache.Delete(dg.ID)
	return &discordgo.WebhookEdit{
		Content: ptr(fmt.Sprintf("Role changes performed by this bot will now be posted to %s.", targetChannel.Mention())),
	}, nil
}

func (b *DiscordBot) handleManageAuditOff(
	ctx context.Context,
	logger zerolog.Logger,
	i *discordgo.InteractionCreate,
	cmd *discordgo.ApplicationCommandInteractionDataOption,
) (*discordgo.WebhookEdit, error) {
	guildID, _, err := getMessageAttributionIDs(i)
	if err != nil {
		return nil, err
	}
	// check if anything changed
	dg, err := b.db.Guild.Get(ctx, guildID)
	if err != nil {
		return nil, err
	}
	if dg.AuditChannel == 0 {
		return &discordgo.WebhookEdit{
			Content: ptr("No audit log channel is configured for this server."),
		}, nil
	}
	err = b.db.Guild.UpdateOneID(dg.ID).
		ClearAuditChannel().
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error unsetting audit log channel: %w", err)
	}
	return &discordgo.WebhookEdit{
		Content: ptr(fmt.Sprintf("<#%d> will no longer receive audit log messages.", dg.AuditChannel)),
	}, nil
}

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
	roleID, err := strconv.ParseUint(role.ID, 10, 64)
	if err != nil {
		return nil, err
	}
	// ensure that the YouTubeTalent exists
	talent, err := b.db.YouTubeTalent.Query().
		Where(youtubetalent.ID(channelID)).
		First(ctx)
	if ent.IsNotFound(err) {
		return &discordgo.WebhookEdit{
			Content: ptr(fmt.Sprintf("Unknown YouTube channel - please check for a typo in the channel ID or add the channel first through https://gentei.tindabox.net/app/enroll?server=%s", i.GuildID)),
		}, nil
	} else if err != nil {
		return nil, err
	}
	// overwrite role mapping
	guildID, _, err := getMessageAttributionIDs(i)
	if err != nil {
		return nil, err
	}
	existingRole, err := b.db.GuildRole.Query().
		WithTalent().
		Where(
			guildrole.HasGuildWith(guild.ID(guildID)),
			guildrole.HasTalentWith(youtubetalent.ID(channelID)),
		).
		Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, err
	} else if err == nil {
		// already mapped
		if existingRole.Edges.Talent.ID == channelID {
			return &discordgo.WebhookEdit{
				Content: ptr(fmt.Sprintf("%s is already the role mapped to this YouTube channel.", role.Mention())),
			}, nil
		} else {
			existingTalent := existingRole.Edges.Talent
			return &discordgo.WebhookEdit{
				Content: ptr(templates.MustRender(templates.RoleAlreadyMapped, templates.RoleAlreadyMappedContext{
					ChannelID:   existingTalent.ID,
					ChannelName: existingTalent.ChannelName,
					RoleMention: role.Mention(),
				})),
			}, nil
		}
	}
	// verify that we can add and remove ourselves from the role
	botUserID, err := strconv.ParseUint(b.session.State.User.ID, 10, 64)
	if err != nil {
		return nil, err
	}
	// add and remove self from role to make sure it works
	err = b.applyRole(ctx, guildID, roleID, botUserID, true, "role permissions test", true)
	if err != nil {
		return &discordgo.WebhookEdit{
			Content: ptr(templates.MustRender(templates.RolePermissionFailure, templates.RolePermissionFailureContext{
				Action:      "add",
				RoleMention: role.Mention(),
			})),
		}, err
	}
	err = b.applyRole(ctx, guildID, roleID, botUserID, false, "role permissions test", true)
	if err != nil {
		return &discordgo.WebhookEdit{
			Content: ptr(templates.MustRender(templates.RolePermissionFailure, templates.RolePermissionFailureContext{
				Action:      "remove",
				RoleMention: role.Mention(),
			})),
		}, err
	}
	// save
	err = b.db.GuildRole.Create().
		SetID(roleID).
		SetName(role.Name).
		SetGuildID(guildID).
		SetTalentID(talent.ID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &discordgo.WebhookEdit{
		Content: ptr(templates.MustRender(templates.RoleApplied, templates.RoleAppliedContext{
			ChannelName: talent.ChannelName,
			RoleMention: role.Mention(),
		})),
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
		dg, err := b.db.Guild.Query().
			WithRoles().
			Where(
				guild.ID(guildID),
				guild.HasYoutubeTalentsWith(youtubetalent.ID(youtubeID)),
			).
			First(ctx)
		if ent.IsNotFound(err) {
			return &discordgo.WebhookEdit{
				Content: ptr(fmt.Sprintf("There is no role mapping in this Discord server to unmap for the specified YouTube channel (`%s`).", youtubeID)),
			}, nil
		} else if err != nil {
			return nil, err
		}
		for _, role := range dg.Edges.Roles {
			if role.Edges.Talent.ID == youtubeID {
				err = b.db.GuildRole.DeleteOne(role).Exec(ctx)
				if err != nil {
					return nil, err
				}
				return &discordgo.WebhookEdit{
					Content: ptr(fmt.Sprintf("<@&%d> has been unmapped.", role.ID)),
				}, nil
			}
		}
		return nil, fmt.Errorf("mysterious fallthrough case on unmap")
	}
	if roleOption, exists := options["role"]; exists {
		roleValue := roleOption.RoleValue(nil, "")
		role := i.ApplicationCommandData().Resolved.Roles[roleValue.ID]
		roleID, err := strconv.ParseUint(role.ID, 10, 64)
		if err != nil {
			return nil, err
		}
		count, err := b.db.GuildRole.Delete().
			Where(
				guildrole.HasGuildWith(guild.ID(guildID)),
				guildrole.ID(roleID),
			).
			Exec(ctx)
		if err != nil {
			return nil, err
		}
		if count == 0 {
			return &discordgo.WebhookEdit{
				Content: ptr(fmt.Sprintf("%s is not currently mapped to a YouTube talent.", role.Mention())),
			}, nil
		}
		return &discordgo.WebhookEdit{
			Content: ptr(fmt.Sprintf("%s has been unmapped.", role.Mention())),
		}, nil
	}
	return &discordgo.WebhookEdit{
		Content: ptr("Please specify either the `youtube-channel-id` or `role` option."),
	}, nil
}
