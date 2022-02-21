// commands.go contains app command delegation and such

package bot

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/bwmarrin/discordgo"
	"github.com/member-gentei/member-gentei/gentei/bot/commands"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
	"github.com/member-gentei/member-gentei/gentei/membership"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	eaCommandName          = "ea-gentei"
	eaCommandDescription   = "Gentei membership management (early access/dev)"
	prodCommandName        = "gentei"
	prodCommandDescription = "Gentei membership management"
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
			{
				Name:        "check",
				Description: "Check membership eligiblity.",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
		},
	}
	globalCommand = &discordgo.ApplicationCommand{
		Name:        prodCommandName,
		Description: prodCommandDescription,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "info",
				Description: "Show server and/or membership info",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "check",
				Description: "Check membership eligiblity.",
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
)

func (b *DiscordBot) PushCommands(global, earlyAccess bool) error {
	// get self - we might not have started a websocket
	self, err := b.session.User("@me")
	if err != nil {
		return fmt.Errorf("error getting @me: %w", err)
	}
	if earlyAccess {
		for _, guildID := range earlyAccessGuilds {
			log.Info().Str("guildID", guildID).Msg("loading early access command")
			_, err := b.session.ApplicationCommandCreate(self.ID, guildID, earlyAccessCommand)
			if err != nil {
				return fmt.Errorf("error loading early access command to guild '%s': %w", guildID, err)
			}
		}
	}
	if global {
		log.Info().Msg("pushing global command - new command set will be available in 1~2 hours")
		pushed, err := b.session.ApplicationCommandBulkOverwrite(self.ID, "", []*discordgo.ApplicationCommand{globalCommand})
		if err != nil {
			return fmt.Errorf("error loading global command: %w", err)
		}
		var versions []string
		for _, cmd := range pushed {
			versions = append(versions, cmd.Version)
		}
		log.Info().Strs("versions", versions).Msg("push global command")
	}
	return nil
}

// slashResponseFunc should return an error message if error != nil. If not, it's treated as a real error.
type slashResponseFunc func(logger zerolog.Logger) (*discordgo.WebhookEdit, error)

const mysteriousErrorMessage = "A mysterious error occured, and this bot's author has been notified. Try again later? :("

func (b *DiscordBot) handleCheck(ctx context.Context, i *discordgo.InteractionCreate) {
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
			u, err := b.db.User.Query().
				Where(
					user.ID(userID),
					user.YoutubeIDNotNil(),
					user.YoutubeIDNEQ(""),
				).
				First(ctx)
			if ent.IsNotFound(err) {
				return &discordgo.WebhookEdit{
					Content: "Please register and link your YouTube account to https://gentei.tindabox.net to check memberships.",
				}, nil
			} else if err != nil {
				return nil, err
			}
			if time.Since(u.LastCheck).Minutes() < 1 {
				return &discordgo.WebhookEdit{Content: mysteriousErrorMessage}, nil
			}
			var (
				logger = log.With().Uint64("userID", userID).Uint64("guildID", guildID).Logger()
			)
			talents, err := b.db.YouTubeTalent.Query().
				Where(youtubetalent.HasGuildsWith(guild.ID(guildID))).
				All(ctx)
			if err != nil {
				return nil, err
			}
			if len(talents) == 0 {
				return &discordgo.WebhookEdit{
					Content: "This server has not configured memberships yet or has paused membership management. Please be discreet until server moderation announces something!",
				}, nil
			}
			talentIDs := make([]string, len(talents))
			for i := range talents {
				talentIDs[i] = talents[i].ID
			}
			logger.Debug().Strs("talentIDs", talentIDs).Msg("performing /gentei check")
			results, err := membership.CheckForUser(ctx, b.db, b.youTubeConfig, userID, &membership.CheckForUserOptions{
				ChannelIDs: talentIDs,
			})
			if err != nil {
				logger.Err(err).Msg("error checking memberships for user")
				return &discordgo.WebhookEdit{Content: mysteriousErrorMessage}, nil
			}

			logger.Debug().Interface("results", results).Msg("/gentei check results")
			err = membership.SaveMemberships(ctx, b.db, userID, results)
			if err != nil {
				logger.Err(err).Msg("error saving UserMembership objects for user")
				return &discordgo.WebhookEdit{Content: mysteriousErrorMessage}, nil
			}
			// apply changes
			var (
				gainedIDs = b.qch.GetGained()
				lostIDs   = b.qch.GetLost()
			)
			for _, gainedID := range gainedIDs {
				err = b.grantMemberships(ctx, b.db, gainedID)
				if err != nil {
					logger.Err(err).Int("gainedID", gainedID).Msg("error granting memberships for user")
					return &discordgo.WebhookEdit{Content: mysteriousErrorMessage}, nil
				}
			}
			for _, lostID := range lostIDs {
				err = b.revokeMemberships(ctx, b.db, lostID)
				if err != nil {
					logger.Err(err).Int("lostID", lostID).Msg("error revoking memberships for user")
					return &discordgo.WebhookEdit{Content: mysteriousErrorMessage}, nil
				}
			}
			embeds, err := commands.CreateMembershipInfoEmbeds(ctx, b.db, userID, guildID, gainedIDs, lostIDs)
			if err != nil {
				logger.Err(err).Msg("error creating embeds for reply")
				return &discordgo.WebhookEdit{Content: mysteriousErrorMessage}, nil
			}
			response = &discordgo.WebhookEdit{
				Content: "Discord roles have been applied - see below for details.",
				Embeds:  embeds,
			}
		}
		return response, nil
	})
}

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

// deferredReply can only be called once, but it'll process each responseFunc in serial. If it gets a WebHookEdit payload, it sends that and stops processing later responseFuncs.
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
