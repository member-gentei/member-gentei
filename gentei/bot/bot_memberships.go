package bot

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/member-gentei/member-gentei/gentei/ent/guildrole"
	"github.com/member-gentei/member-gentei/gentei/ent/predicate"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/member-gentei/member-gentei/gentei/ent/usermembership"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func (b *DiscordBot) enforceAllRoles(ctx context.Context, dryRun bool, reason string) error {
	b.roleEnforcementMutex.Lock()
	defer b.roleEnforcementMutex.Unlock()
	// iterate through all configured roles
	grs, err := b.db.GuildRole.Query().
		WithGuild().
		Where(
			guildrole.HasTalentWith(youtubetalent.DisabledIsNil()),
		).
		Order(ent.Asc(guildrole.FieldID)).
		All(ctx)
	if err != nil {
		return fmt.Errorf("error getting GuildRoles: %w", err)
	}
	log.Info().Int("count", len(grs)).Msg("enforcing members-only roles")
	for i := range grs {
		var (
			gr         = grs[i]
			guildIDStr = strconv.FormatUint(gr.Edges.Guild.ID, 10)
			roleIDStr  = strconv.FormatUint(gr.ID, 10)
			logger     = log.With().Str("guildID", guildIDStr).Str("roleID", roleIDStr).Logger()
		)
		logger.Info().Int("i", i).Msg("enforcing members-only role")
		if err = b.enforceRole(ctx, gr, dryRun, reason); err != nil {
			// TODO: inform point of contact about permissions failure
			var restErr *discordgo.RESTError
			if errors.As(err, &restErr) {
				if restErr.Message.Code == discordgo.ErrCodeMissingAccess {
					logger.Warn().Err(err).Msg("permissions error working with Guild")
				}
			} else if errors.Is(err, discordgo.ErrStateNotFound) {
				logger.Warn().Msg("Guild relevant to enforcement not in session state")
				pruned, err := b.pruneGuildIfAbsent(ctx, gr.Edges.Guild.ID)
				if err != nil {
					return fmt.Errorf("error pruning Guild: %w", err)
				} else if pruned {
					logger.Info().Msg("pruned departed Guild, no membership changes will be communicated")
				} else {
					logger.Error().Msg("Guild was absent from state, consider restarting")
				}
			} else if err != nil {
				logger.Err(err).Msg("error enforcing GuildRole")
				return fmt.Errorf("error enforcing GuildRole: %w", err)
			}
			// n.b. rollup log is emitted by the enforceRole() call
		}
	}
	log.Info().Int("count", len(grs)).Msg("members-only role enforcement run complete")
	return nil
}

