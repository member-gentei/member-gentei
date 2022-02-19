package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/bwmarrin/discordgo"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/member-gentei/member-gentei/gentei/async"
	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/member-gentei/member-gentei/gentei/ent/guildrole"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
	"github.com/rs/zerolog/log"
	"github.com/ziflex/lecho/v3"
	"golang.org/x/oauth2"
)

func ServeAPI(db *ent.Client, discordConfig *oauth2.Config, youTubeConfig *oauth2.Config, topic *pubsub.Topic, jwtKey []byte, address string, debug bool) error {
	// create a copy of discordConfig that has the enroll endpoint
	enrollDiscordConfig := *discordConfig
	enrollDiscordConfig.RedirectURL = strings.Replace(discordConfig.RedirectURL, "login/discord", "app/enroll", 1)
	e := echo.New()
	e.Debug = debug
	xffExtract := echo.ExtractIPFromXFFHeader()
	e.IPExtractor = func(r *http.Request) string {
		cfcip := r.Header.Get("CF-Connecting-IP")
		if cfcip != "" {
			return cfcip
		}
		// fall back to x-forwarded-for, which is what the ingress writes
		return xffExtract(r)
	}
	// configure CORS
	corsConfig := middleware.CORSConfig{
		AllowOrigins:     []string{"https://gentei.tindabox.net", "https://member-gentei.tindabox.net"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Cookie", "Content-Type"},
	}
	if strings.Contains(address, "localhost:") {
		corsConfig.AllowOrigins = append(corsConfig.AllowOrigins, "http://localhost:3000")
		log.Debug().Interface("allowOrigins", corsConfig.AllowOrigins).Msg("CORS modified for local use")
	}
	// don't log healthz requests
	e.GET("/healthz", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	g := e.Group("")
	g.Use(lecho.Middleware(
		lecho.Config{Logger: lecho.From(log.Logger, lecho.WithLevel(2))},
	))
	g.Use(middleware.CORSWithConfig(corsConfig))
	g.POST(
		"/login/discord",
		loginDiscord(db, discordConfig, jwtKey, topic, !strings.Contains(address, "localhost:")),
	)
	loginRequired := g.Group("")
	loginRequired.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		TokenLookup: "cookie:token",
		SigningKey:  jwtKey,
		Claims:      &jwt.StandardClaims{},
	}))
	loginRequired.POST("/login/youtube", loginYouTube(db, youTubeConfig))
	loginRequired.POST("/logout", logout())
	loginRequired.GET("/me", getMe(db))
	loginRequired.DELETE("/me/youtube", deleteYouTube(db))
	loginRequired.DELETE("/me", deleteMe(db))
	loginRequired.POST("/enroll-guild", enrollGuild(db, &enrollDiscordConfig))
	loginRequired.GET("/guild/:id", getGuild(db))
	loginRequired.PATCH("/guild/:id", patchGuild(db))
	loginRequired.GET("/talents", getTalents(db))
	return e.Start(address)
}

type loginDiscordData struct {
	Code string `json:"code"`
}

