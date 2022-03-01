package bot

import (
	"context"
	"strconv"
	"time"

	"github.com/member-gentei/member-gentei/gentei/bot/commands"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/rs/zerolog/log"
)

var (
	// when we last nagged the server owner (!) about audit logs
	lastAuditLogNagTimes = map[uint64]time.Time{}
)

// auditLog sends audit logs, if configured on the Discord Guild.
func (b *DiscordBot) auditLog(ctx context.Context, guildID, userID, roleID uint64, add bool, reason string) {
	// TODO: cache the audit log channel ID heavily
	var (
		logger = log.With().
			Str("guildID", strconv.FormatUint(guildID, 10)).
			Str("userID", strconv.FormatUint(userID, 10)).
			Str("roleID", strconv.FormatUint(roleID, 10)).
			Logger()
		avatarURL string
	)
	dg, err := b.db.Guild.Query().
		Where(
			guild.ID(guildID),
			guild.AuditChannelNotNil(),
			guild.Not(guild.AuditChannel(0)),
		).
		Only(ctx)
	if ent.IsNotFound(err) {
		return
	}
	if err != nil {
		logger.Err(err).Msg("error querying for audit log channel")
		return
	}
	dgUser, err := b.session.User(strconv.FormatUint(userID, 10))
	if err != nil {
		logger.Warn().Err(err).Msg("error getting Discord user avatar - user deleted?")
		avatarURL = ""
		return
	} else {
		avatarURL = dgUser.AvatarURL("")
	}
	// send audit log message
	logger = logger.With().Str("auditChannel", strconv.FormatUint(dg.AuditChannel, 10)).Logger()
	_, err = b.session.ChannelMessageSendEmbed(
		strconv.FormatUint(dg.AuditChannel, 10),
		commands.CreateAuditLogEmbed(userID, avatarURL, roleID, reason, add),
	)
	// TODO: nag admins instead of me about things not working
	if err != nil && lastAuditLogNagTimes[dg.AuditChannel].Before(time.Now().Add(-time.Duration(time.Hour*24))) {
		lastAuditLogNagTimes[dg.AuditChannel] = time.Now()
		logger.Err(err).Msg("audit log delivery failure")
	}
}
