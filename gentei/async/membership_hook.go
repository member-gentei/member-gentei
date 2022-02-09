package async

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/member-gentei/member-gentei/gentei/membership"
	"github.com/rs/zerolog/log"
)

type psMembershipChangeHandler struct {
	ctx   context.Context
	topic *pubsub.Topic
	membership.ChangeHandler
}

func (p *psMembershipChangeHandler) GainedMembership(userMembershipID int) {
	p.publishAsync(
		ApplyMembershipPSMessage{
			UserMembershipID: userMembershipID,
			Gained:           true,
		},
	)
}
func (p *psMembershipChangeHandler) LostMembership(userMembershipID int) {
	p.publishAsync(
		ApplyMembershipPSMessage{
			UserMembershipID: userMembershipID,
			Gained:           false,
		},
	)
}

func (p *psMembershipChangeHandler) publishAsync(message ApplyMembershipPSMessage) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Err(err).
			Interface("message", message).
			Msg("error marshalling ApplyMemsershipPSMessage")
		return
	}
	p.topic.Publish(p.ctx, &pubsub.Message{
		Attributes: map[string]string{
			"type": string(ApplyMembershipType),
		},
		Data: data,
	})

}

func NewPubSubMembershipChangeHandler(ctx context.Context, topic *pubsub.Topic) membership.ChangeHandler {
	return &psMembershipChangeHandler{
		ctx:   ctx,
		topic: topic,
	}
}
