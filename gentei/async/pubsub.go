// Pubsub-backed asynchronous processing.

package async

import (
	"context"
	"encoding/json"
	"strconv"

	"cloud.google.com/go/pubsub"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

// PSMessageType is a Pub/Sub message type enum, packed into the attributes field.
type PSMessageType string

const (
	GeneralType         PSMessageType = "general"
	ApplyMembershipType PSMessageType = "apply-membership"
)

type GeneralPSMessage struct {
	YouTubeRegistration *YouTubeRegistrationMessage `json:",omitempty"`
}

type ApplyMembershipPSMessage struct {
	UserMembershipID int
	Gained           bool
}

type YouTubeRegistrationMessage struct {
	UserID json.Number
}

// ListenGeneral polls a pubsub subscription for GeneralType messages.
func ListenGeneral(
	parentCtx context.Context,
	db *ent.Client,
	youTubeConfig *oauth2.Config,
	subscription *pubsub.Subscription,
) error {
	return subscription.Receive(parentCtx, func(ctx context.Context, m *pubsub.Message) {
		if typeAttribute := m.Attributes["type"]; typeAttribute != string(GeneralType) {
			log.Warn().Str("typeAttribute", typeAttribute).Msg("non-general message made it past the filter?")
			m.Ack()
			return
		}
		var message GeneralPSMessage
		err := json.Unmarshal(m.Data, &message)
		if err != nil {
			log.Warn().Str("data", string(m.Data)).Msg("acking message that cannot be decoded as JSON")
			m.Ack()
			return
		}
		if message.YouTubeRegistration != nil {
			userID, err := strconv.ParseUint(message.YouTubeRegistration.UserID.String(), 10, 64)
			if err != nil {
				log.Err(err).
					Str("unparsedUserID", message.YouTubeRegistration.UserID.String()).
					Msg("error decoding UserID as uint64")
				return
			}
			if err := ProcessYouTubeRegistration(ctx, db, youTubeConfig, userID); err != nil {
				log.Err(err).Uint64("userID", userID).Msg("error processing user registration")
			} else {
				m.Ack()
			}
		}
		m.Nack()
		log.Warn().Str("data", string(m.Data)).Msg("nacking unhandled message - is this mid-deploy?")
	})
}
