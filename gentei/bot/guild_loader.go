package bot

import (
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mark-ignacio/gsync"
	"github.com/rs/zerolog/log"
)

// largeGuildLoader tracks the loading process for servers that have huge member lists.
type largeGuildLoader struct {
	session                  *discordgo.Session
	guildMemberLoadMutexes   *gsync.Map[string, *sync.Mutex]
	guildMemberRequests      map[string]*lglStats
	guildMemberRequestsMutex sync.Mutex
	inactiveInterval         time.Duration

	// consume the channel before taking a lock
	// reset the ticker after taking a lock
	inactiveTicker *time.Ticker
}

type lglStats struct {
	FirstRequest      time.Time
	LastChunkReceived time.Time
	Retries           int
}

func (l *largeGuildLoader) StartWatchdog(retryBatchSize int, maxRetries int) {
	l.inactiveTicker = time.NewTicker(l.inactiveInterval)
	var alreadyStartedMode bool
	for {
		<-l.inactiveTicker.C
		if !l.session.DataReady {
			continue
		}
		left := l.doWatchdog(retryBatchSize, maxRetries)
		if left == 0 && !alreadyStartedMode {
			alreadyStartedMode = true
			idleInterval := l.inactiveInterval * 10
			l.inactiveTicker.Reset(idleInterval)
			log.Info().
				Dur("interval", idleInterval).
				Msg("all guild members loaded, changing watchdog to longer check interval")
		}
	}
}

// Returns how many guilds need to be loaded.
func (l *largeGuildLoader) doWatchdog(retryBatchSize int, maxRetries int) int {
	l.guildMemberRequestsMutex.Lock()
	defer l.guildMemberRequestsMutex.Unlock()
	var (
		retriedGuildCount int
		deleteGuildIDs    []string
		chunkThreshold    = time.Now().Add(-time.Second * 10)
	)
	// go through all uncompleted requests
	for guildID, stats := range l.guildMemberRequests {
		logger := log.With().Str("guildID", guildID).Logger()
		if stats.Retries > maxRetries {
			logger.Warn().Msg("reached max retires for loading guild member list")
			deleteGuildIDs = append(deleteGuildIDs, guildID)
			continue
		} else if retriedGuildCount < retryBatchSize {
			m, _ := l.guildMemberLoadMutexes.LoadOrStore(guildID, &sync.Mutex{})
			m.TryLock() // allows working through the occasional duplicate on startup
			if err := l.session.RequestGuildMembers(guildID, "", 0, "rgc-"+guildID, false); err != nil {
				logger.Err(err).Msg("error re-requesting guild members")
			}
			logger.Debug().Msg("re-requested guild members")
			stats.Retries++
			retriedGuildCount++
		} else if stats.LastChunkReceived.After(chunkThreshold) {
			logger.Debug().Time("lastChunk", stats.LastChunkReceived).Msg("last chunk rather recent, no retry required")
		} else {
			logger.Debug().Msg("skipping retry for now, too many to load")
		}
	}
	for _, guildID := range deleteGuildIDs {
		delete(l.guildMemberRequests, guildID)
	}
	lgms := len(l.guildMemberRequests)
	if lgms > 0 {
		log.Info().Int("count", lgms).Msg("guilds left to retry loading member lists")
	}
	return lgms
}

func (l *largeGuildLoader) GuildCreateHandler(s *discordgo.Session, gc *discordgo.GuildCreate) {
	logger := log.With().
		Str("guildID", gc.ID).
		Str("guildName", gc.Name).
		Logger()
	logger.Info().Msg("joined Guild")
	m, _ := l.guildMemberLoadMutexes.LoadOrStore(gc.ID, &sync.Mutex{})
	// start guild member load if >= largeThreshold
	// (see why at https://discord.com/developers/docs/topics/gateway-events#request-guild-members)
	if gc.MemberCount < s.Identify.LargeThreshold {
		return
	}
	if !m.TryLock() {
		logger.Info().Msg("something else locked the guildMemberLoadMutex?")
		m.Lock()
	}
	logger.Info().Int("memberCount", gc.MemberCount).Msg("big server; requesting Guild members")
	l.guildMemberRequestsMutex.Lock()
	defer l.guildMemberRequestsMutex.Unlock()
	l.guildMemberRequests[gc.ID] = &lglStats{
		FirstRequest: time.Now(),
	}
	if err := l.session.RequestGuildMembers(gc.ID, "", 0, l.getNonce(gc.ID), false); err != nil {
		logger.Err(err).Msg("error requesting guild members")
	}
}

func (l *largeGuildLoader) GuildMembersChunkHandler(s *discordgo.Session, gmc *discordgo.GuildMembersChunk) {
	l.guildMemberRequestsMutex.Lock()
	defer l.guildMemberRequestsMutex.Unlock()
	l.inactiveTicker.Reset(l.inactiveInterval)
	logger := log.With().
		Str("guildID", gmc.GuildID).
		Int("total", gmc.ChunkCount).
		Logger()
	logger.Trace().
		Int("chunkIndex", gmc.ChunkIndex).
		Int("chunkCount", gmc.ChunkCount).
		Msg("got guild member chunk")
	if gmc.Nonce != l.getNonce(gmc.GuildID) {
		return
	}
	if gmc.ChunkIndex == 0 {
		logger.Info().Msg("got first guild member chunk")
		gmr, found := l.guildMemberRequests[gmc.GuildID]
		if !found {
			l.guildMemberRequests[gmc.GuildID] = &lglStats{LastChunkReceived: time.Now()}
		} else {
			gmr.LastChunkReceived = time.Now()
		}
	}
	if gmc.ChunkIndex == gmc.ChunkCount-1 {
		logger.Info().Msg("got all guild member chunks")
		delete(l.guildMemberRequests, gmc.GuildID)
		m, ok := l.guildMemberLoadMutexes.Load(gmc.GuildID)
		if ok {
			m.Unlock()
		}
	}
}

func (l *largeGuildLoader) getNonce(guildID string) string {
	return "rgc-" + guildID
}
