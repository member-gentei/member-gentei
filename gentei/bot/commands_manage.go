package bot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/member-gentei/member-gentei/gentei/bot/templates"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/member-gentei/member-gentei/gentei/ent/guildrole"
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
			Content: fmt.Sprintf("Unknown channel - please check for a typo in the channel ID or add the channel first through https://gentei.tindabox.net/app/enroll?server=%s", i.GuildID),
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
		Where(guildrole.ID(roleID)).
		Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, err
	} else if err == nil {
		// already mapped
		if existingRole.Edges.Talent.ID == channelID {
			return &discordgo.WebhookEdit{
				Content: fmt.Sprintf("%s is already the role mapped to this YouTube channel.", role.Mention()),
			}, nil
		} else {
			existingTalent := existingRole.Edges.Talent
			return &discordgo.WebhookEdit{
				Content: templates.MustRender(templates.RoleAlreadyMapped, templates.RoleAlreadyMappedContext{
					ChannelID:   existingTalent.ID,
					ChannelName: existingTalent.ChannelName,
					RoleMention: role.Mention(),
				}),
			}, nil
		}
	}
	// verify that we can add and remove ourselves from the role
	botUserID, err := strconv.ParseUint(b.session.State.User.ID, 10, 64)
	if err != nil {
		return nil, err
	}
	// add and remove self from role to make sure it works
	err = b.applyRole(ctx, guildID, roleID, botUserID, true)
	if err != nil {
		return &discordgo.WebhookEdit{
			Content: templates.MustRender(templates.RolePermissionFailure, templates.RolePermissionFailureContext{
				Action:      "add",
				RoleMention: role.Mention(),
			}),
		}, err
	}
	err = b.applyRole(ctx, guildID, roleID, botUserID, false)
	if err != nil {
		return &discordgo.WebhookEdit{
			Content: templates.MustRender(templates.RolePermissionFailure, templates.RolePermissionFailureContext{
				Action:      "remove",
				RoleMention: role.Mention(),
			}),
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
		Content: templates.MustRender(templates.RoleApplied, templates.RoleAppliedContext{
			ChannelName: talent.ChannelName,
			RoleMention: role.Mention(),
		}),
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
				Content: fmt.Sprintf("There is no role mapping in this Discord server to unmap for the specified YouTube channel (`%s`).", youtubeID),
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
					Content: fmt.Sprintf("<@&%d> has been unmapped.", role.ID),
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
				Content: fmt.Sprintf("%s is not currently mapped to a YouTube talent.", role.Mention()),
			}, nil
		}
		return &discordgo.WebhookEdit{
			Content: fmt.Sprintf("%s has been unmapped.", role.Mention()),
		}, nil
	}
	return &discordgo.WebhookEdit{
		Content: "Please specify either the `youtube-channel-id` or `role` option.",
	}, nil
}
