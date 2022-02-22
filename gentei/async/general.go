package async

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/membership"
	"github.com/rs/zerolog/log"
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

// ProcessYouTubeRegistration really only checks memberships and triggers changes. One day it might do something else?
func ProcessYouTubeRegistration(ctx context.Context, db *ent.Client, youTubeConfig *oauth2.Config, userID uint64) error {
	crs, err := membership.CheckForUser(ctx, db, youTubeConfig, userID, nil)
	if err != nil {
		return fmt.Errorf("error checking memberships for user: %w", err)
	}
	lastCheckTime := time.Now()
	err = membership.SaveMemberships(ctx, db, userID, crs)
	if err != nil {
		return fmt.Errorf("error saving memberships for user: %w", err)
	}
	log.Info().Uint64("userID", userID).Time("lastCheck", lastCheckTime).Msg("setting last check time")
	err = db.User.UpdateOneID(userID).
		SetLastCheck(time.Now()).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("error saving last check time for user: %w", err)
	}
	return err
}
