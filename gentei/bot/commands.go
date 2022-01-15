// commands.go contains app command delegation and such

package bot

import (
	"context"
	"strconv"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/bwmarrin/discordgo"
	"github.com/member-gentei/member-gentei/gentei/bot/commands"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	eaCommandName        = "ea-gentei"
	eaCommandDescription = "Gentei membership management (early access/dev)"
	prodCommandName      = "gentei"
)

var (
	// TODO: maybe not hardcode this
	earlyAccessGuilds = []string{
		"929085334430556240",
		"497603499190779914",
	}
	earlyAccessCommand = &discordgo.ApplicationCommand{
		Name:        eaCommandName,
		Description: eaCommandDescription,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "info",
				Description: "Show server and/or membership info",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "manage",
				Description: "Admin-only: manage memberships",
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "map",
						Description: "Change role mapping of a channel -> Discord role.",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Options: []*discordgo.ApplicationCommandOption{
							{
								Name:        "youtube-channel-id",
								Description: "The YouTube channel ID whose memberships should be monitored. (e.g. UCAL_ZudIZXhCDrniD4ZQobQ)",
								Type:        discordgo.ApplicationCommandOptionString,
								Required:    true,
							},
							{
								Name:        "role",
								Description: "The Discord role for members of this YouTube channel",
								Type:        discordgo.ApplicationCommandOptionRole,
								Required:    true,
							},
						},
					},
					{
						Name:        "unmap",
						Description: "Remove role mapping by referencing either the YouTube channel or Discord role.",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Options: []*discordgo.ApplicationCommandOption{
							{
								Name:        "youtube-channel-id",
								Description: "The YouTube channel ID whose memberships we should stop monitoring. (e.g. UCAL_ZudIZXhCDrniD4ZQobQ)",
								Type:        discordgo.ApplicationCommandOptionString,
							},
							{
								Name:        "role",
								Description: "The Discord role for members of a YouTube channel",
								Type:        discordgo.ApplicationCommandOptionRole,
							},
						},
					},
				},
			},
		},
	}
	globalCommand = &discordgo.ApplicationCommand{}
)

// slashResponseFunc should return an error message if error != nil. If not, it's treated as a real error.
type slashResponseFunc func(logger zerolog.Logger) (*discordgo.WebhookEdit, error)

func (b *DiscordBot) handleInfo(ctx context.Context, i *discordgo.InteractionCreate) {
	b.deferredReply(ctx, i, "info", true, func(logger zerolog.Logger) (*discordgo.WebhookEdit, error) {
		var response *discordgo.WebhookEdit
		// if this is a DM, fetch user info
		if i.User != nil {
			// TODO
		} else {
			// fetch guild info + user-relevant info
			userID, err := strconv.ParseUint(i.Member.User.ID, 10, 64)
			if err != nil {
				return nil, err
			}
			guildID, err := strconv.ParseUint(i.GuildID, 10, 64)
			if err != nil {
				return nil, err
			}
			isThisUser := func(uq *ent.UserQuery) {
				uq.Where(user.ID(userID))
			}
			dg, err := b.db.Guild.Query().
				WithYoutubeTalents().
				WithAdmins(isThisUser).
				WithMembers(isThisUser).
				Where(guild.ID(guildID)).
				First(ctx)
			if err != nil {
				return nil, err
			}
			response = &discordgo.WebhookEdit{
				Content: "Here's how this server is configured.",
				Embeds:  commands.GetGuildInfoEmbeds(dg, len(dg.Edges.Admins) > 0),
			}
		}
		return response, nil
	})
}

