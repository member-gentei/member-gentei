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
	"github.com/member-gentei/member-gentei/gentei/ent/guildrole"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
	"github.com/member-gentei/member-gentei/gentei/membership"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	commandName            = "gentei"
	eaCommandDescription   = "Gentei membership management (early access)"
	prodCommandDescription = "Gentei membership management"
)

var (
	// TODO: maybe not hardcode this
	earlyAccessGuilds = []string{
		"497603499190779914",
	}
	earlyAccessCommand = &discordgo.ApplicationCommand{
		Name:        commandName,
		Description: eaCommandDescription,
		Options: []*discordgo.ApplicationCommandOption{
			// info
			globalCommand.Options[0],
			// check
			globalCommand.Options[1],
			{
				Name:        "manage",
				Description: "Admin-only: manage memberships",
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Options: []*discordgo.ApplicationCommandOption{
					// audit
					globalCommand.Options[2].Options[0],
					// map
					globalCommand.Options[2].Options[1],
					// unmap
					globalCommand.Options[2].Options[2],
				},
			},
		},
	}
	globalCommand = &discordgo.ApplicationCommand{
		Name:        commandName,
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
				Description: "Admin-only: manage memberships and server settings",
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "audit-set",
						Description: "Set/change role management audit log settings.",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Options: []*discordgo.ApplicationCommandOption{
							{
								Name:         "channel",
								Description:  "The Discord channel to recieve audit logs.",
								Type:         discordgo.ApplicationCommandOptionChannel,
								ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
								Required:     true,
							},
						},
					},
					{
						Name:        "audit-off",
						Description: "Turns off role management audit logs.",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
					},
					{
						Name:        "map",
						Description: "Set/change role mapping of a channel -> Discord role.",
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

func (b *DiscordBot) ClearCommands(global, earlyAccess bool) error {
	// get self - we might not have started a websocket
	self, err := b.session.User("@me")
	if err != nil {
		return fmt.Errorf("error getting @me: %w", err)
	}
	if global {
		globalCommands, err := b.session.ApplicationCommands(self.ID, "")
		if err != nil {
			return fmt.Errorf("error getting global commands: %w", err)
		}
		for _, c := range globalCommands {
			err = b.session.ApplicationCommandDelete(c.ApplicationID, "", c.ID)
			if err != nil {
				return fmt.Errorf("error deleting global command: %w", err)
			}
			log.Info().Msg("cleared global command - it should take effect very soon")
		}
	}
	if earlyAccess {
		for _, guildID := range earlyAccessGuilds {
			guildCommands, err := b.session.ApplicationCommands(self.ID, guildID)
			if err != nil {
				return fmt.Errorf("error getting commands for guild '%s': %w", guildID, err)
			}
			for _, c := range guildCommands {
				err = b.session.ApplicationCommandDelete(c.ApplicationID, guildID, c.ID)
				if err != nil {
					return fmt.Errorf("error deleting guild command: %w", err)
				}
			}
			log.Info().Str("guildID", guildID).Msg("cleared guild command")
		}
	}
	return nil
}

// slashResponseFunc should return an error message if error != nil. If not, it's treated as a real error.
type slashResponseFunc func(logger zerolog.Logger) (*discordgo.WebhookEdit, error)

var (
	mysteriousErrorMessage = ptr("A mysterious error occured, and this bot's author has been notified. Try again later? :(")
)

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
					Content: ptr("Please register and link your YouTube account to https://gentei.tindabox.net to check memberships."),
				}, nil
			} else if err != nil {
				return nil, err
			}
			ensureRegisteredUserHasGuildEdge(context.Background(), b.db, guildID, userID)
			if time.Since(u.LastCheck).Minutes() < 1 {
				tryAgainIn := u.LastCheck.Add(time.Minute).Unix()
				return &discordgo.WebhookEdit{
					Content: ptr(fmt.Sprintf("Membership checks are rate limited. Please try again <t:%d:R>.", tryAgainIn)),
				}, nil
			}
			var (
				logger = log.With().
					Str("userID", strconv.FormatUint(userID, 10)).
					Str("guildID", strconv.FormatUint(guildID, 10)).
					Logger()
			)
			talents, err := b.db.YouTubeTalent.Query().
				Where(youtubetalent.HasGuildsWith(guild.ID(guildID))).
				All(ctx)
			if err != nil {
				return nil, err
			}
			if len(talents) == 0 {
				return &discordgo.WebhookEdit{
					Content: ptr("This server has not configured memberships yet or has paused membership management. Please be discreet until server moderation announces something!"),
				}, nil
			}
			talentIDs := make([]string, len(talents))
			for i := range talents {
				talentIDs[i] = talents[i].ID
			}
			// add an User <-> Guild edge if they don't already have one
			isGuildMember, err := b.db.Guild.Query().
				Where(
					guild.ID(guildID),
					guild.HasMembersWith(user.ID(userID)),
				).
				Exist(ctx)
			if err != nil {
				return nil, err
			}
			if !isGuildMember {
				err = b.db.Guild.UpdateOneID(guildID).
					AddMemberIDs(userID).
					Exec(ctx)
				if err != nil {
					return nil, err
				}
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
			// apply all managed roles for this server only
			grs, err := b.db.GuildRole.Query().
				WithTalent().
				Where(
					guildrole.HasGuildWith(guild.ID(guildID)),
				).
				All(ctx)
			if err != nil {
				logger.Err(err).Msg("error fetching GuildRole objects for /gentei check apply")
				return &discordgo.WebhookEdit{Content: mysteriousErrorMessage}, nil
			}
			for _, role := range grs {
				shouldHaveRole := results.IsMember(role.Edges.Talent.ID)
				logger.Debug().
					Str("roleID", strconv.FormatUint(role.ID, 10)).
					Bool("shouldHaveRole", shouldHaveRole).
					Msg("check: applying role")
				err = b.applyRole(ctx, guildID, role.ID, userID, shouldHaveRole, "/gentei check (on-demand)", true)
				if err != nil {
					if IsDiscordError(err, discordgo.ErrCodeMissingPermissions) {
						logger.Warn().Err(err).Msg("Gentei lost permissions to manage a role - informing user")
						response = &discordgo.WebhookEdit{
							Content: ptr(fmt.Sprintf("The bot has lost permission to manage <@&%d>, so Gentei cannot apply all role changes! Please contact admins/moderators to restore permissions - once that's done, you can run `/gentei check` again to apply changes.", role.ID)),
						}
						return response, nil
					}
					if IsDiscordError(err, discordgo.ErrCodeUnknownRole) {
						logger.Warn().Err(err).Msg("Gentei was assigned a now removed role - informing user")
						response = &discordgo.WebhookEdit{
							Content: ptr(fmt.Sprintf("The role assigned to Gentei to manage membership to \"%s\" - which used to be <@&%d> - no longer exists, so Gentei cannot apply all role changes. Please contact admins/moderators to either assign a new role or remove role management - once that's done, you can run `/gentei check` again if a new role has been assigned.", role.Edges.Talent.ChannelName, role.ID)),
						}
						return response, nil
					}
					logger.Err(err).Msg("error applying role during check")
				}
			}
			embeds, err := commands.CreateMembershipInfoEmbeds(ctx, b.db, userID, guildID)
			if err != nil {
				logger.Err(err).Msg("error creating embeds for reply")
				return &discordgo.WebhookEdit{Content: mysteriousErrorMessage}, nil
			}
			if len(results.DisabledChannels) > 0 {
				// append embeds for channels
				disabledEmbeds, err := commands.GetDisabledChannelEmbeds(ctx, b.db, results.DisabledChannels)
				if err != nil {
					logger.Err(err).Msg("error getting disabled channel embeds")
					return &discordgo.WebhookEdit{Content: mysteriousErrorMessage}, nil
				}
				response = &discordgo.WebhookEdit{
					Content: ptr("Discord roles have been applied, but membership checks are currently disabled for one or more channels. See below for details."),
					Embeds:  ptr(append(embeds, disabledEmbeds...)),
				}
			} else {
				response = &discordgo.WebhookEdit{
					Content: ptr("Discord roles have been applied - see below for details."),
					Embeds:  ptr(embeds),
				}
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
			b.replyNoDM(i)
		} else {
			// fetch guild info + user-relevant info
			guildID, userID, err := getMessageAttributionIDs(i)
			if err != nil {
				return nil, err
			}
			go ensureRegisteredUserHasGuildEdge(context.Background(), b.db, guildID, userID)
			dg, err := b.db.Guild.Query().
				WithRoles(func(grq *ent.GuildRoleQuery) { grq.WithTalent() }).
				WithYoutubeTalents().
				Where(guild.ID(guildID)).
				Only(ctx)
			if err != nil {
				return nil, err
			}
			var (
				embeds = commands.GetGuildInfoEmbeds(dg)
			)
			if len(embeds) > 10 {
				response = &discordgo.WebhookEdit{
					Content: ptr(fmt.Sprintf("Due to technical Discord limitations, we can only show 10 of the %d memberships configured for this server. Please go to https://gentei.tindabox.net/app to view membership info for the rest.", len(embeds))),
					Embeds:  ptr(embeds[:10]),
				}
			} else {
				response = &discordgo.WebhookEdit{
					Content: ptr("Here's how this server is configured."),
					Embeds:  ptr(embeds),
				}
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
			case "audit-set":
				return b.handleManageAuditSet(ctx, logger, i, subcommand)
			case "audit-off":
				return b.handleManageAuditOff(ctx, logger, i, subcommand)
			case "map":
				return b.handleManageMap(ctx, logger, i, subcommand)
			case "unmap":
				return b.handleManageUnmap(ctx, logger, i, subcommand)
			default:
				return &discordgo.WebhookEdit{
					Content: ptr("You've somehow sent an unknown `/gentei manage` command. Discord is not supposed to allow this to happen so... try reloading this browser window or your Discord client? :thinking:"),
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
			Logger()
	)
	if i.Member != nil {
		logger = logger.With().
			Str("guildID", i.GuildID).
			Str("userID", i.Member.User.ID).
			Logger()
	} else {
		logger = logger.With().
			Str("userID", i.User.ID).
			Bool("dm", true).
			Logger()
	}
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
					Content: mysteriousErrorMessage,
				}
			} else {
				logger.Warn().Err(err).Msg("sending responseFunc error response")
			}
		}
		if response == nil {
			continue
		}
		_, err = b.session.InteractionResponseEdit(i.Interaction, response)
		if err != nil {
			logger.Err(err).Msg("error generating response")
		}
		return
	}
}

