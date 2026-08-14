package timezonecfg

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/core/xrcorn"
	"xr-game-server/core/xrlog"
	"xr-game-server/entity"
	"xr-game-server/gameevent"
)

var (
	mu            sync.RWMutex
	timezoneCrons = make(map[int]*timezoneCron) // key: timezone offset
)

type timezoneCron struct {
	timezone   int
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// Init 从 live_guilds / live_rooms 收集已使用时区并启动对应定时任务
func Init() {
	for _, tz := range loadUsedTimezones() {
		StartCron(tz)
	}
}

func loadUsedTimezones() []int {
	seen := make(map[int]struct{})
	var result []int

	add := func(tz int) {
		if _, ok := seen[tz]; ok {
			return
		}
		seen[tz] = struct{}{}
		result = append(result, tz)
	}

	var guildTzs []int
	if err := g.DB().Model(string(entity.TbLiveGuild)).
		Fields("DISTINCT timezone").
		Where("status", entity.LiveGuildStatusNormal).
		Scan(&guildTzs); err != nil {
		xrlog.ErrorWithErr(gctx.New(), "TimezoneCfg", fmt.Sprintf("load guild timezones failed: %v", err), err)
	} else {
		for _, tz := range guildTzs {
			add(tz)
		}
	}

	var roomTzs []int
	if err := g.DB().Model(string(entity.TbLiveRoom)).
		Fields("DISTINCT timezone").
		Scan(&roomTzs); err != nil {
		xrlog.ErrorWithErr(gctx.New(), "TimezoneCfg", fmt.Sprintf("load room timezones failed: %v", err), err)
	} else {
		for _, tz := range roomTzs {
			add(tz)
		}
	}

	if len(result) == 0 {
		add(0)
	}
	return result
}

// EnsureCron 确保指定时区的定时任务已启动(设置时区时调用)
func EnsureCron(timezone int) {
	StartCron(timezone)
}

// StartCron 启动指定时区的定时任务
func StartCron(timezone int) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := timezoneCrons[timezone]; exists {
		xrlog.DetailLog.Info(gctx.New(), "timezonecfg: timezone", timezone, "already running")
		return
	}

	ctx, cancel := context.WithCancel(gctx.New())
	tc := &timezoneCron{
		timezone:   timezone,
		ctx:        ctx,
		cancelFunc: cancel,
	}

	dayPattern := buildCronPattern(timezone, "day")
	_, err := xrcorn.AddSingleton(ctx, dayPattern, func(ctx context.Context) {
		xrlog.DetailLog.Info(ctx, "timezonecfg: executing day task for timezone", timezone)
		event.Pub(gameevent.TimezoneDayEvent, gameevent.NewTimezoneDayEventData(timezone))
		xrlog.DetailLog.Info(ctx, "timezonecfg: day task done for timezone", timezone)
	})
	if err != nil {
		xrlog.ErrorWithErr(gctx.New(), "TimezoneCfg", fmt.Sprintf("failed to add day cron for timezone %d: %v", timezone, err), err)
		cancel()
		return
	}

	weekPattern := buildCronPattern(timezone, "week")
	_, err = xrcorn.AddSingleton(ctx, weekPattern, func(ctx context.Context) {
		xrlog.DetailLog.Info(ctx, "timezonecfg: executing week task for timezone", timezone)
		event.Pub(gameevent.TimezoneWeekEvent, gameevent.NewTimezoneWeekEventData(timezone))
		xrlog.DetailLog.Info(ctx, "timezonecfg: week task done for timezone", timezone)
	})
	if err != nil {
		xrlog.ErrorWithErr(gctx.New(), "TimezoneCfg", fmt.Sprintf("failed to add week cron for timezone %d: %v", timezone, err), err)
		cancel()
		return
	}

	timezoneCrons[timezone] = tc

	dayNext, _ := xrcorn.GetNextTimeWithReference(dayPattern, time.Now().In(time.UTC))
	weekNext, _ := xrcorn.GetNextTimeWithReference(weekPattern, time.Now().In(time.UTC))
	xrlog.DetailLog.Infof(gctx.New(), "timezonecfg: started for timezone %d, day next: %v, week next: %v", timezone, dayNext, weekNext)
}

// StopCron 停止指定时区的定时任务
func StopCron(timezone int) {
	mu.Lock()
	defer mu.Unlock()

	tc, exists := timezoneCrons[timezone]
	if !exists {
		return
	}

	tc.cancelFunc()
	delete(timezoneCrons, timezone)
	xrlog.DetailLog.Infof(gctx.New(), "timezonecfg: stopped for timezone %d", timezone)
}

func buildCronPattern(timezone int, taskType string) string {
	hour := (24 + timezone) % 24

	if taskType == "week" {
		return fmt.Sprintf("0 0 %d * * 1", hour)
	}

	return fmt.Sprintf("0 0 %d * * *", hour)
}
