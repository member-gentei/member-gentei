package async

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/membership"
	"golang.org/x/oauth2"
)

func PublishGeneralMessage(ctx context.Context, topic *pubsub.Topic, message GeneralPSMessage) error {
	marshalled, err := json.Marshal(message)
	if err != nil {
		return err
	}
	pr := topic.Publish(ctx, &pubsub.Message{
		Attributes: map[string]string{
			"type": string(GeneralType),
		},
		Data: marshalled,
	})
	_, err = pr.Get(ctx)
	return err
}

// ProcessUserRegistration really only checks memberships and triggers changes. One day it might do something else?
func ProcessUserRegistration(ctx context.Context, db *ent.Client, youTubeConfig *oauth2.Config, userID uint64) error {
	crs, err := membership.CheckForUser(ctx, db, youTubeConfig, userID, nil)
	if err != nil {
		return fmt.Errorf("error checking memberships for user: %w", err)
	}
	return membership.SaveMemberships(ctx, db, userID, crs)
}
