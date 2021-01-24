package discord

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/member-gentei/member-gentei/bot/discord/api"
	"github.com/member-gentei/member-gentei/bot/discord/lang"
	"github.com/member-gentei/member-gentei/pkg/common"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"github.com/Lukaesebrot/dgc"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// discordBot is the whole Discord bot.
type discordBot struct {
	ctx       context.Context
	apiClient api.ClientWithResponsesInterface
	dgSession *discordgo.Session
	fs        *firestore.Client
	bundle    *i18n.Bundle

	lastMemberCheck                map[string]time.Time  // global rate limiter for user member checks
	guildStates                    map[string]guildState // holds state for a Discord guild
	ytChannels                     map[string]common.Channel
	ytChannelsMutex                sync.RWMutex
	ytChannelMembershipsLastLoaded map[string]time.Time
	ytChannelMemberships           map[string]map[string]struct{} // holds memberships corresponding to a particular YouTube channel
	ytChannelMembershipsMutex      sync.RWMutex

	// newMemberRoleApplier() stuff
	// key is "guildID-userID"
	// map[string](chan *discordgo.GuildMemberUpdate)
	guildMemberUpdateChannels sync.Map
}

// error returns the result of the first round of listening to changes.
func (d *discordBot) listenToGuildAssociations() error {
	var (
		firstErrChan = make(chan error)
		firstErrSent bool
	)
	d.guildStates = map[string]guildState{}
	go func() {
		snapsIter := d.fs.Collection(common.DiscordGuildCollection).Snapshots(d.ctx)
		for {
			snaps, err := snapsIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				if c := status.Code(err); c != codes.OK {
					log.Fatal().Err(err).Msg("rpc error getting DiscordGuild snapshots")
				}
				log.Err(err).Msg("error getting DiscordGuild snapshots, still listening")
				continue
			}
			log.Debug().Interface("changes", snaps.Changes).
				Msgf("At %s there were %d results", snaps.ReadTime, snaps.Size)
			for _, change := range snaps.Changes {
				var guild common.DiscordGuild
				err = change.Doc.DataTo(&guild)
				if err != nil {
					log.Err(err).Msg("error unmarshalling DiscordGuild")
				}
				// removal from this map deactivates the bot for this guild
				if change.Kind == firestore.DocumentRemoved {
					delete(d.guildStates, guild.ID)
					continue
				}
				log.Debug().Interface("guild", guild).Msg("received new Guild info")
				switch state := d.guildStates[guild.ID]; state.LoadState {
				case guildFirstEncounter:
					// totally new guild
					d.guildStates[guild.ID] = guildState{
						Doc:       guild,
						LoadState: guildWaitingForCreateEvent,
						localizer: makeLocalizer(d.bundle, guild.BCP47),
					}
				case guildWaitingForAssociationData, guildLoaded:
					var shouldCheckRoles bool
					if state.LoadState == guildWaitingForAssociationData {
						shouldCheckRoles = true
					} else {
						// only check roles if mappings change
						incomingMemberInfo := guildState{Doc: guild}.GetMembershipInfo()
						currentMemberInfo := state.GetMembershipInfo()
						if len(incomingMemberInfo) != len(currentMemberInfo) {
							shouldCheckRoles = true
						} else {
							for channelSlug, memberRoleID := range state.GetMembershipInfo() {
								if incomingMemberInfo[channelSlug] != memberRoleID {
									shouldCheckRoles = true
									break
								}
							}
						}
					}
					// only change the localizer if the language changes
					if state.Doc.BCP47 != guild.BCP47 {
						log.Debug().Str("bcp47", guild.BCP47).Msg("changing language for guild")
						state.localizer = makeLocalizer(d.bundle, guild.BCP47)
					}
					state.Doc = guild
					d.guildStates[guild.ID] = state
					if shouldCheckRoles {
						d.checkRoles(d.dgSession, guild.ID, nil)
						d.loadMemberships(guild.ID)
						err := d.enforceMembershipsAsync(guild.ID)
						if err != nil {
							log.Err(err).Str("guildID", guild.ID).Msg("error requesting guild members for async enforcement")
						}
					}
					state.LoadState = guildLoaded
					d.guildStates[guild.ID] = state
				}
			}
			if !firstErrSent {
				firstErrChan <- err
				close(firstErrChan)
				firstErrSent = true
			}
		}
	}()
	return <-firstErrChan
}

