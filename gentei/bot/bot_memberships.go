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
	"golang.org/x/sync/semaphore"
)

func (b *DiscordBot) enforceAllRoles(ctx context.Context, dryRun bool, reason string) error {
	// iterate through all configured roles
	grs, err := b.db.GuildRole.Query().
		WithGuild().
		Where(
			guildrole.HasTalentWith(youtubetalent.DisabledIsNil()),
		).
		All(ctx)
	if err != nil {
		return fmt.Errorf("error getting GuildRoles: %w", err)
	}
	for _, gr := range grs {
		var (
			guildIDStr = strconv.FormatUint(gr.Edges.Guild.ID, 10)
			roleIDStr  = strconv.FormatUint(gr.ID, 10)
			logger     = log.With().Str("guildID", guildIDStr).Str("roleID", roleIDStr).Logger()
		)
		if err = b.enforceRole(ctx, gr, dryRun, reason); err != nil {
			// TODO: inform point of contact about permissions failure
			var restErr *discordgo.RESTError
			if errors.As(err, &restErr) {
				if restErr.Message.Code == discordgo.ErrCodeMissingAccess {
					logger.Warn().Err(err).Msg("permissions error working with Guild")
				}
			} else {
				logger.Err(err).Msg("error enforcing GuildRole")
				return fmt.Errorf("error enforcing GuildRole: %w", err)
			}
		}
	}
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
		mutex      = b.roleRWMutex.GetOrCreate(roleIDStr)
	)
	logger.Debug().Msg("acquiring RWMutex for role")
	mutex.Lock()
	defer mutex.Unlock()
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
	// compile the changeset
	var (
		toAdd    []string
		toRemove []string
	)
	dGuild, err := b.session.Guild(guildIDStr)
	if err != nil {
		return fmt.Errorf("error getting Guild %d from discordgo session: %w", guildID, err)
	}
	for _, member := range dGuild.Members {
		uid, err := strconv.ParseUint(member.User.ID, 10, 64)
		if err != nil {
			return err
		}
		// determine if this user should be granted / removed the role
		if shouldHaveRole[uid] && !sSliceContains(roleIDStr, member.Roles) {
			toAdd = append(toAdd, member.User.ID)
		} else if sSliceContains(roleIDStr, member.Roles) {
			toRemove = append(toRemove, member.User.ID)
		}
	}
	logger.Info().
		Int("addCount", len(toAdd)).
		Int("removeCount", len(toRemove)).
		Bool("dryRun", dryRun).
		Msg("role enforcement rollup")
	if dryRun {
		return nil
	}
	// smash through the queue
	var (
		sem       = semaphore.NewWeighted(4)
		eg, egCtx = errgroup.WithContext(ctx)
	)
	for _, uid := range toAdd {
		userID, err := strconv.ParseUint(uid, 10, 64)
		if err != nil {
			return err
		}
		eg.Go(func() error {
			e := sem.Acquire(egCtx, 1)
			if e != nil {
				return e
			}
			defer sem.Release(1)
			return b.applyRole(egCtx, guildID, roleID, userID, true, reason)
		})
	}
	for _, uid := range toRemove {
		userID, err := strconv.ParseUint(uid, 10, 64)
		if err != nil {
			return err
		}
		eg.Go(func() error {
			e := sem.Acquire(egCtx, 1)
			if e != nil {
				return e
			}
			defer sem.Release(1)
			return b.applyRole(egCtx, guildID, roleID, userID, false, reason)
		})
	}
	if err = eg.Wait(); err != nil {
		return fmt.Errorf("error applying changes for role, aborted: %w", err)
	}
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

func sSliceContains(needle string, haystack []string) bool {
	for _, hay := range haystack {
		if needle == hay {
			return true
		}
	}
	return false
}
