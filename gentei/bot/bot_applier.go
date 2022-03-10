package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"cloud.google.com/go/pubsub"
	"github.com/bwmarrin/discordgo"
	"github.com/member-gentei/member-gentei/gentei/async"
	"github.com/member-gentei/member-gentei/gentei/bot/templates"
	"github.com/rs/zerolog/log"
)

func (b *DiscordBot) StartPSApplier(parentCtx context.Context, sub *pubsub.Subscription) {
	var (
		pCtx, cancel = context.WithCancel(parentCtx)
	)
	b.cancelPSApplier = cancel
	go func() {
		defer cancel()
		err := sub.Receive(pCtx, b.handlePSMessage)
		if err != nil {
			log.Err(err).Msg("bot PSApplier crashed?")
		}
	}()
}

func (b *DiscordBot) handlePSMessage(ctx context.Context, m *pubsub.Message) {
	if typeAttribute := m.Attributes["type"]; typeAttribute != string(async.ApplyMembershipType) {
		log.Warn().Str("typeAttribute", typeAttribute).Msg("non apply-membership message made it past the filter?")
		m.Ack()
		return
	}
	var message async.ApplyMembershipPSMessage
	err := json.Unmarshal(m.Data, &message)
	if err != nil {
		log.Warn().Str("data", string(m.Data)).Msg("acking message that cannot be decoded as JSON")
		m.Ack()
		return
	}
	switch {
	case message.DeleteUserID != "":
		var (
			userIDStr    = message.DeleteUserID.String()
			userID, err  = strconv.ParseUint(message.DeleteUserID.String(), 10, 64)
			reasonDetail = message.Reason
			reason       string
		)
		if err != nil {
			log.Err(err).
				Str("unparsedUserID", userIDStr).
				Msg("error decoding UserID as uint64")
			m.Ack()
		}
		if reasonDetail != "" {
			reason = fmt.Sprintf("user deleted (%s)", reasonDetail)
		} else {
			reason = "user deleted"
		}
		logger := log.With().Str("userID", strconv.FormatUint(userID, 10)).Logger()
		err = b.revokeMembershipsByUserID(ctx, userID, reason)
		if err != nil {
			logger.Err(err).Msg("error revoking all memberships before deletion")
			return
		}
		// now actually delete the user
		err = b.db.User.DeleteOneID(userID).Exec(ctx)
		if err != nil {
			return
		}
		m.Ack()
		// best-effort attempt at sending the user deletion DM
		ch, err := b.session.UserChannelCreate(userIDStr)
		if err != nil {
			logger.Err(err).Msg("error creating UserChannel to inform of deletion")
			return
		}
		msg, err := b.session.ChannelMessageSend(ch.ID, templates.PlaintextUserDeleted)
		if err != nil {
			var restErr *discordgo.RESTError
			if errors.As(err, &restErr) && restErr.Message.Code == discordgo.ErrCodeCannotSendMessagesToThisUser {
				logger.Warn().Err(err).
					Msg("unable to send deletion confirmation message")
			} else {
				logger.Err(err).
					Msg("error sending deletion confirmation message")
			}
		} else {
			logger.Info().
				Interface("messageMetadata", msg).
				Msg("sent deletion confirmation message")
		}
	case message.Gained:
		err = b.grantMemberships(ctx, b.db, message.UserMembershipID, message.Reason)
		if err != nil {
			log.Err(err).Msg("error granting memberships")
			return
		}
	default:
		err = b.revokeMemberships(ctx, b.db, message.UserMembershipID, message.Reason)
		if err != nil {
			log.Err(err).Msg("error revoking memberships")
			return
		}
	}
	m.Ack()
}