func (d *discordBot) listenToChannelChanges() error {
	var (
		firstErrChan = make(chan error)
		firstErrSent bool
	)
	go func() {
		snapsIter := d.fs.Collection(common.ChannelCollection).Snapshots(d.ctx)
		for {
			snaps, err := snapsIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				if c := status.Code(err); c != codes.OK {
					log.Fatal().Err(err).Msg("rpc error getting DiscordGuild snapshots")
				}
				log.Err(err).Msg("error getting DiscordGuild snapshots, still listening")
				continue
			}
			func() {
				d.ytChannelsMutex.Lock()
				defer d.ytChannelsMutex.Unlock()
				for _, change := range snaps.Changes {
					var channel common.Channel
					err = change.Doc.DataTo(&channel)
					if err != nil {
						log.Err(err).Msg("error unmarshalling DiscordGuild")
						break
					}
					if change.Kind == firestore.DocumentRemoved {
						delete(d.ytChannels, channel.ChannelID)
						continue
					}
					log.Debug().Interface("channel", channel).Msg("received new YouTube channel info")
					d.ytChannels[change.Doc.Ref.ID] = channel
				}
			}()
			if !firstErrSent {
				firstErrChan <- err
				close(firstErrChan)
				firstErrSent = true
			}
		}
	}()
	return <-firstErrChan
}

// reloads memberships for any present in guildState
func (d *discordBot) listenToMemberCheckUpdates(checkSubscription *pubsub.Subscription) {
	// this ring buffer is a mild O(4) defense against at-least-once delivery
	deliveredTimestamps := make([]string, 4)
	deliveredTimestampIndex := 0
	go func() {
		checkSubscription.Receive(d.ctx, func(ctx context.Context, msg *pubsub.Message) {
			ts := string(msg.Data)
			logger := log.With().Str("checkMessage", ts).Logger()
			logger.Debug().Msg("received member check message")
			msg.Ack()
			var noReload bool
			for _, storedTs := range deliveredTimestamps {
				if storedTs == ts {
					noReload = true
					break
				}
			}
			if noReload {
				logger.Debug().Msg("discarding duplicate member check message")
				return
			}
			deliveredTimestamps[deliveredTimestampIndex] = ts
			deliveredTimestampIndex = (deliveredTimestampIndex + 1) % 4
			// load all memberships in a mildly-threadsafe manner
			guildIDs := make([]string, 0, len(d.guildStates))
			for key := range d.guildStates {
				guildIDs = append(guildIDs, key)
			}
			logger.Info().Strs("guildIDs", guildIDs).Msg("check message received, reloading memberships")
			d.loadMemberships(guildIDs...)
			for _, guildID := range guildIDs {
				err := d.enforceMembershipsAsync(guildID)
				if err != nil {
					logger.Err(err).Str("guildID", guildID).Msg("error initiating enforcement after membership reload")
				}
			}
		})
	}()
}

func (d *discordBot) handleGuildCreate(s *discordgo.Session, g *discordgo.GuildCreate) {
	var (
		state  guildState
		exists bool
		logger = log.With().Str("guildID", g.ID).Str("name", g.Name).Logger()
	)
	// create guild if it doesn't exist
	if state, exists = d.guildStates[g.ID]; !exists {
		state = guildState{
			LoadState: guildWaitingForAssociationData,
			localizer: makeLocalizer(d.bundle, language.AmericanEnglish.String()),
		}
		d.guildStates[g.ID] = state
		logger.Info().Interface("state", state).Msg("guildWaitingForAssociationData")
		return
	}
	logger.Info().Interface("state", state).Msg("joined guild")
	d.checkRoles(s, g.ID, g.Guild)
	d.loadMemberships(g.ID)
	state.LoadState = guildLoaded
	d.guildStates[g.ID] = state
	err := d.enforceMembershipsAsync(g.ID)
	if err != nil {
		log.Err(err).Str("guildID", g.ID).Msg("error requesting guild members for async enforcement")
	}
	return
}

