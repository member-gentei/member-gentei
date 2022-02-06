package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func getDiscordTokenMe(token *oauth2.Token) (*discordgo.User, error) {
	dg, err := discordgo.New(fmt.Sprintf("Bearer %s", token.AccessToken))
	if err != nil {
		return nil, err
	}
	return dg.User("@me")
}

type tokenEmbeddedGuild struct {
	ID      string
	Name    string
	Icon    string
	OwnerID string `mapstructure:"owner_id"`
}

func parseAndSaveGuild(ctx context.Context, db *ent.Client, userID uint64, guildMap map[string]interface{}) (guildResponse, error) {
	var embed tokenEmbeddedGuild
	err := mapstructure.Decode(guildMap, &embed)
	if err != nil {
		return guildResponse{}, fmt.Errorf("error parsing embedded guild: %w", err)
	}
	guildID, err := strconv.ParseUint(embed.ID, 10, 64)
	if err != nil {
		return guildResponse{}, fmt.Errorf("error parsing embedded guild.ID as uint64: %w", err)
	}
	ownerID, err := strconv.ParseUint(embed.OwnerID, 10, 64)
	if err != nil {
		return guildResponse{}, fmt.Errorf("error parsing embedded guild.ID as uint64: %w", err)
	}
	// update guild - if the guild exists, we do not update the owner ID.
	// we don't know if the owner is a user, so don't do anything with this yet
	dg, err := db.Guild.Query().
		WithYoutubeTalents().
		Where(guild.ID(guildID)).
		First(ctx)
	if ent.IsNotFound(err) {
		var adminSnowflakes = []uint64{ownerID}
		if ownerID != userID {
			adminSnowflakes = append(adminSnowflakes, userID)
		}
		dg, err = db.Guild.Create().
			SetID(guildID).
			SetName(embed.Name).
			SetIconHash(embed.Icon).
			SetAdminSnowflakes(adminSnowflakes).
			AddAdminIDs(userID).
			Save(ctx)
		if err != nil {
			return guildResponse{}, fmt.Errorf("error creating Guild: %w", err)
		}
		log.Debug().Interface("guild", dg).Msg("created Guild")
	} else {
		// update edges and admin fields to include this and other user
		var connectingUserIsAdmin bool
		for _, snowflake := range dg.AdminSnowflakes {
			if snowflake == userID {
				connectingUserIsAdmin = true
				break
			}
		}
		if !connectingUserIsAdmin {
			dg.AdminSnowflakes = append(dg.AdminSnowflakes, userID)
			err = db.Guild.UpdateOneID(dg.ID).SetAdminSnowflakes(dg.AdminSnowflakes).Exec(ctx)
			if err != nil {
				return guildResponse{}, fmt.Errorf("error updating Guild admin snowflakes with enrolling user: %w", err)
			}
		}
		exists, err := db.Guild.QueryAdmins(dg).Where(user.ID(userID)).Exist(ctx)
		if err != nil {
			return guildResponse{}, fmt.Errorf("error querying Guild admins for enrolling user: %w", err)
		}
		if !exists {
			err = db.Guild.UpdateOneID(guildID).
				AddAdminIDs(userID).
				Exec(ctx)
			if err != nil {
				return guildResponse{}, fmt.Errorf("error updating Guild admin edges with enrolling user: %w", err)
			}
		}
	}
	// populate response
	talents := dg.Edges.YoutubeTalents
	talentIDs := make([]string, len(talents))
	for i := range talents {
		talentIDs[i] = talents[i].ID
	}
	return guildResponse{
		ID:        embed.ID,
		Name:      embed.Name,
		Icon:      embed.Icon,
		TalentIDs: talentIDs,
		Settings:  dg.Settings,
	}, nil
}

func createOrAssociateTalentsToGuild(ctx context.Context, db *ent.Client, guildID uint64, talentIDs []string) error {
	// create some placeholders if we don't have the channel on file
	existingTalentIDs, err := db.YouTubeTalent.Query().
		Where(youtubetalent.IDIn(talentIDs...)).
		IDs(ctx)
	if err != nil {
		return err
	}
	existingTalentsMap := make(map[string]bool, len(talentIDs))
	for _, talentID := range existingTalentIDs {
		existingTalentsMap[talentID] = true
	}
	for _, talentID := range talentIDs {
		if existingTalentsMap[talentID] {
			continue
		}
		err = UpsertYouTubeChannelID(ctx, db, talentID)
		if err != nil {
			return err
		}
	}
	tx, err := db.Tx(ctx)
	if err != nil {
		return err
	}
	err = func() error {
		// add nonexistent edges
		err := tx.YouTubeTalent.Update().
			Where(
				youtubetalent.IDIn(talentIDs...),
				youtubetalent.Not(
					youtubetalent.HasGuildsWith(guild.ID(guildID)),
				),
			).
			AddGuildIDs(guildID).
			Exec(ctx)
		if err != nil {
			return err
		}
		// remove others
		err = tx.YouTubeTalent.Update().
			Where(
				youtubetalent.IDNotIn(talentIDs...),
				youtubetalent.HasGuildsWith(guild.ID(guildID)),
			).
			RemoveGuildIDs(guildID).
			Exec(ctx)
		if err != nil {
			return err
		}
		return tx.Commit()
	}()
	if err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return fmt.Errorf("%w: %v", err, rerr)
		}
		return err
	}
	return nil
}

func makeGuildResponse(dg *ent.Guild) guildResponse {
	var talentIDs []string
	if len(dg.Edges.YoutubeTalents) > 0 {
		talentIDs = make([]string, 0, len(dg.Edges.YoutubeTalents))
		for _, t := range dg.Edges.YoutubeTalents {
			talentIDs = append(talentIDs, t.ID)
		}
	}
	var adminIDs []string
	for _, id := range dg.AdminSnowflakes {
		adminIDs = append(adminIDs, strconv.FormatUint(id, 10))
	}
	return guildResponse{
		ID:        strconv.FormatUint(dg.ID, 10),
		Name:      dg.Name,
		Icon:      dg.IconHash,
		TalentIDs: talentIDs,
		AdminIDs:  adminIDs,
		Settings:  dg.Settings,
	}
}

func getChannelIDForYouTubeToken(ctx context.Context, ts oauth2.TokenSource) (string, error) {
	svc, err := youtube.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return "", err
	}
	clr, err := svc.Channels.List([]string{"id"}).Mine(true).Do()
	if err != nil {
		return "", err
	}
	if len(clr.Items) == 0 {
		log.Err(err).Msg("new token cannot get own channel ID")
		return "", err
	}
	return clr.Items[0].Id, nil
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
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		var jbody struct {
			Error       string `json:"error"`
			Description string `json:"error_description"`
		}
		json.Unmarshal(body, &jbody)
		if jbody.Description == "Token expired or revoked" {
			return nil
		}
		log.Error().
			Int("status", r.StatusCode).
			Str("body", string(body)).
			Msg(">=400 status code revoking YouTube token")
	}
	return nil
}

func revokeDiscordToken(ctx context.Context, token *oauth2.Token) error {
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
	r, err := http.PostForm("https://discord.com/api/oauth2/token/revoke", values)
	if err != nil {
		return err
	}
	// 400 error happens if the token was already revoked by a user.
	// TODO: actually code in the check after someone does this
	if r.StatusCode >= 400 {
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		log.Error().
			Int("status", r.StatusCode).
			Str("body", string(body)).
			Msg(">=400 status code revoking Discord token")
	}
	return nil
}