func loginDiscord(
	db *ent.Client,
	discordConfig *oauth2.Config,
	jwtKey []byte,
	topic *pubsub.Topic,
	secureCookie bool,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			data loginDiscordData
			ctx  = c.Request().Context()
		)
		if err := c.Bind(&data); err != nil {
			return err
		}
		if data.Code == "" {
			return errors.New("must specify code")
		}
		oauthToken, err := discordConfig.Exchange(ctx, data.Code)
		var (
			retErr *oauth2.RetrieveError
		)
		if errors.As(err, &retErr) {
			log.Err(err).Msg("oauth2.RetrieveError")
			var body struct {
				Error       string `json:"error"`
				Description string `json:"error_description"`
			}
			err = json.Unmarshal(retErr.Body, &body)
			if err != nil {
				return fmt.Errorf("error decoding Discord OAuth login repsonse: %w", err)
			}
			return c.JSON(http.StatusBadRequest, body)
		} else if err != nil {
			log.Err(err).Msgf("concrete type: %T", err)
			return err
		}
		svc, err := discordgo.New(fmt.Sprintf("Bearer %s", oauthToken.AccessToken))
		if err != nil {
			return err
		}
		discordUser, err := svc.User("@me")
		if err != nil {
			return err
		}
		// create token
		expiry := time.Now().Add(time.Hour * 24 * 14)
		token := jwt.NewWithClaims(
			jwt.SigningMethodHS256,
			&jwt.StandardClaims{
				Id:        discordUser.ID,
				Audience:  "https://gentei.tindabox.net",
				ExpiresAt: expiry.Unix(),
			},
		)
		tokenStr, err := token.SignedString(jwtKey)
		if err != nil {
			log.Err(err).Msgf("%T", err)
			return err
		}
		// save user to db
		userID, err := strconv.ParseUint(discordUser.ID, 10, 64)
		if err != nil {
			return err
		}
		// on create, we check all enrolled servers for privs and presence
		isUpdate, err := db.User.Query().Where(user.ID(userID)).Exist(ctx)
		if err != nil {
			return err
		}
		userDBID, err := db.User.Create().
			SetID(userID).
			SetFullName(fmt.Sprintf("%s#%s", discordUser.Username, discordUser.Discriminator)).
			SetAvatarHash(discordUser.Avatar).
			SetDiscordToken(oauthToken).
			OnConflictColumns(user.FieldID).
			UpdateFullName().UpdateAvatarHash().
			UpdateDiscordToken().
			ID(ctx)
		if err != nil {
			return err
		}
		if !isUpdate {
			userGuilds, err := svc.UserGuilds(0, "", "")
			if err != nil {
				return err
			}
			var userGuildIDs = make([]uint64, len(userGuilds))
			for i, ug := range userGuilds {
				guildID, err := strconv.ParseUint(ug.ID, 10, 64)
				if err != nil {
					return err
				}
				userGuildIDs[i] = guildID
			}
			// link guild members
			guildIDs, err := db.Guild.Query().Where(
				guild.IDIn(userGuildIDs...),
			).IDs(ctx)
			if err != nil {
				return err
			}
			log.Debug().
				Uints64("guildIDs", userGuildIDs).
				Uints64("ourGuildIDs", guildIDs).
				Msg("user guilds")
			// link admins
			adminIDs, err := db.Guild.Query().Where(func(s *sql.Selector) {
				s.Where(sqljson.ValueContains(guild.FieldAdminSnowflakes, userID))
			}).IDs(ctx)
			if err != nil {
				return err
			}
			err = db.User.UpdateOneID(userDBID).
				AddGuildIDs(guildIDs...).
				AddGuildsAdminIDs(adminIDs...).
				Exec(ctx)
			if err != nil {
				return err
			}
			// the rest can happen later.
			if topic == nil {
				log.Warn().Uint64("userID", userID).
					Msg("async pubsub topic unspecified, would've sent general message")
			} else {
				err = async.PublishGeneralMessage(ctx, topic, async.GeneralPSMessage{
					UserRegistration: &async.UserRegistrationMessage{
						UserID: json.Number(discordUser.ID),
					},
				})
				if err != nil {
					return err
				}
			}
		}
		c.SetCookie(&http.Cookie{
			Name:     "token",
			Value:    tokenStr,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
			Secure:   secureCookie,
			Expires:  expiry,
			HttpOnly: true,
		})
		me, err := meResponseFromUser(
			db.User.Query().
				Where(user.ID(userDBID)).
				WithGuilds().
				WithGuildsAdmin().
				WithMemberships(func(umq *ent.UserMembershipQuery) {
					umq.WithRoles()
				}).
				OnlyX(ctx),
		)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusAccepted, me)
	}
}

func logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.SetCookie(&http.Cookie{
			Name:    "token",
			Value:   "delete-this",
			Path:    "/",
			Expires: time.Now().Add(-time.Hour),
			MaxAge:  0,
		})
		return c.JSON(http.StatusAccepted, nil)
	}
}

type loginYouTubeData struct {
	Code string `json:"code"`
}

type loginYouTubeResponse struct {
	ChannelID string
}

func loginYouTube(db *ent.Client, youtubeConfig *oauth2.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			ctx     = c.Request().Context()
			jwtUser = c.Get("user").(*jwt.Token)
			claims  = jwtUser.Claims.(*jwt.StandardClaims)
			data    loginYouTubeData
		)
		err := c.Bind(&data)
		if err != nil {
			return err
		}
		userID, err := strconv.ParseUint(claims.Id, 10, 64)
		if err != nil {
			return err
		}
		logger := log.With().Uint64("userID", userID).Logger()
		token, err := youtubeConfig.Exchange(ctx, data.Code)
		var (
			retErr *oauth2.RetrieveError
		)
		if errors.As(err, &retErr) {
			logger.Err(err).Msg("oauth2.RetrieveError")
			var body struct {
				Error       string `json:"error"`
				Description string `json:"error_description"`
			}
			err = json.Unmarshal(retErr.Body, &body)
			if err != nil {
				return fmt.Errorf("error decoding Discord OAuth login repsonse: %w", err)
			}
			return c.JSON(http.StatusBadRequest, body)
		} else if err != nil {
			logger.Err(err).Msgf("concrete type: %T", err)
			return err
		}
		// check if this YouTube channel is already associated with a different user
		ts := youtubeConfig.TokenSource(ctx, token)
		userChannelID, err := getChannelIDForYouTubeToken(ctx, ts)
		if err != nil {
			logger.Err(err).Msg("error getting channel ID with new token")
			return err
		}
		first, err := db.User.Query().Where(
			user.YoutubeID(userChannelID),
		).First(ctx)
		if ent.IsNotFound(err) {
			// great
		} else if err != nil {
			return err
		} else if first.ID != userID {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "YouTube channel belongs to a different user",
			})
		}
		// save
		token, err = ts.Token()
		if err != nil {
			return err
		}
		err = db.User.UpdateOneID(userID).
			SetYoutubeID(userChannelID).
			SetYoutubeToken(token).
			Exec(ctx)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, loginYouTubeResponse{
			ChannelID: userChannelID,
		})
	}
}

type meResponse struct {
	ID            string
	FullName      string
	AvatarHash    string
	YouTube       meResponseYouTube
	LastRefreshed int64
	Memberships   map[string]meResponseMembership `json:",omitempty"`
	ServerAdmin   []string                        `json:",omitempty"`
	Servers       []string                        `json:",omitempty"`
}

type meResponseMembership struct {
	LastVerified int64
	Failed       bool
}

type meResponseYouTube struct {
	ID    string
	Valid bool
}

func meResponseFromUser(user *ent.User) (meResponse, error) {
	yt := meResponseYouTube{
		Valid: user.YoutubeToken != nil,
	}
	if user.YoutubeID != nil {
		yt.ID = *user.YoutubeID
	}
	var (
		memberships map[string]meResponseMembership
		serverAdmin []string
		servers     []string
	)
	if len(user.Edges.GuildsAdmin) > 0 {
		for _, dg := range user.Edges.GuildsAdmin {
			serverAdmin = append(serverAdmin, strconv.FormatUint(dg.ID, 10))
		}
	}
	if len(user.Edges.Guilds) > 0 {
		for _, dg := range user.Edges.Guilds {
			servers = append(servers, strconv.FormatUint(dg.ID, 10))
		}
	}
	if len(user.Edges.Memberships) > 0 {
		memberships = make(map[string]meResponseMembership, len(user.Edges.Memberships))
		for _, membership := range user.Edges.Memberships {
			guildRoles, err := membership.Edges.RolesOrErr()
			if err != nil {
				return meResponse{}, err
			}
			for _, guildRole := range guildRoles {
				roleID := strconv.FormatUint(guildRole.ID, 10)
				memberships[roleID] = meResponseMembership{
					LastVerified: membership.LastVerified.Unix(),
				}
			}
		}
	}
	return meResponse{
		ID:            strconv.FormatUint(user.ID, 10),
		FullName:      user.FullName,
		AvatarHash:    user.AvatarHash,
		LastRefreshed: user.LastCheck.Unix(),
		YouTube:       yt,
		Memberships:   memberships,
		ServerAdmin:   serverAdmin,
		Servers:       servers,
	}, nil
}