// usually initiated by a enforceMembershipsAsync call.
func (d *discordBot) handleGuildMembersChunk(s *discordgo.Session, chunk *discordgo.GuildMembersChunk) {
	logger := log.With().Str("guildID", chunk.GuildID).Int("chunkIndex", chunk.ChunkIndex).Logger()
	state, exists := d.guildStates[chunk.GuildID]
	if !exists || state.LoadState != guildLoaded {
		logger.Warn().Int("loadState", int(state.LoadState)).
			Msg("received GuildMembersChunk for non-ready GuildState")
		return
	}
	memberInfo := state.GetMembershipInfo()
	if len(memberInfo) == 0 {
		logger.Warn().Int("loadState", int(state.LoadState)).
			Msg("received GuildMembersChunk for guild without membership mappings")
		return
	}
	d.ytChannelMembershipsMutex.RLock()
	defer d.ytChannelMembershipsMutex.RUnlock()
	for channelSlug := range memberInfo {
		logger.Debug().Str("channelSlug", channelSlug).Msg("enforcing member role for chunk")
		verifiedMembers := d.ytChannelMemberships[channelSlug]
		d.enforceRole(chunk.GuildID, channelSlug, verifiedMembers, chunk.Members, len(memberInfo) > 1)
	}
	if chunk.ChunkIndex == chunk.ChunkCount-1 {
		refreshTime := time.Now()
		_, err := d.fs.Collection(common.DiscordGuildCollection).Doc(state.Doc.ID).Update(d.ctx, []firestore.Update{
			firestore.Update{
				Path:  "LastMembershipRefresh",
				Value: refreshTime,
			},
		})
		if err != nil {
			logger.Err(err).Msg("error updating membership refresh time in Firestore")
		}
		logger.Info().Str("refreshTime", refreshTime.Format(time.RFC3339)).Msg("membership refresh complete")
	}
}

func (d *discordBot) enforceRole(
	guildID string,
	channelSlug string,
	verifiedIDs map[string]struct{},
	guildMembers []*discordgo.Member,
	reasonHasName bool,
) {
	roleID := d.guildStates[guildID].GetMembershipRoleID(channelSlug)
	for _, guildMember := range guildMembers {
		_, shouldHaveRole := verifiedIDs[guildMember.User.ID]
		var reason string
		if reasonHasName {
			reason = fmt.Sprintf("periodic membership refresh (%s)", channelSlug)
		} else {
			reason = "periodic membership refresh"
		}
		if shouldHaveRole && !userHasRole(guildMember, roleID) {
			// user needs role
			d.newRoleApplier(
				guildID, channelSlug, guildMember.User, roleAdd, reason,
				5, defaultRoleApplyPeriod, defaultRoleApplyTimeout,
			)
		} else if !shouldHaveRole && userHasRole(guildMember, roleID) {
			// user needs role removed
			d.newRoleApplier(
				guildID, channelSlug, guildMember.User, roleRevoke, reason,
				5, defaultRoleApplyPeriod, defaultRoleApplyTimeout,
			)
		}
	}
}

func (d *discordBot) handleGuildMemberUpdate(s *discordgo.Session, update *discordgo.GuildMemberUpdate) {
	for _, roleID := range update.Member.Roles {
		updateKey := fmt.Sprintf("%s-%s-%s", update.GuildID, update.User.ID, roleID)
		value, exists := d.guildMemberUpdateChannels.Load(updateKey)
		if exists {
			updateChannel := value.(chan *discordgo.GuildMemberUpdate)
			updateChannel <- update
			log.Debug().Str("updateKey", updateKey).Msg("sending candidate guild member update")
		}
	}
}

