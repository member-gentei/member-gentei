package roles

import (
	"context"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type ApplyRoleResult struct {
	Attempts int
	Error    error
}

var (
	DefaultTimeout       = time.Second * 10
	DefaultRetryInterval = time.Second * 2
)

// ApplyRole starts a retry loop to try to apply a role before a specified deadline.
//
// Give this a context that can be canceled by a GUILD_MEMBER_UPDATE event pre-empting the next retry.
// https://discord.com/developers/docs/topics/gateway#guild-member-update
//
// From observed behavior, we can infer that role application is an asynchronous/eventually consistent API call that
// sometimes fails to work on large (>5000 user) servers. So we just keep applying the role until it finally works
// or *we* time out trying to apply it.
func ApplyRole(applyCtx context.Context, session *discordgo.Session, guildID, userID, roleID uint64, add bool) (result <-chan ApplyRoleResult) {
	var (
		guildIDStr = strconv.FormatUint(guildID, 10)
		roleIDStr  = strconv.FormatUint(roleID, 10)
		userIDStr  = strconv.FormatUint(userID, 10)
		resultChan = make(chan ApplyRoleResult, 1)
		baseLogger = log.With().Uint64("guildID", guildID).Uint64("userID", userID).Uint64("roleID", roleID).Logger()
		applyRole  func() error
	)
	if add {
		applyRole = func() error {
			return session.GuildMemberRoleAdd(guildIDStr, userIDStr, roleIDStr)
		}
	} else {
		applyRole = func() error {
			return session.GuildMemberRoleRemove(guildIDStr, userIDStr, roleIDStr)
		}
	}
	deadline := time.Now().Add(DefaultTimeout)
	ctx, cancel := context.WithDeadline(applyCtx, deadline)
	go func() {
		var (
			attempts int
			err      error
		)
		defer cancel()
		for attempts = 1; ; attempts++ {
			var (
				breakOut bool
				logger   = baseLogger.With().Int("attempt", attempts).Logger()
			)
			logger.Debug().Msg("attempting to apply role")
			select {
			case <-ctx.Done():
				logger.Debug().Msg("role apply timed out/cancelled")
				err = ctx.Err()
				breakOut = true
			default:
				if attempts > 1 {
					// check if the last attempt worked
					var member *discordgo.Member
					member, err = session.GuildMember(guildIDStr, userIDStr)
					if err != nil {
						logger.Debug().Msg("error fetching user roles during apply")
						breakOut = true
						break
					}
					for _, role := range member.Roles {
						if role == roleIDStr {
							logger.Debug().Msg("query informed succesful role apply")
							breakOut = true
							break
						}
					}
					if breakOut {
						break
					}
				}
				err = applyRole()
				if err != nil {
					logger.Debug().Msg("error attempting to apply role")
					breakOut = true
					break
				}
				time.Sleep(DefaultRetryInterval)
			}
			if breakOut {
				break
			}
		}
		resultChan <- ApplyRoleResult{
			Attempts: attempts,
			Error:    err,
		}
	}()
	return resultChan
}