func getMe(db *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			ctx     = c.Request().Context()
			jwtUser = c.Get("user").(*jwt.Token)
			claims  = jwtUser.Claims.(*jwt.StandardClaims)
		)
		userID, err := strconv.ParseUint(claims.Id, 10, 64)
		if err != nil {
			return err
		}
		// get user by ID
		u, err := db.User.Query().
			WithGuilds().
			WithGuildsAdmin().
			Where(user.ID(userID)).
			First(ctx)
		if ent.IsNotFound(err) {
			// this happens with multiple sessions
			c.SetCookie(&http.Cookie{
				Name:    "token",
				Value:   "delete-this",
				Path:    "/",
				Expires: time.Now().Add(-time.Hour),
				MaxAge:  0,
			})
			return c.NoContent(http.StatusUnauthorized)
		} else if err != nil {
			return err
		}
		// TODO: cache management on this response
		me, err := meResponseFromUser(u)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusAccepted, me)
	}
}

func deleteYouTube(db *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			ctx     = c.Request().Context()
			jwtUser = c.Get("user").(*jwt.Token)
			claims  = jwtUser.Claims.(*jwt.StandardClaims)
		)
		userID, err := strconv.ParseUint(claims.Id, 10, 64)
		if err != nil {
			return err
		}
		err = db.User.UpdateOneID(userID).
			ClearYoutubeID().
			ClearMemberships().
			Exec(ctx)
		if err != nil {
			return err
		}
		me, err := meResponseFromUser(
			db.User.Query().
				Where(user.ID(userID)).
				WithGuilds().
				WithGuildsAdmin().
				OnlyX(ctx),
		)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusAccepted, me)
	}
}

func deleteMe(db *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			ctx     = c.Request().Context()
			jwtUser = c.Get("user").(*jwt.Token)
			claims  = jwtUser.Claims.(*jwt.StandardClaims)
		)
		userID, err := strconv.ParseUint(claims.Id, 10, 64)
		if err != nil {
			return err
		}
		// revoke tokens
		u, err := db.User.Get(ctx, userID)
		if err != nil {
			return err
		}
		if u.DiscordToken != nil {
			err = revokeDiscordToken(ctx, u.DiscordToken)
			if err != nil {
				log.Err(err).Uint64("userID", userID).Msg("error revoking Discord token, proceeding to delete")
			}
			err = nil
		}
		if u.YoutubeToken != nil {
			err = revokeYouTubeToken(ctx, u.YoutubeToken)
			if err != nil {
				log.Err(err).Uint64("userID", userID).Msg("error revoking YouTube token, proceeding to delete")
			}
			err = nil
		}
		err = db.User.DeleteOneID(userID).Exec(ctx)
		if err != nil {
			return err
		}
		return logout()(c)
	}
}

type enrollGuildData struct {
	Code        string `json:"code"`
	Permissions string `json:"permissions"`
}

type guildResponse struct {
	ID            string
	Name          string
	Icon          string
	TalentIDs     []string `json:",omitempty"`
	AdminIDs      []string
	RolesByTalent map[string]roleInfo
}

type roleInfo struct {
	ID   string
	Name string
}

