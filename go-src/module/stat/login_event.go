package stat

import (
	"time"

	"xr-game-server/core/event"
	"xr-game-server/dao/statdao"
	"xr-game-server/entity/stat"
	"xr-game-server/gameevent"
	"xr-game-server/module/promotionattribution"
)

func initLoginEvent() {
	event.Sub(gameevent.LoginEvent, onLoginEvent)
}

func onLoginEvent(data any) {
	enqueue(&statJob{Kind: jobLogin, At: time.Now(), Payload: data})
}

func consumeLoginJob(job *statJob) {
	userId, ok := gameevent.ParseLoginEventUserId(job.Payload)
	if !ok {
		return
	}
	now := job.At
	if now.IsZero() {
		now = time.Now()
	}
	clientPlatform := resolveUserClientReportPlatform(userId)
	recordDailyLogin(now, userId, clientPlatform)
	recordWeeklyLogin(now, userId, clientPlatform)
	recordMonthlyLogin(now, userId, clientPlatform)
}

func recordDailyLogin(now time.Time, userId uint64, clientPlatform string) {
	date := entity.FormatDailyLoginStatDate(now)
	if !statdao.TryRecordDailyLogin(date, userId) {
		return
	}
	stat := statdao.GetDailyLoginStatByDate(date)
	stat.AddLoginCount(1)
	switch clientPlatform {
	case promotionattribution.RegisterPlatformApp:
		stat.AddAppDeviceLoginCount(1)
	case promotionattribution.RegisterPlatformH5:
		stat.AddH5DeviceLoginCount(1)
	}
	statdao.PublishDailyLoginStat(stat)
}

func recordWeeklyLogin(now time.Time, userId uint64, clientPlatform string) {
	week := entity.FormatWeeklyLoginStatKey(now)
	if !statdao.TryRecordWeeklyLogin(week, userId) {
		return
	}
	stat := statdao.GetWeeklyLoginStatByWeek(week)
	stat.AddLoginCount(1)
	switch clientPlatform {
	case promotionattribution.RegisterPlatformApp:
		stat.AddAppDeviceLoginCount(1)
	case promotionattribution.RegisterPlatformH5:
		stat.AddH5DeviceLoginCount(1)
	}
	statdao.PublishWeeklyLoginStat(stat)
}

func recordMonthlyLogin(now time.Time, userId uint64, clientPlatform string) {
	month := entity.FormatMonthlyLoginStatKey(now)
	if !statdao.TryRecordMonthlyLogin(month, userId) {
		return
	}
	stat := statdao.GetMonthlyLoginStatByMonth(month)
	stat.AddLoginCount(1)
	switch clientPlatform {
	case promotionattribution.RegisterPlatformApp:
		stat.AddAppDeviceLoginCount(1)
	case promotionattribution.RegisterPlatformH5:
		stat.AddH5DeviceLoginCount(1)
	}
	statdao.PublishMonthlyLoginStat(stat)
}