func (b *DiscordBot) handleManage(ctx context.Context, i *discordgo.InteractionCreate) {
	if i.User != nil {
		b.replyNoDM(i)
		return
	}
	b.deferredReply(ctx, i, "manage", true,
		// the calling user has to be an admin of some sort to run this command
		func(logger zerolog.Logger) (*discordgo.WebhookEdit, error) {
			return b.deferredRequireModeratorOrAdmin(ctx, i)
		},
		func(logger zerolog.Logger) (*discordgo.WebhookEdit, error) {
			manageCmd := i.ApplicationCommandData().Options[0]
			subcommand := manageCmd.Options[0]
			switch subcommand.Name {
			case "map":
				return b.handleManageMap(ctx, logger, i, subcommand)
			case "unmap":
				return b.handleManageUnmap(ctx, logger, i, subcommand)
			default:
				return &discordgo.WebhookEdit{
					Content: "You've somehow sent an unknown `/gentei manage` command. Discord is not supposed to allow this to happen so... try reloading this browser window or your Discord client? :thinking:",
				}, nil
			}
		},
	)
}

// helpers

// deferredReply can only be called once, but it'll process each responseFunc until it gets a WebHookEdit payload to send from one of them as an early termination mechanism
func (b *DiscordBot) deferredReply(ctx context.Context, i *discordgo.InteractionCreate, commandName string, ephemeral bool, responseFuncs ...slashResponseFunc) {
	var (
		logger = log.With().
			Str("command", commandName).
			Str("userID", i.Member.User.ID).
			Str("guildID", i.GuildID).
			Logger()
		appID = b.session.State.User.ID
	)
	err := b.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: 1 << 6,
		},
	})
	if err != nil {
		logger.Err(err).Msg("error sending message deferral, dropping reply")
		return
	}
	for _, responseFunc := range responseFuncs {
		response, err := responseFunc(logger)
		if err != nil {
			if response == nil {
				logger.Err(err).Msg("responseFunc did not return an error response, sending oops message")
				response = &discordgo.WebhookEdit{
					Content: "A mysterious error occured, and this bot's author has been notified. Try again later? :(",
				}
			} else {
				logger.Warn().Err(err).Msg("sending responseFunc error response")
			}
		}
		if response == nil {
			continue
		}
		_, err = b.session.InteractionResponseEdit(appID, i.Interaction, response)
		if err != nil {
			logger.Err(err).Msg("error generating response")
		}
		return
	}
}

func (b *DiscordBot) replyNoDM(i *discordgo.InteractionCreate) {
	b.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "This command can only be used in a Discord server",
		},
	})
}

var (
	insufficientPermissionsWebhookEdit = &discordgo.WebhookEdit{
		Content: "You do not have permission to run this command in this server.",
	}
)

func (b *DiscordBot) deferredRequireModeratorOrAdmin(ctx context.Context, i *discordgo.InteractionCreate) (*discordgo.WebhookEdit, error) {
	guildID, userID, err := getMessageAttributionIDs(i)
	if err != nil {
		log.Err(err).Msg("error converting IDs to uint64")
		return nil, err
	}
	var rolePredicates = []*sql.Predicate{
		sqljson.ValueContains(guild.FieldAdminSnowflakes, userID),
		sqljson.ValueContains(guild.FieldModeratorSnowflakes, userID),
	}
	for _, roleIDStr := range i.Member.Roles {
		roleID, err := strconv.ParseUint(roleIDStr, 10, 64)
		if err != nil {
			log.Err(err).Msg("error converting role ID to uint64")
			return nil, err
		}
		rolePredicates = append(
			rolePredicates,
			sqljson.ValueContains(guild.FieldAdminSnowflakes, roleID),
			sqljson.ValueContains(guild.FieldModeratorSnowflakes, roleID),
		)
	}
	hasPermission, err := b.db.Guild.Query().
		Where(
			guild.ID(guildID),
		).
		Where(func(s *sql.Selector) {
			s.Where(sql.Or(rolePredicates...))
		}).
		Exist(ctx)
	if err != nil {
		log.Err(err).Msg("error checking for admin/mod privs")
		return nil, err
	}
	if !hasPermission {
		return insufficientPermissionsWebhookEdit, nil
	}
	return nil, nil
}

func subcommandOptionsMap(cmd *discordgo.ApplicationCommandInteractionDataOption) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	options := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(cmd.Options))
	for _, option := range cmd.Options {
		options[option.Name] = option
	}
	return options
}