func (d *discordBot) handleCmdCheck(ctx *dgc.Ctx) {
	var (
		m          = ctx.Event
		guildState = d.guildStates[m.GuildID]
	)
	switch guildState.LoadState {
	case guildWaitingForAssociationData:
		msg := mustLocalizeMessage(guildState.localizer, &i18n.Message{
			ID:    "GuildNotRegisteredReply",
			Other: "This Discord server isn't registered for membership management yet. Please wait until the server owner gets this sorted out!",
		})
		d.reply(
			log.Logger, m.GuildID, ctx.Event.ChannelID, ctx.Event.Author.ID,
			ctx.Event.Reference(), msg,
		)
	case guildWaitingForCreateEvent:
		msg := mustLocalizeMessage(guildState.localizer, &i18n.Message{
			ID:    "BotRestartedReply",
			Other: "This bot has secretly, recently restarted and is still loading - please try again in a minute!",
		})
		d.reply(
			log.Logger, m.GuildID, ctx.Event.ChannelID, ctx.Event.Author.ID,
			ctx.Event.Reference(), msg,
		)
	case guildLoaded:
		logger := log.With().Str("userID", m.Author.ID).Str("guildID", m.GuildID).Logger()
		if timeout := time.Now().Sub(d.lastMemberCheck[m.Author.ID]).Seconds(); timeout < 30 {
			logger.Debug().Float64("timeout", timeout).Msg("rate limited membership check")
			msg := mustLocalizeMessage(guildState.localizer, &i18n.Message{
				ID:    "RateLimitReply",
				Other: "Your replies are rate limited to prevent abuse - please try again in a minute!",
			})
			err := d.reply(
				logger, m.GuildID, ctx.Event.ChannelID, ctx.Event.Author.ID,
				ctx.Event.Reference(), msg,
			)
			if err != nil {
				logger.Err(err).Msg("error communicating rate limit")
				return
			}
			break
		}
		d.checkMembershipReply(logger, guildState, m)
	default:
		log.Debug().Interface("state", guildState).Msg("unsolicited message")
	}
}

type discordRoleCheck struct {
	channelSlug string
	action      roleAction
	required    bool
}

const multiMembersConfirmed = `Memberships confirmed! You will be granted roles corresponding to the following channels:
{{- range .titles }}
◦ ` + "`" + "{{ . }}" + "`" + `
{{- end }}`

