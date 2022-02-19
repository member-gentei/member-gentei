package bot

import (
	"context"

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
func (b *DiscordBot) grantMemberships(ctx context.Context, db *ent.Client, userMembershipID int) error {
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
		err := b.grantRole(ctx, db, logger, userID, guildRole, userMembershipID)
		_ = err
	}
	// re-apply existing roles
	for _, guildRole := range existingRoles {
		err := b.grantRole(ctx, db, logger, userID, guildRole, 0)
		_ = err
	}
	return err
}

func (b *DiscordBot) revokeMemberships(ctx context.Context, db *ent.Client, userMembershipID int) error {
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
					Uint64("guildID", guildID).
					Uint64("roleID", roleID).
					Str("talentID", guildRole.Edges.Talent.ID).
					Logger()
		)
		err = b.applyRole(ctx, guildID, roleID, userID, false)
		if err != nil {
			roleLogger.Err(err).Msg("failed to revoke role membership")
			continue
		}
		roleLogger.Info().Msg("revoked role membership")
	}
	return nil
}

func (b *DiscordBot) grantRole(
	ctx context.Context,
	db *ent.Client,
	logger zerolog.Logger,
	userID uint64,
	guildRole *ent.GuildRole,
	addEdgeUserMembershipID int,
) error {
	var (
		guildID    = guildRole.Edges.Guild.ID
		roleID     = guildRole.ID
		roleLogger = logger.With().
				Uint64("guildID", guildID).
				Uint64("roleID", roleID).
				Str("talentID", guildRole.Edges.Talent.ID).
				Logger()
	)
	err := b.applyRole(ctx, guildID, roleID, userID, true)
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
