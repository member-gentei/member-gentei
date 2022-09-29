package roles

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type RoleUpdateTracker struct {
	session *discordgo.Session

	trackHookMutex     *sync.Mutex
	trackHooks         map[string]RoleUpdateTrackFunc
	trackHookPollMutex *sync.Mutex
	trackHookPollCtx   map[string]context.CancelFunc
}

type RoleUpdateTrackData struct {
	// Fields extracted from *discordgo.GuildMemberUpdate
	GuildID string   // gmu.GuildID
	UserID  string   // gmu.User.ID
	Roles   []string // gmu.Roles

	// If this information comes from polling instead of an actual GUILD_MEMBER_UPDATE event
	Polled bool
}

type RoleUpdateTrackFunc func(RoleUpdateTrackData) (removeHook bool)

// TrackHook executes a function on GUILD_MEMBER_UPDATE.
func (r *RoleUpdateTracker) TrackHook(guildID, userID string, trackFunc RoleUpdateTrackFunc) {
	r.trackHookMutex.Lock()
	r.trackHookPollMutex.Lock()
	defer r.trackHookMutex.Unlock()
	defer r.trackHookPollMutex.Unlock()
	var (
		pollCtx, cancel = context.WithTimeout(context.Background(), time.Second*60)
		key             = r.trackHookKey(guildID, userID)
	)
	// poll for changes every ~10 seconds for 60 seconds for good measure
	r.trackHookPollCtx[key] = cancel
	r.trackHooks[key] = trackFunc
	go func() {
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()
		cleanUp := func() {
			r.trackHookPollMutex.Lock()
			defer r.trackHookPollMutex.Unlock()
			delete(r.trackHookPollCtx, key)
		}
		for {
			select {
			case <-pollCtx.Done():
				cleanUp()
				return
			case <-ticker.C:
				member, err := r.session.GuildMember(guildID, userID)
				if err != nil {
					log.Err(err).
						Str("guildID", guildID).
						Str("userID", userID).
						Msg("error polling for GuildMember, stopping")
					cleanUp()
					return
				}
				r.trackHookMutex.Lock()
				trackHook := r.trackHooks[key]
				r.trackHookMutex.Unlock()
				if trackHook == nil {
					// if the hook disappeared, we're done
					cleanUp()
					return
				} else if trackHook(RoleUpdateTrackData{
					GuildID: member.GuildID,
					UserID:  userID,
					Roles:   member.Roles,
					Polled:  true,
				}) {
					r.trackHookMutex.Lock()
					defer r.trackHookMutex.Unlock()
					delete(r.trackHooks, key)
					cleanUp()
					return
				}
			}
		}
	}()
}

func (r *RoleUpdateTracker) trackHookKey(guildID, userID string) string {
	return fmt.Sprintf("%s-%s", guildID, userID)
}

func (r *RoleUpdateTracker) start() {
	r.session.AddHandler(func(s *discordgo.Session, gmu *discordgo.GuildMemberUpdate) {
		key := r.trackHookKey(gmu.GuildID, gmu.User.ID)
		r.trackHookMutex.Lock()
		defer r.trackHookMutex.Unlock()
		f, exists := r.trackHooks[key]
		if !exists {
			return
		}
		if f(RoleUpdateTrackData{
			GuildID: gmu.GuildID,
			UserID:  gmu.User.ID,
			Roles:   gmu.Roles,
		}) {
			delete(r.trackHooks, key)
		}
	})
}

func NewRoleUpdateTracker(session *discordgo.Session) *RoleUpdateTracker {
	tracker := &RoleUpdateTracker{
		session:            session,
		trackHookMutex:     &sync.Mutex{},
		trackHooks:         map[string]RoleUpdateTrackFunc{},
		trackHookPollMutex: &sync.Mutex{},
		trackHookPollCtx:   map[string]context.CancelFunc{},
	}
	tracker.start()
	return tracker
}