func (b *DiscordBot) replyNoDM(i *discordgo.InteractionCreate) {
	err := b.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "This command can only be used in a Discord server.",
		},
	})
	if err != nil {
		log.Err(err).Msg("error replying to slash command DM")
	}
}

var (
	insufficientPermissionsWebhookEdit = &discordgo.WebhookEdit{
		Content: ptr("You do not have permission to run this command in this server."),
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

// ctx is expected to be context.Background() most of the time.
func ensureRegisteredUserHasGuildEdge(ctx context.Context, db *ent.Client, guildID, userID uint64) {
	logger := log.With().
		Str("guildID", strconv.FormatUint(guildID, 10)).
		Str("userID", strconv.FormatUint(guildID, 10)).
		Logger()
	isMember, err := db.User.Query().
		Where(
			user.ID(userID),
			user.HasGuildsWith(guild.ID(guildID)),
		).
		Exist(ctx)
	if err != nil {
		logger.Err(err).Msg("error checking for edge between User and Guild")
	}
	if !isMember {
		err := db.User.UpdateOneID(userID).
			AddGuildIDs(guildID).
			Exec(ctx)
		if err != nil && ent.IsNotFound(err) {
			logger.Err(err).Msg("error adding edge between User and Guild")
		}
	}
}

func ptr[T any](o T) *T {
	return &o
}