func enrollGuild(db *ent.Client, discordConfig *oauth2.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			data   enrollGuildData
			ctx    = c.Request().Context()
			user   = c.Get("user").(*jwt.Token)
			claims = user.Claims.(*jwt.StandardClaims)
		)
		if err := c.Bind(&data); err != nil {
			return err
		}
		if data.Code == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "must specify code",
			})
		}
		if data.Permissions == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "must specify permissions",
			})
		}
		oauthToken, err := discordConfig.Exchange(ctx, data.Code)
		var (
			retErr *oauth2.RetrieveError
		)
		if errors.As(err, &retErr) {
			log.Err(err).Msg("oauth2.RetrieveError")
			var body struct {
				Error       string `json:"error"`
				Description string `json:"error_description"`
			}
			err = json.Unmarshal(retErr.Body, &body)
			if err != nil {
				return fmt.Errorf("error decoding Discord OAuth login repsonse: %w", err)
			}
			return c.JSON(http.StatusBadRequest, body)
		} else if err != nil {
			log.Err(err).Msgf("concrete type: %T", err)
			return err
		}
		discordUser, err := getDiscordTokenMe(oauthToken)
		if err != nil {
			return err
		}
		// the user must match the userID we have in the JWT
		if discordUser.ID != claims.Id {
			// revoke! we don't want it!
			values := url.Values{}
			values.Set("client_id", discordConfig.ClientID)
			values.Set("client_secret", discordConfig.ClientSecret)
			values.Set("access_token", oauthToken.AccessToken)
			values.Set("refresh_token", oauthToken.RefreshToken)
			_, err := http.PostForm("https://discord.com/api/oauth2/token/revoke", values)
			if err != nil {
				log.Err(err).Msg("failed to revoke OAuth2 token")
			}
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "user used to add bot does not match logged-in user",
			})
		}
		// okay, *now* we can save it all
		guildMap := oauthToken.Extra("guild").(map[string]interface{})
		userID, err := strconv.ParseUint(claims.Id, 10, 64)
		if err != nil {
			return err
		}
		response, err := parseAndSaveGuild(ctx, db, userID, guildMap)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, response)
	}
}

type singleGuildData struct {
	ID uint64 `param:"id"`
}

func getGuild(db *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			data    singleGuildData
			ctx     = c.Request().Context()
			jwtUser = c.Get("user").(*jwt.Token)
			claims  = jwtUser.Claims.(*jwt.StandardClaims)
		)
		if err := c.Bind(&data); err != nil {
			return err
		}
		userID, err := strconv.ParseUint(claims.Id, 10, 64)
		if err != nil {
			return err
		}
		// only return the guild if the user has some association with it
		dg, err := db.Guild.Query().
			WithYoutubeTalents().
			WithRoles(func(grq *ent.GuildRoleQuery) { grq.WithTalent() }).
			WithAdmins().
			Where(
				guild.ID(data.ID),
				guild.Or(
					guild.HasMembersWith(user.ID(userID)),
					guild.HasAdminsWith(user.ID(userID)),
				),
			).First(ctx)
		if ent.IsNotFound(err) {
			return c.NoContent(http.StatusNotFound)
		} else if err != nil {
			return err
		}
		talentIDs := make([]string, len(dg.Edges.YoutubeTalents))
		for i := range dg.Edges.YoutubeTalents {
			talentIDs[i] = dg.Edges.YoutubeTalents[i].ID
		}
		return c.JSON(http.StatusOK, makeGuildResponse(dg))
	}
}

type patchGuildData struct {
	ID      uint64   `param:"id"`
	Talents []string `json:"talents"`
}

type patchGuildErrorResponse struct {
	Error patchGuildErrorResponseError `json:"error"`
}

type patchGuildErrorResponseError struct {
	Message string            `json:"message,omitempty"`
	Talents map[string]string `json:"talents,omitempty"`
}

const maxTalentCount = 16

