package bot

import (
	"context"
	"strconv"

	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/member-gentei/member-gentei/gentei/ent/guildrole"
	"github.com/member-gentei/member-gentei/gentei/ent/predicate"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/member-gentei/member-gentei/gentei/ent/usermembership"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// grantMemberships makes the edges of userMembershipID authoritative.
//
// Because we don't know the failure modes of role application yet, we mask all apply errors so we can deduce it from logs + alerts.
func (b *DiscordBot) grantMemberships(ctx context.Context, db *ent.Client, userMembershipID int, reason string) error {
	var (
		logger = log.With().Int("userMembershipID", userMembershipID).Logger()
	)
	// big joins!
	userID, err := db.User.Query().
		Where(user.HasMembershipsWith(usermembership.ID(userMembershipID))).
		OnlyID(ctx)
	if err != nil {
		logger.Err(err).Msg("failed to query relevant User object")
		return ent.ConstraintError{}
	}
	logger = logger.With().Logger()
	existingRoles, err := db.GuildRole.Query().
		Where(guildrole.HasUserMembershipsWith(usermembership.ID(userMembershipID))).
		WithGuild().
		WithTalent().
		All(ctx)
	if err != nil {
		logger.Err(err).Msg("failed to query existing GuildRole objects")
		return err
	}
	// apply any new roles
	// (we need to query this separately to add edges)
	guildRoles, err := db.GuildRole.Query().
		Where(userNewlyVerifiedFor(userMembershipID)...).
		WithGuild().
		WithTalent().
		All(ctx)
	if err != nil {
		logger.Err(err).Msg("failed to query newly verified GuildRole objects")
		return err
	}
	for _, guildRole := range guildRoles {
		err := b.grantRole(ctx, db, logger, userID, guildRole, userMembershipID, reason)
		_ = err
	}
	// re-apply existing roles
	for _, guildRole := range existingRoles {
		err := b.grantRole(ctx, db, logger, userID, guildRole, 0, reason)
		_ = err
	}
	return err
}

func (b *DiscordBot) revokeMemberships(ctx context.Context, db *ent.Client, userMembershipID int, reason string) error {
	var (
		logger = log.With().Int("userMembershipID", userMembershipID).Logger()
	)
	// big joins!
	userID, err := db.User.Query().
		Where(user.HasMembershipsWith(usermembership.ID(userMembershipID))).
		OnlyID(ctx)
	if err != nil {
		logger.Err(err).Msg("failed to query relevant User object")
		return err
	}
	existingRoles, err := db.GuildRole.Query().
		Where(guildrole.HasUserMembershipsWith(usermembership.ID(userMembershipID))).
		WithGuild().
		WithTalent().
		All(ctx)
	if err != nil {
		logger.Err(err).Msg("failed to query existing GuildRole objects")
		return err
	}
	// revoke everything
	for _, guildRole := range existingRoles {
		var (
			guildID    = guildRole.Edges.Guild.ID
			roleID     = guildRole.ID
			roleLogger = logger.With().
					Str("talentID", guildRole.Edges.Talent.ID).
					Logger()
		)
		b.revokeMembership(ctx, roleLogger, guildID, roleID, userID, reason)
	}
	return nil
}

func (b *DiscordBot) revokeMembershipsByUserID(ctx context.Context, userID uint64, reason string) error {
	guildRoles, err := b.db.GuildRole.Query().
		WithGuild().
		WithTalent().
		Where(
			guildrole.HasUserMembershipsWith(
				usermembership.HasUserWith(user.ID(userID)),
			),
		).
		All(ctx)
	if err != nil {
		log.Err(err).Str("userID", strconv.FormatUint(userID, 10)).
			Msg("failed to query existing GuildRole objects for wide revoke")
		return err
	}
	for _, guildRole := range guildRoles {
		var (
			guildID = guildRole.Edges.Guild.ID
			roleID  = guildRole.ID
			logger  = log.With().
				Str("talentID", guildRole.Edges.Talent.ID).
				Logger()
		)
		b.revokeMembership(ctx, logger, guildID, roleID, userID, reason)
	}
	return nil
}

func (b *DiscordBot) revokeMembership(ctx context.Context, baseLogger zerolog.Logger, guildID, roleID, userID uint64, reason string) {
	var (
		logger = baseLogger.With().
			Str("guildID", strconv.FormatUint(guildID, 10)).
			Str("roleID", strconv.FormatUint(roleID, 10)).
			Logger()
	)
	err := b.applyRole(ctx, guildID, roleID, userID, false, reason)
	if err != nil {
		logger.Err(err).Msg("failed to revoke role membership")
	} else {
		logger.Info().Msg("revoked role membership")
	}
}

func (b *DiscordBot) grantRole(
	ctx context.Context,
	db *ent.Client,
	logger zerolog.Logger,
	userID uint64,
	guildRole *ent.GuildRole,
	addEdgeUserMembershipID int,
	reason string,
) error {
	var (
		guildID    = guildRole.Edges.Guild.ID
		roleID     = guildRole.ID
		roleLogger = logger.With().
				Str("guildID", strconv.FormatUint(guildID, 10)).
				Str("roleID", strconv.FormatUint(roleID, 10)).
				Str("talentID", guildRole.Edges.Talent.ID).
				Logger()
	)
	err := b.applyRole(ctx, guildID, roleID, userID, true, reason)
	if err != nil {
		roleLogger.Err(err).Msg("failed to grant role membership")
		return err
	}
	roleLogger.Info().Msg("granted newly verified membership")
	if addEdgeUserMembershipID != 0 {
		err = db.GuildRole.UpdateOneID(guildRole.ID).
			AddUserMembershipIDs(addEdgeUserMembershipID).
			Exec(ctx)
		if err != nil {
			roleLogger.Err(err).Msg("failed to add UserMembership <-> GuildRole edge")
			return err
		}
	}
	return nil
}

func userNewlyVerifiedFor(userMembershipID int) []predicate.GuildRole {
	return append(
		userVerifiedFor(userMembershipID),
		guildrole.Not(
			guildrole.HasUserMembershipsWith(usermembership.ID(userMembershipID)),
		),
	)
}

func userVerifiedFor(userMembershipID int) []predicate.GuildRole {
	return []predicate.GuildRole{
		guildrole.HasTalentWith(
			youtubetalent.HasMembershipsWith(usermembership.ID(userMembershipID)),
		),
		guildrole.HasGuildWith(
			guild.Or(
				guild.HasAdminsWith(user.HasMembershipsWith(usermembership.ID(userMembershipID))),
				guild.HasMembersWith(user.HasMembershipsWith(usermembership.ID(userMembershipID))),
			),
		),
	}
}