func (d *discordBot) checkMembershipReply(
	logger zerolog.Logger,
	state guildState,
	m *discordgo.MessageCreate,
) {
	// send typing status as a loading indicator
	if err := d.dgSession.ChannelTyping(m.ChannelID); err != nil {
		logger.Err(err).Msg("error sending ChannelTyping status")
		return
	}
	var (
		membershipInfo = state.GetMembershipInfo()
		checks         = make([]discordRoleCheck, 0, len(membershipInfo))
	)
	for channelSlug, roleID := range membershipInfo {
		checkLogger := logger.With().Str("channelSlug", channelSlug).Logger()
		checkLogger.Debug().Str("channelSlug", channelSlug).Msg("checking membership for user")
		response, err := d.apiClient.CheckMembershipWithResponse(
			d.ctx,
			api.ChannelSlugPathParam(channelSlug),
			api.CheckMembershipJSONRequestBody{Snowflake: m.Author.ID},
		)
		if err != nil {
			checkLogger.Err(err).Msg("error checking user membership")
			return
		}
		d.lastMemberCheck[m.Author.ID] = time.Now()
		if response.JSON200 != nil {
			var changeRequired = false
			if response.JSON200.Member {
				checkLogger.Info().Msg("membership confirmed")
				// make change if role is not yet assigned
				if !userHasRole(m.Member, roleID) {
					changeRequired = true
				}
				checks = append(checks, discordRoleCheck{
					channelSlug: channelSlug,
					action:      roleAdd,
					required:    changeRequired,
				})
			} else {
				reason := *response.JSON200.Reason
				checkLogger.Debug().Str("reason", reason).Msg("user is not a member")
				if reason == "not connected" {
					msg := mustLocalizeMessage(state.localizer, &i18n.Message{
						ID:    "SignupRequiredReply",
						Other: "Please sign up on https://member-gentei.tindabox.net/app and run this command a few minutes after connecting your YouTube account!",
					})
					err = d.reply(checkLogger, m.GuildID, m.ChannelID, m.Author.ID, m.Reference(), msg)
					if err != nil {
						checkLogger.Err(err).Msg("error replying")
					}
					return
				}
				// make change if role is assigned
				if userHasRole(m.Member, roleID) {
					changeRequired = true
				}
				checks = append(checks, discordRoleCheck{
					channelSlug: channelSlug,
					action:      roleRevoke,
					required:    changeRequired,
				})
			}
		} else {
			checkLogger.Warn().Bytes("body", response.Body).Interface("headers", response.HTTPResponse.Header).Int("code", response.StatusCode()).Msg("json200 is nil")
			msg := mustLocalizeMessage(state.localizer, &i18n.Message{
				ID:    "ErrorCheckingReply",
				Other: "Error checking memberships! Please try again later - an alert has been sent to the developer.",
			})
			err = d.reply(checkLogger, m.GuildID, m.ChannelID, m.Author.ID, m.Reference(), msg)
			if err != nil {
				checkLogger.Err(err).Msg("error replying")
			}
			return
		}
	}
	// sort for consistency
	sort.Slice(checks, func(i, j int) bool {
		return checks[i].channelSlug < checks[j].channelSlug
	})
	var (
		confirmedChannelTitles = make([]string, 0, len(checks))
		unconfirmedCount       int
	)
	// perform required role changes`
	d.ytChannelsMutex.RLock()
	defer d.ytChannelsMutex.RUnlock()
	for _, check := range checks {
		channel := d.ytChannels[check.channelSlug]
		if check.action == roleAdd {
			confirmedChannelTitles = append(confirmedChannelTitles, channel.ChannelTitle)
		}
		if !check.required {
			continue
		}
		if check.action == roleRevoke {
			unconfirmedCount++
		}
		var actionReason string
		if check.action == roleAdd {
			actionReason = "`!mg check` verified"
		} else {
			actionReason = "`!mg check` un-verified"
		}
		if len(checks) > 1 {
			actionReason = fmt.Sprintf("%s (%s)", actionReason, channel.ChannelTitle)
		}
		d.newRoleApplier(
			m.GuildID, check.channelSlug, m.Author, check.action, actionReason,
			5, defaultRoleApplyPeriod, defaultRoleApplyTimeout,
		)
	}
	// reply!
	var replyMessage string
	if len(confirmedChannelTitles) > 0 {
		replyMessage = state.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "MembershipConfirmedReply",
				One:   "Membership confirmed! You will be added as a member shortly.",
				Other: multiMembersConfirmed,
			},
			TemplateData: map[string]interface{}{
				"titles": confirmedChannelTitles,
			},
			PluralCount: len(membershipInfo), // print the plural form when there are multiple membership possibilities
		})
	} else {
		replyMessage = mustLocalizeMessage(state.localizer, &i18n.Message{
			ID:    "MembershipUnconfirmedReply",
			Other: "We just checked, and you don't seem to be a member.",
		})
	}
	err := d.reply(
		logger, m.GuildID, m.ChannelID, m.Author.ID,
		m.Reference(), replyMessage,
	)
	if err != nil {
		logger.Err(err).Msg("error replying")
		return
	}
}