func (b *DiscordBot) enforceRole(ctx context.Context, gr *ent.GuildRole, dryRun bool, reason string) error {
	dg, err := gr.Edges.GuildOrErr()
	if err != nil {
		return err
	}
	var (
		guildID    = dg.ID
		roleID     = gr.ID
		guildIDStr = strconv.FormatUint(guildID, 10)
		roleIDStr  = strconv.FormatUint(roleID, 10)
		logger     = log.With().Str("guildID", guildIDStr).Str("roleID", roleIDStr).Logger()
	)
	// gather users who should have this role
	var (
		shouldHaveRole = map[uint64]bool{}
		yesterdayIsh   = time.Now().Add(-time.Hour * 24)
	)
	ums, err := b.db.GuildRole.QueryUserMemberships(gr).
		WithUser().
		All(ctx)
	if err != nil {
		return fmt.Errorf("error getting users that should have role: %w", err)
	}
	for _, um := range ums {
		userID := um.Edges.User.ID
		if um.FailCount == 0 {
			// the last check worked
			shouldHaveRole[userID] = true
		} else if !um.LastVerified.IsZero() && um.FirstFailed.After(yesterdayIsh) {
			// a check worked before, but it failed today. They have [insert grace period] to fix it.
			// TODO: insert configurable grace period here
			shouldHaveRole[userID] = true
			logger.Info().Str("userID", strconv.FormatUint(userID, 10)).Msg("role membership in grace period")
		}
	}
	// wait until we have all members
	mutex := b.guildMemberLoadMutexes[guildIDStr]
	if !mutex.TryLock() {
		logger.Info().Msg("waiting until Guild Members are done loading")
		mutex.Lock()
	}
	defer mutex.Unlock()
	// compile the changeset
	var (
		toAdd     []string
		toRemove  []string
		keepCount int
	)
	dGuild, err := b.session.State.Guild(guildIDStr)
	if err != nil {
		return fmt.Errorf("error getting Guild %d from discordgo session state: %w", guildID, err)
	}
	for _, member := range dGuild.Members {
		uid, err := strconv.ParseUint(member.User.ID, 10, 64)
		if err != nil {
			return err
		}
		// determine if this user should be granted / removed the role
		if shouldHaveRole[uid] {
			if sliceContains(roleIDStr, member.Roles) {
				keepCount++
			} else {
				// "should have the role and don't already have the role"
				toAdd = append(toAdd, member.User.ID)
			}
		} else if sliceContains(roleIDStr, member.Roles) {
			// "should not have the role and has the role"
			toRemove = append(toRemove, member.User.ID)
		}
	}
	logger.Info().
		Int("addCount", len(toAdd)).
		Int("removeCount", len(toRemove)).
		Int("keepCount", keepCount).
		Int("guildMembers", len(dGuild.Members)).
		Bool("dryRun", dryRun).
		Msg("role enforcement rollup")
	if dryRun {
		return nil
	}
	// smash through the queue, 4 applications at a time
	var (
		eg, egCtx = errgroup.WithContext(ctx)
	)
	eg.SetLimit(4)
	for _, uid := range toAdd {
		userID, err := strconv.ParseUint(uid, 10, 64)
		if err != nil {
			return err
		}
		eg.Go(func() error {
			return b.applyRole(egCtx, guildID, roleID, userID, true, reason, false)
		})
	}
	for _, uid := range toRemove {
		userID, err := strconv.ParseUint(uid, 10, 64)
		if err != nil {
			return err
		}
		eg.Go(func() error {
			return b.applyRole(egCtx, guildID, roleID, userID, false, reason, false)
		})
	}
	if err = eg.Wait(); err != nil {
		return fmt.Errorf("error applying changes for role, aborted: %w", err)
	}
	logger.Info().Msg("role enforcement complete")
	return nil
}

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
		logger.Warn().Err(err).Msg("failed to query relevant User object")
		return &ent.ConstraintError{}
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
		logger.Warn().Err(err).Msg("failed to query relevant User object")
		return err
	}
	existingRoles, err := db.GuildRole.Query().
		Where(guildrole.HasUserMembershipsWith(usermembership.ID(userMembershipID))).
		WithGuild().
		WithTalent().
		All(ctx)
	if err != nil {
		logger.Warn().Err(err).Msg("failed to query existing GuildRole objects")
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
	err := b.applyRole(ctx, guildID, roleID, userID, false, reason, true)
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
	err := b.applyRole(ctx, guildID, roleID, userID, true, reason, true)
	if err != nil {
		roleLogger.Warn().Err(err).Msg("failed to grant role membership")
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

// returns whether the Guild and its roles were pruned.
func (b *DiscordBot) pruneGuildIfAbsent(ctx context.Context, guildID uint64) (bool, error) {
	var (
		guildIDStr = strconv.FormatUint(guildID, 10)
		restErr    *discordgo.RESTError
	)
	_, err := b.session.Guild(guildIDStr)
	if errors.As(err, &restErr) {
		if restErr.Message.Code == discordgo.ErrCodeMissingAccess {
			// we're not in the guild, so remove it and its roles
			// (this should CASCADE DELETE to all appropriate objects)
			err = b.db.Guild.DeleteOneID(guildID).Exec(ctx)
			if ent.IsNotFound(err) {
				err = nil
			}
			return true, err
		}
	}
	return false, err
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
