// Pubsub-backed asynchronous processing.

package async

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/member-gentei/member-gentei/gentei/apis"
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
	// TODO: remove next week.
	UserRegistration    interface{}                 `json:",omitempty"`
	UserDelete          json.Number                 `json:",omitempty"`
	YouTubeRegistration *YouTubeRegistrationMessage `json:",omitempty"`
}

type ApplyMembershipPSMessage struct {
	UserMembershipID int
	Gained           bool
	Reason           string

	DeleteUserID json.Number `json:",omitempty"`
}

type YouTubeRegistrationMessage struct {
	UserID json.Number
}

type UserDeleteMessage struct {
	UserID json.Number
}

var (
	trashUserRegistrationMessageDeadline time.Time
	trashUserRegistrationHours           = 24 * 7.0
)

// ListenGeneral polls a pubsub subscription for GeneralType messages.
func ListenGeneral(
	parentCtx context.Context,
	db *ent.Client,
	youTubeConfig *oauth2.Config,
	subscription *pubsub.Subscription,
	botTopic *pubsub.Topic,
	setChangeReason func(string),
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
		if message.UserRegistration != nil && time.Since(trashUserRegistrationMessageDeadline).Hours() < trashUserRegistrationHours {
			m.Ack()
			return
		}
		if message.YouTubeRegistration != nil {
			userID, err := strconv.ParseUint(message.YouTubeRegistration.UserID.String(), 10, 64)
			if err != nil {
				log.Err(err).
					Str("unparsedUserID", message.YouTubeRegistration.UserID.String()).
					Msg("error decoding UserID as uint64")
				m.Ack()
				return
			}
			if err := ProcessYouTubeRegistration(ctx, db, youTubeConfig, userID, setChangeReason); err != nil {
				if errors.Is(err, apis.ErrNoYouTubeTokenForUser) {
					log.Warn().Err(err).Uint64("userID", userID).Msg("acking YouTube registration with errors")
					m.Ack()
					return
				}
				log.Err(err).Uint64("userID", userID).Msg("error processing YouTube registration")
			} else {
				m.Ack()
				return
			}
		}
		if message.UserDelete.String() != "" {
			userID, err := strconv.ParseUint(message.UserDelete.String(), 10, 64)
			if err != nil {
				log.Err(err).
					Str("unparsedUserID", message.YouTubeRegistration.UserID.String()).
					Msg("error decoding UserID as uint64")
				m.Ack()
				return
			}
			if err = ProcessUserDelete(ctx, db, botTopic, userID); err != nil {
				log.Err(err).Uint64("userID", userID).Msg("error processing user deletion")
			} else {
				m.Ack()
				return
			}
		}
		m.Nack()
		log.Warn().Str("data", string(m.Data)).Msg("nacking unhandled message - is this mid-deploy?")
	})
}

func publishPubSubMessage(ctx context.Context, topic *pubsub.Topic, messageType PSMessageType, message interface{}) error {
	marshalled, err := json.Marshal(message)
	if err != nil {
		return err
	}
	pr := topic.Publish(ctx, &pubsub.Message{
		Attributes: map[string]string{
			"type": string(messageType),
		},
		Data: marshalled,
	})
	_, err = pr.Get(ctx)
	return err
}

func init() {
	var err error
	trashUserRegistrationMessageDeadline, err = time.Parse(time.UnixDate, "Sun Feb 27 15:50:06 UTC 2022")
	if err != nil {
		panic(err)
	}
}
