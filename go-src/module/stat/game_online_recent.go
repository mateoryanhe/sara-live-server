package stat

import (
	"sync"
	"time"

	"xr-game-server/core/event"
	"xr-game-server/dto/statdto"
	"xr-game-server/gameevent"
)

const (
	gameRecentTrendRetention = 2 * time.Hour
	gameRecentTrendMinute    = time.Minute
	gameRecentTrendPoints    = 60
)

type gameOnlineRecentSample struct {
	At    time.Time
	Count uint64
}

type gameRecentTrend struct {
	mu       sync.Mutex
	samples  []gameOnlineRecentSample
	betUsers map[string]map[uint64]struct{}
}

var gameRecent gameRecentTrend

func initGameOnlineRecentTrend() {
	gameRecent.betUsers = make(map[string]map[uint64]struct{})
	event.Sub(gameevent.GameBetCreatedEvent, onGameBetCreatedForRecentTrend)
}

func onGameBetCreatedForRecentTrend(data any) {
	ev, ok := data.(*gameevent.GameBetCreatedEventData)
	if !ok || ev == nil || ev.UserId == 0 || ev.Amount <= 0 {
		return
	}
	if !shouldCountUserStat(ev.UserId) {
		return
	}
	recordGameRecentBetUser(time.Now(), ev.UserId)
}

func recordGameRecentBetUser(at time.Time, userId uint64) {
	key := formatGameRecentMinuteKey(at)
	gameRecent.mu.Lock()
	defer gameRecent.mu.Unlock()
	if gameRecent.betUsers == nil {
		gameRecent.betUsers = make(map[string]map[uint64]struct{})
	}
	bucket, ok := gameRecent.betUsers[key]
	if !ok {
		bucket = make(map[uint64]struct{})
		gameRecent.betUsers[key] = bucket
	}
	bucket[userId] = struct{}{}
	pruneGameRecentLocked(at.Add(-gameRecentTrendRetention))
}

func appendGameOnlineRecentSample(at time.Time, count uint64) {
	gameRecent.mu.Lock()
	defer gameRecent.mu.Unlock()
	gameRecent.samples = append(gameRecent.samples, gameOnlineRecentSample{At: at, Count: count})
	pruneGameRecentLocked(at.Add(-gameRecentTrendRetention))
}

func pruneGameRecentLocked(cutoff time.Time) {
	for len(gameRecent.samples) > 0 && gameRecent.samples[0].At.Before(cutoff) {
		gameRecent.samples = gameRecent.samples[1:]
	}
	for key := range gameRecent.betUsers {
		t, ok := parseGameRecentMinuteKey(key)
		if !ok || t.Before(cutoff) {
			delete(gameRecent.betUsers, key)
		}
	}
}

func formatGameRecentMinuteKey(at time.Time) string {
	at = at.Truncate(gameRecentTrendMinute)
	return at.Format("2006-01-02 15:04")
}

func parseGameRecentMinuteKey(key string) (time.Time, bool) {
	t, err := time.ParseInLocation("2006-01-02 15:04", key, time.Local)
	return t, err == nil
}

func gameOnlineCountAtMinute(samples []gameOnlineRecentSample, minute time.Time) uint64 {
	minute = minute.Truncate(gameRecentTrendMinute)
	var best gameOnlineRecentSample
	found := false
	for _, s := range samples {
		sm := s.At.Truncate(gameRecentTrendMinute)
		if sm.After(minute) {
			break
		}
		if sm.Equal(minute) || sm.Before(minute) {
			best = s
			found = true
			if sm.Equal(minute) {
				break
			}
		}
	}
	if !found {
		return 0
	}
	return best.Count
}

// SnapshotGameTrendCountsForMetric 与资源细采样(10s)对齐：瞬时游戏在线 + 当前分钟累计下注人数
func SnapshotGameTrendCountsForMetric(at time.Time) (gameOnline, minuteBetUsers uint64) {
	gameOnline = CurrentGameOnlineCount()
	key := formatGameRecentMinuteKey(at)
	gameRecent.mu.Lock()
	if bucket := gameRecent.betUsers[key]; len(bucket) > 0 {
		minuteBetUsers = uint64(len(bucket))
	}
	gameRecent.mu.Unlock()
	return gameOnline, minuteBetUsers
}

func buildLastHourGameTrendPoints(now time.Time) []*statdto.CMSGameBetUserTrendPoint {
	now = now.Truncate(gameRecentTrendMinute)
	gameRecent.mu.Lock()
	samples := append([]gameOnlineRecentSample(nil), gameRecent.samples...)
	betCopy := make(map[string]map[uint64]struct{}, len(gameRecent.betUsers))
	for k, v := range gameRecent.betUsers {
		inner := make(map[uint64]struct{}, len(v))
		for uid := range v {
			inner[uid] = struct{}{}
		}
		betCopy[k] = inner
	}
	gameRecent.mu.Unlock()

	liveOnline := CurrentGameOnlineCount()
	out := make([]*statdto.CMSGameBetUserTrendPoint, 0, gameRecentTrendPoints)
	for i := gameRecentTrendPoints - 1; i >= 0; i-- {
		at := now.Add(-time.Duration(i) * gameRecentTrendMinute)
		key := formatGameRecentMinuteKey(at)
		var betCount uint64
		if bucket := betCopy[key]; len(bucket) > 0 {
			betCount = uint64(len(bucket))
		}
		online := gameOnlineCountAtMinute(samples, at)
		if i == 0 {
			online = maxUint64(online, liveOnline)
		}
		out = append(out, &statdto.CMSGameBetUserTrendPoint{
			Time:        key,
			Count:       betCount,
			OnlineCount: online,
		})
	}
	return out
}
