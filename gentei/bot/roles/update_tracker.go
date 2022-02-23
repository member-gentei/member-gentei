package roles

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type RoleUpdateTracker struct {
	session    *discordgo.Session
	mutex      *sync.Mutex
	trackHooks map[string]RoleUpdateTrackFunc
}

type RoleUpdateTrackFunc func(*discordgo.GuildMemberUpdate) (removeHook bool)

// TrackHook executes a function on GUILD_MEMBER_UPDATE.
func (r *RoleUpdateTracker) TrackHook(guildID, userID string, trackFunc RoleUpdateTrackFunc) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.trackHooks[r.trackHookKey(guildID, userID)] = trackFunc
}

func (r *RoleUpdateTracker) trackHookKey(guildID, userID string) string {
	return fmt.Sprintf("%s-%s", guildID, userID)
}

func (r *RoleUpdateTracker) start() {
	r.session.AddHandler(func(s *discordgo.Session, gmu *discordgo.GuildMemberUpdate) {
		key := r.trackHookKey(gmu.GuildID, gmu.User.ID)
		r.mutex.Lock()
		defer r.mutex.Unlock()
		f, exists := r.trackHooks[key]
		if !exists {
			return
		}
		if f(gmu) {
			delete(r.trackHooks, key)
		}
	})
}

func NewRoleUpdateTracker(session *discordgo.Session) *RoleUpdateTracker {
	tracker := &RoleUpdateTracker{
		session:    session,
		mutex:      &sync.Mutex{},
		trackHooks: map[string]RoleUpdateTrackFunc{},
	}
	tracker.start()
	return tracker
}
