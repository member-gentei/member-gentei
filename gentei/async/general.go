package async

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/member-gentei/member-gentei/gentei/ent/usermembership"
	"github.com/member-gentei/member-gentei/gentei/membership"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

func PublishGeneralMessage(ctx context.Context, topic *pubsub.Topic, message GeneralPSMessage) error {
	return publishPubSubMessage(ctx, topic, GeneralType, message)
}

func PublishApplyMembershipMessage(ctx context.Context, topic *pubsub.Topic, message ApplyMembershipPSMessage) error {
	return publishPubSubMessage(ctx, topic, ApplyMembershipType, message)
}

// ProcessYouTubeRegistration really only checks memberships and triggers changes. One day it might do something else?
func ProcessYouTubeRegistration(ctx context.Context, db *ent.Client, youTubeConfig *oauth2.Config, userID uint64, setChangeReason func(string)) error {
	crs, err := membership.CheckForUser(ctx, db, youTubeConfig, userID, nil)
	if err != nil {
		return fmt.Errorf("error checking memberships for user: %w", err)
	}
	lastCheckTime := time.Now()
	setChangeReason("new user / YouTube channel change")
	err = membership.SaveMemberships(ctx, db, userID, crs)
	if err != nil {
		return fmt.Errorf("error saving memberships for user: %w", err)
	}
	log.Info().Str("userID", strconv.FormatUint(userID, 10)).Time("lastCheck", lastCheckTime).Msg("setting last check time")
	err = db.User.UpdateOneID(userID).
		SetLastCheck(time.Now()).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("error saving last check time for user: %w", err)
	}
	return err
}

// ProcessUserDelete revokes tokens and tells the bot to delete the user.
// The bot has to delete the user because it'll communicate the final role removals + user deletion, and at that point why have the async queue require another message?
func ProcessUserDelete(ctx context.Context, db *ent.Client, topic *pubsub.Topic, userID uint64, reason string) error {
	logger := log.With().Str("userID", strconv.FormatUint(userID, 10)).Logger()
	// revoke tokens
	u, err := db.User.Get(ctx, userID)
	if err != nil {
		// already gone!
		if ent.IsNotFound(err) {
			return nil
		}
		return err
	}
	if u.DiscordToken != nil {
		err = revokeDiscordToken(ctx, u.DiscordToken)
		if err != nil {
			logger.Err(err).Msg("error revoking Discord token, proceeding to delete")
		}
		err = nil
	}
	if u.YoutubeToken != nil {
		err = revokeYouTubeToken(ctx, u.YoutubeToken)
		if err != nil {
			logger.Err(err).Msg("error revoking YouTube token, proceeding to delete")
		}
		err = nil
	}
	// tell the bot to delete the user
	err = PublishApplyMembershipMessage(ctx, topic, ApplyMembershipPSMessage{
		DeleteSingle: &DeleteUserMessage{
			UserID: json.Number(strconv.FormatUint(userID, 10)),
			Reason: reason,
		},
	})
	if err != nil {
		return fmt.Errorf("error publishing role revoke message: %w", err)
	}
	return nil
}

// ProcessUserDelete revokes tokens and deletes all membership edges.
func ProcessYouTubeDelete(ctx context.Context, db *ent.Client, topic *pubsub.Topic, userID uint64) error {
	logger := log.With().Str("userID", strconv.FormatUint(userID, 10)).Logger()
	u, err := db.User.Get(ctx, userID)
	if err != nil {
		if !ent.IsNotFound(err) {
			return err
		}
	}
	// no-op
	if u == nil || u.YoutubeID == nil || *u.YoutubeID == "" {
		return nil
	}
	err = revokeYouTubeToken(ctx, u.YoutubeToken)
	if err != nil {
		logger.Err(err).Msg("error revoking YouTube token, proceeding to delete")
	}
	// delete UserMembership edges first
	deleteCount, err := db.UserMembership.Delete().
		Where(
			usermembership.HasUserWith(user.ID(userID)),
		).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("error deleting UserMembership objects: %w", err)
	}
	logger.Info().
		Int("count", deleteCount).
		Msg("deleted UserMembership objects associated with YouTube user")
	return db.User.UpdateOneID(u.ID).
		ClearYoutubeID().
		ClearYoutubeToken().
		Exec(ctx)
}

func revokeYouTubeToken(ctx context.Context, token *oauth2.Token) error {
	var toRevoke string
	if time.Since(token.Expiry) > 0 {
		toRevoke = token.RefreshToken
	} else {
		toRevoke = token.AccessToken
	}
	// why does google just decide to put this in params instead of the body
	// https://developers.google.com/identity/protocols/oauth2/web-server#tokenrevoke
	r, err := http.Post(
		fmt.Sprintf("https://oauth2.googleapis.com/revoke?token=%s", toRevoke),
		"application/x-www-form-urlencoded",
		nil,
	)
	if err != nil {
		return err
	}
	// 400 error happens if the token was already revoked by a user
	if r.StatusCode >= 400 {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		var jbody struct {
			Error       string `json:"error"`
			Description string `json:"error_description"`
		}
		json.Unmarshal(body, &jbody)
		if jbody.Description == "Token expired or revoked" {
			return nil
		}
		if jbody.Error == "invalid_token" {
			return nil
		}
		log.Error().
			Int("status", r.StatusCode).
			Str("body", string(body)).
			Msg(">=400 status code revoking YouTube token")
	}
	return nil
}

var revokeMutex sync.Mutex

func revokeDiscordToken(ctx context.Context, token *oauth2.Token) error {
	// rudimentary rate limit
	revokeMutex.Lock()
	defer revokeMutex.Unlock()
	var (
		toRevoke string
		values   = url.Values{}
	)
	if time.Since(token.Expiry) > 0 {
		toRevoke = token.AccessToken
	} else {
		toRevoke = token.RefreshToken
	}
	values.Add("token", toRevoke)
	r, err := retryablehttp.PostForm("https://discord.com/api/oauth2/token/revoke", values)
	if err != nil {
		return err
	}
	// 400 error happens if the token was already revoked by a user.
	if r.StatusCode >= 400 {
		body, _ := io.ReadAll(r.Body)
		if r.StatusCode == http.StatusUnauthorized {
			// 401 happens when the token is already revoked/expired
			if strings.Contains(string(body), "invalid_client") {
				return nil
			}
		}
		defer r.Body.Close()
		log.Error().
			Int("status", r.StatusCode).
			Str("body", string(body)).
			Msg(">=400 status code revoking Discord token")
	}
	return nil
}
