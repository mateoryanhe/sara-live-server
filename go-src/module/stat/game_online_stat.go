package stat

import (
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dao/statdao"
	statentity "xr-game-server/entity/stat"
	"xr-game-server/gameevent"
)

const (
	gameOnlineIdleTimeout        = time.Minute
	gameOnlineSweepInterval      = 10 * time.Second
	gameOnlinePeakSampleInterval = time.Minute
)

type gameOnlineTracker struct {
	mu        sync.Mutex
	lastBetAt map[uint64]time.Time
}

var gameOnline gameOnlineTracker

func initGameOnlineStat() {
	gameOnline.lastBetAt = make(map[uint64]time.Time)
	event.Sub(gameevent.GameBetCreatedEvent, onGameBetCreatedForOnlineStat)
	ctx := gctx.New()
	xrtimer.AddSingleton(ctx, gameOnlineSweepInterval, gameOnlineSweepTick)
	xrtimer.AddSingleton(ctx, gameOnlinePeakSampleInterval, gameOnlinePeakSampleTick)
}

func onGameBetCreatedForOnlineStat(data any) {
	ev, ok := data.(*gameevent.GameBetCreatedEventData)
	if !ok || ev == nil || ev.UserId == 0 || ev.Amount <= 0 {
		return
	}
	if !shouldCountUserStat(ev.UserId) {
		return
	}
	gameOnlineTouchBet(ev.UserId, time.Now())
}

func gameOnlineTouchBet(userId uint64, at time.Time) {
	gameOnline.mu.Lock()
	if gameOnline.lastBetAt == nil {
		gameOnline.lastBetAt = make(map[uint64]time.Time)
	}
	gameOnline.lastBetAt[userId] = at
	count := uint64(len(gameOnline.lastBetAt))
	gameOnline.mu.Unlock()
	recordGameOnlinePeakAt(at, count)
}

func gameOnlineSweepTick(_ context.Context) {
	now := time.Now()
	cutoff := now.Add(-gameOnlineIdleTimeout)
	gameOnline.mu.Lock()
	for userId, lastAt := range gameOnline.lastBetAt {
		if lastAt.Before(cutoff) {
			delete(gameOnline.lastBetAt, userId)
		}
	}
	count := uint64(len(gameOnline.lastBetAt))
	gameOnline.mu.Unlock()
	recordGameOnlinePeakAt(now, count)
	appendGameOnlineRecentSample(now, count)
}

func gameOnlinePeakSampleTick(_ context.Context) {
	now := time.Now()
	count := CurrentGameOnlineCount()
	recordGameOnlinePeakAt(now, count)
	appendGameOnlineRecentSample(now, count)
}

// CurrentGameOnlineCount 当前游戏在线人数(连续下注视为在线, 超过 1 分钟无下注视为离线)
func CurrentGameOnlineCount() uint64 {
	now := time.Now()
	cutoff := now.Add(-gameOnlineIdleTimeout)
	gameOnline.mu.Lock()
	defer gameOnline.mu.Unlock()
	for userId, lastAt := range gameOnline.lastBetAt {
		if lastAt.Before(cutoff) {
			delete(gameOnline.lastBetAt, userId)
		}
	}
	return uint64(len(gameOnline.lastBetAt))
}

func recordGameOnlinePeakAt(at time.Time, count uint64) {
	if count == 0 {
		return
	}
	date := statentity.FormatDailyLoginStatDate(at)
	daily := statdao.GetDailyGameOnlinePeakStatByDate(date)
	daily.RaisePeakCount(count)
	statdao.PublishDailyGameOnlinePeakStat(daily)

	week := statentity.FormatWeeklyLoginStatKey(at)
	weekly := statdao.GetWeeklyGameOnlinePeakStatByWeek(week)
	weekly.RaisePeakCount(count)
	statdao.PublishWeeklyGameOnlinePeakStat(weekly)

	month := statentity.FormatMonthlyGameBetUserStatKey(at)
	monthly := statdao.GetMonthlyGameOnlinePeakStatByMonth(month)
	monthly.RaisePeakCount(count)
	statdao.PublishMonthlyGameOnlinePeakStat(monthly)
}