func (d *discordBot) loadMemberships(guildIDs ...string) {
	// de-duplicate desired channelSlugs
	channelSlugs := []string{}
	for _, guildID := range guildIDs {
		state := d.guildStates[guildID]
		for channelSlug := range state.GetMembershipInfo() {
			i := sort.SearchStrings(channelSlugs, channelSlug)
			if i >= len(channelSlugs) || channelSlugs[i] != channelSlug {
				channelSlugs = append(channelSlugs[:i], append([]string{channelSlug}, channelSlugs[i:]...)...)
			}
		}
	}
	d.ytChannelMembershipsMutex.Lock()
	defer d.ytChannelMembershipsMutex.Unlock()
	for _, channelSlug := range channelSlugs {
		// skip if this membership list was loaded in the last 2 minutes
		logger := log.With().Str("channelSlug", channelSlug).Logger()
		lastLoaded := d.ytChannelMembershipsLastLoaded[channelSlug]
		if lastLoaded.Add(time.Minute * 2).After(time.Now()) {
			logger.Debug().Str("lastLoaded", lastLoaded.Format(time.RFC3339)).
				Msg("skipping channel membership load - happened in the last 2 minutes")
			continue
		}
		// get all channel members
		memberIDs := map[string]struct{}{}
		iter := d.fs.
			Collection(common.ChannelCollection).Doc(channelSlug).
			Collection(common.ChannelMemberCollection).
			Select().Documents(d.ctx)
		for {
			snap, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				if c := status.Code(err); c != codes.OK {
					log.Err(err).Msg("rpc error getting channel memberships")
					break
				}
			}
			memberIDs[snap.Ref.ID] = struct{}{}
		}
		logger.Debug().Int("count", len(memberIDs)).Msg("loaded members for channel")
		d.ytChannelMemberships[channelSlug] = memberIDs
		d.ytChannelMembershipsLastLoaded[channelSlug] = time.Now()
	}
}

func (d *discordBot) enforceMembershipsAsync(guildID string) error {
	log.Info().Str("guildID", guildID).Msg("requesting guild members for enforcement")
	return d.dgSession.RequestGuildMembers(guildID, "", 0, false)
}

func (d *discordBot) cmdMiddleware(next dgc.ExecutionHandler) dgc.ExecutionHandler {
	return func(ctx *dgc.Ctx) {
		// guild must be enrolled
		guildState := d.guildStates[ctx.Event.GuildID]
		if guildState.LoadState == guildFirstEncounter {
			return
		}
		// TODO: channel-specific restrictions
		next(ctx)
	}
}

// checkRoles checks for the appropriate admin/etc roles.
// it assumes that guildState has a common.DiscordGuild and that we have the appropriate bot permissions
func (d *discordBot) checkRoles(session *discordgo.Session, guildID string, guild *discordgo.Guild) {
	var err error
	if guild == nil {
		guild, err = session.Guild(guildID)
		if err != nil {
			log.Err(err).Str("guildID", guildID).Msg("error loading guild while checking roles")
			return
		}
	}
	var logger = log.With().Str("guildID", guild.ID).Str("name", guild.Name).Logger()
	state := d.guildStates[guildID]
	// xref roles for the channel against admin and member roles
	if state.GetMembershipInfo() == nil {
		logger.Warn().Msg("guild has no registered members-only role ID")
	}
}

func (d *discordBot) reply(
	logger zerolog.Logger,
	guildID, channelID, userID string,
	messageRef *discordgo.MessageReference,
	message string,
) error {
	guildState := d.guildStates[guildID]
	if !guildState.noFancyReply {
		_, err := d.dgSession.ChannelMessageSendReply(channelID, message, messageRef)
		if err != nil {
			errString := err.Error()
			if strings.Contains(errString, "Cannot reply without permission to read message history") {
				logger.Debug().Msg("falling back to simple replies in this Discord guild")
				guildState.noFancyReply = true
				d.guildStates[guildID] = guildState
			} else if strings.Contains(errString, `{"message_reference": ["Unknown message"]}`) {
				logger.Debug().Err(err).Msg("fancy reply message reference probably deleted, falling back to simple reply")
			} else {
				logger.Err(err).Msg("error sending fancy reply")
				return err
			}
		} else {
			return nil
		}
	}
	// if noFancyReply || !readMessageHistoryPermission || "Unknown message"
	_, err := d.dgSession.ChannelMessageSend(channelID, fmt.Sprintf("<@%s> %s", userID, message))
	if err != nil {
		logger.Err(err).Msg("error sending simple reply")
	}
	return err
}