func patchGuild(db *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			data    patchGuildData
			ctx     = c.Request().Context()
			jwtUser = c.Get("user").(*jwt.Token)
			claims  = jwtUser.Claims.(*jwt.StandardClaims)
		)
		if err := c.Bind(&data); err != nil {
			return err
		}
		userID, err := strconv.ParseUint(claims.Id, 10, 64)
		if err != nil {
			return err
		}
		// only allow PATCH if the user is an admin
		isAdmin, err := db.Guild.Query().
			Where(
				guild.ID(data.ID),
				guild.HasAdminsWith(user.ID(userID)),
			).
			Exist(ctx)
		if ent.IsNotFound(err) {
			log.Debug().Msg("guild not found")
			return c.NoContent(http.StatusForbidden)
		} else if err != nil {
			return err
		} else if !isAdmin {
			log.Debug().Msg("user is not an admin")
			return c.NoContent(http.StatusForbidden)
		}
		// check constraints
		if len(data.Talents) > maxTalentCount {
			return c.JSON(http.StatusBadRequest, patchGuildErrorResponse{
				Error: patchGuildErrorResponseError{
					Message: "servers can track a maximum of 16 channels",
				},
			})
		}
		// perform patch
		err = createOrAssociateTalentsToGuild(ctx, db, data.ID, data.Talents)
		var nmpErr ErrNoMembershipPlaylist
		if errors.As(err, &nmpErr) {
			return c.JSON(http.StatusBadRequest, patchGuildErrorResponse{
				Error: patchGuildErrorResponseError{
					Talents: map[string]string{
						nmpErr.ChannelID: "memberships not open",
					},
				},
			})
		}
		if err != nil {
			return err
		}
		dg, err := db.Guild.Query().
			WithYoutubeTalents().
			WithRoles(func(grq *ent.GuildRoleQuery) { grq.WithTalent() }).
			Where(guild.ID(data.ID)).
			First(ctx)
		if err != nil {
			return err
		}
		// remove talents and any associated roles
		var (
			existingTalentMap = make(map[string]uint64, len(dg.Edges.Roles))
			missingTalentIDs  []string
		)
		for _, role := range dg.Edges.Roles {
			existingTalentMap[role.Edges.Talent.ID] = role.ID
		}
		for _, talent := range dg.Edges.YoutubeTalents {
			existingTalentMap[talent.ID] = 1
		}
		for _, talentID := range data.Talents {
			if existingTalentMap[talentID] == 0 {
				missingTalentIDs = append(missingTalentIDs, talentID)
			}
		}
		if len(missingTalentIDs) > 0 {
			update := db.Guild.UpdateOneID(data.ID).
				RemoveYoutubeTalentIDs(missingTalentIDs...)
			removeRoleIDs, err := db.Guild.QueryRoles(dg).
				Where(guildrole.HasTalentWith(
					youtubetalent.IDIn(missingTalentIDs...),
				)).
				IDs(ctx)
			if err != nil {
				return err
			}
			if len(removeRoleIDs) > 0 {
				update = update.RemoveRoleIDs(removeRoleIDs...)
			}
			err = update.Exec(ctx)
			if err != nil {
				return err
			}
		}
		return c.JSON(http.StatusOK, makeGuildResponse(dg))
	}
}

type talentResponseItem struct {
	ID        string
	Name      string
	Thumbnail string
}

func getTalents(db *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			ctx = c.Request().Context()
		)
		talents, err := db.YouTubeTalent.Query().
			Order(ent.Desc(youtubetalent.FieldLastUpdated)).
			All(ctx)
		if err != nil {
			return err
		}
		talentItems := make([]talentResponseItem, 0, len(talents))
		for _, talent := range talents {
			talentItems = append(talentItems, talentResponseItem{
				ID:        talent.ID,
				Name:      talent.ChannelName,
				Thumbnail: talent.ThumbnailURL,
			})
		}
		return c.JSON(http.StatusOK, talentItems)
	}
}
