// Pubsub-backed asynchronous processing.

package async

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"cloud.google.com/go/pubsub"
	"github.com/member-gentei/member-gentei/gentei/apis"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi"
)

// PSMessageType is a Pub/Sub message type enum, packed into the attributes field.
type PSMessageType string

const (
	GeneralType         PSMessageType = "general"
	ApplyMembershipType PSMessageType = "apply-membership"
)

type GeneralPSMessage struct {
	UserDelete          *DeleteUserMessage          `json:",omitempty"`
	YouTubeRegistration *YouTubeRegistrationMessage `json:",omitempty"`
	YouTubeDelete       json.Number                 `json:",omitempty"`
}

type ApplyMembershipPSMessage struct {
	ApplySingle  *ApplySingleMessage `json:",omitempty"`
	DeleteSingle *DeleteUserMessage  `json:",omitempty"`
	EnforceAll   *EnforceAllMessage  `json:",omitempty"`
}

type YouTubeRegistrationMessage struct {
	UserID json.Number
}

type DeleteUserMessage struct {
	UserID json.Number
	Reason string `json:",omitempty"`
}

type ApplySingleMessage struct {
	UserMembershipID int
	Gained           bool
	Reason           string
}

type EnforceAllMessage struct {
	DryRun bool
	Reason string
}

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
		if message.YouTubeRegistration != nil {
			userID, err := strconv.ParseUint(message.YouTubeRegistration.UserID.String(), 10, 64)
			if err != nil {
				log.Err(err).
					Str("unparsedUserID", message.YouTubeRegistration.UserID.String()).
					Msg("error decoding UserID as uint64")
				m.Ack()
				return
			}
			logger := log.With().Str("userID", strconv.FormatUint(userID, 10)).Logger()
			if err := ProcessYouTubeRegistration(ctx, db, youTubeConfig, userID, setChangeReason); err != nil {
				if ent.IsNotFound(err) {
					m.Ack()
					logger.Warn().Err(err).Msg("registered user not found, acking YouTube registration")
					return
				}
				if errors.Is(err, apis.ErrNoYouTubeTokenForUser) {
					logger.Warn().Err(err).Msg("acking YouTube registration with errors")
					m.Ack()
					return
				}
				var gErr *googleapi.Error
				if errors.As(err, &gErr) {
					if gErr.Message == "Invalid Credentials" {
						logger.Warn().Err(err).Msg("discarding bad registration creds")
						m.Ack()
						return
					}
				}
				logger.Err(err).Msg("error processing YouTube registration")
			} else {
				m.Ack()
				return
			}
		}
		if message.UserDelete != nil {
			udIDStr := message.UserDelete.UserID.String()
			userID, err := strconv.ParseUint(udIDStr, 10, 64)
			if err != nil {
				log.Err(err).
					Str("unparsedUserID", udIDStr).
					Msg("error decoding UserID as uint64")
				m.Ack()
				return
			}
			logger := log.With().Str("userID", strconv.FormatUint(userID, 10)).Logger()
			if err = ProcessUserDelete(ctx, db, botTopic, userID, message.UserDelete.Reason); err != nil {
				logger.Err(err).Msg("error processing user deletion")
			} else {
				logger.Info().Msg("issued delete for user")
				m.Ack()
				return
			}
		}
		if message.YouTubeDelete.String() != "" {
			userID, err := strconv.ParseUint(message.YouTubeDelete.String(), 10, 64)
			if err != nil {
				log.Err(err).
					Str("unparsedUserID", message.YouTubeDelete.String()).
					Msg("error decoding UserID as uint64")
				m.Ack()
				return
			}
			if err = ProcessYouTubeDelete(ctx, db, botTopic, userID); err != nil {
				log.Err(err).Str("userID", strconv.FormatUint(userID, 10)).
					Msg("error processing YouTube user deletion")
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