func (d *discordBot) startHeartbeat() *time.Ticker {
	// sloppy because this gets cleaned up on program exit
	ticker := time.NewTicker(time.Second * 30)
	go func() {
		for {
			<-ticker.C
			log.Debug().Msg("heartbeat")
		}
	}()
	return ticker
}

func makeLocalizer(bundle *i18n.Bundle, languageTag string) *i18n.Localizer {
	if languageTag == "" {
		i18n.NewLocalizer(bundle, "en-US")
	}
	return i18n.NewLocalizer(bundle, languageTag, "en-US")
}

// localizer.LocalizeMessage that panics
func mustLocalizeMessage(localizer *i18n.Localizer, message *i18n.Message) string {
	msg, err := localizer.LocalizeMessage(message)
	if err != nil {
		panic(err)
	}
	return msg
}

func userHasRole(member *discordgo.Member, roleID string) bool {
	for _, role := range member.Roles {
		if roleID == role {
			return true
		}
	}
	return false
}

const largeThreshold = 50

// StartOptions is an options struct to provide to Start()
type StartOptions struct {
	Token                        string
	APIClient                    api.ClientWithResponsesInterface
	FirestoreClient              *firestore.Client
	MembershipReloadSubscription *pubsub.Subscription
	Heartbeat                    bool
}

// Start does what you think it does.
func Start(
	ctx context.Context,
	options *StartOptions,
) error {
	dg, err := discordgo.New("Bot " + options.Token)
	if err != nil {
		return err
	}
	// set LargeThreshold to a guaranteed ceiling
	dg.Identify.LargeThreshold = largeThreshold
	// add the GUILD_MEMBERS intent so that we can get member stuff
	*dg.Identify.Intents |= discordgo.IntentsGuildMembers
	bot := discordBot{
		ctx:                            ctx,
		apiClient:                      options.APIClient,
		fs:                             options.FirestoreClient,
		dgSession:                      dg,
		lastMemberCheck:                map[string]time.Time{},
		ytChannels:                     map[string]common.Channel{},
		ytChannelMemberships:           map[string]map[string]struct{}{},
		ytChannelMembershipsLastLoaded: map[string]time.Time{},
		bundle:                         lang.NewBundle(),
	}
	// create a Firestore listener for guild associations and channel info
	err = bot.listenToGuildAssociations()
	if err != nil {
		log.Err(err).Msg("error initalizing bot guilds")
		return err
	}
	err = bot.listenToChannelChanges()
	if err != nil {
		log.Err(err).Msg("error initalizing bot channels")
		return err
	}
	// start the membership check notification listener
	bot.listenToMemberCheckUpdates(options.MembershipReloadSubscription)
	// construct router
	router := dgc.Create(&dgc.Router{
		Prefixes: []string{"!mg "},
		PingHandler: func(ctx *dgc.Ctx) {
			ctx.RespondText("Pong!")
		},
	})
	router.RegisterMiddleware(bot.cmdMiddleware)
	router.RegisterCmd(&dgc.Command{
		Name:        "check",
		Description: "Request a check for membership updates. Do this if you just became a member for a YouTube channel!",
		Usage:       "check",
		Handler:     bot.handleCmdCheck,
	})
	dg.AddHandler(bot.handleGuildCreate)
	dg.AddHandler(bot.handleGuildMembersChunk)
	dg.AddHandler(bot.handleGuildMemberUpdate)
	dg.AddHandler(router.Handler())
	router.RegisterDefaultHelpCommand(dg, nil)
	err = dg.Open()
	if err != nil {
		log.Err(err).Msg("error starting discordgo session")
		return err
	}
	bot.startHeartbeat()
	defer dg.Close()
	fmt.Println("Bot running - press CTRL-C ()to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-sc
	fmt.Println("interrupt/kill/term received")
	return nil
}
