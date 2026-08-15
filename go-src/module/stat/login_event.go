package stat

import (
	"github.com/gogf/gf/v2/os/gmlock"
	"time"
	"xr-game-server/core/event"
	"xr-game-server/dao/statdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/entity/stat"
	"xr-game-server/gameevent"
)

func initLoginEvent() {
	event.Sub(gameevent.LoginEvent, onLoginEvent)
}

func onLoginEvent(data any) {
	val := data.(uint64)
	userInfo := userinfodao.GetUserInfoByUserId(val)
	now := time.Now()
	userInfo.SetLastLoginTime(&now)

	lockName := "stat_login"
	gmlock.Lock(lockName)
	defer gmlock.Unlock(lockName)

	recordDailyLogin(now, val)
	recordWeeklyLogin(now, val)
	recordMonthlyLogin(now, val)
}

func recordDailyLogin(now time.Time, userId uint64) {
	date := entity.FormatDailyLoginStatDate(now)
	if statdao.TryRecordDailyLogin(date, userId) {
		stat := statdao.GetDailyLoginStatByDate(date)
		stat.AddLoginCount(1)
	}
}

func recordWeeklyLogin(now time.Time, userId uint64) {
	week := entity.FormatWeeklyLoginStatKey(now)
	if statdao.TryRecordWeeklyLogin(week, userId) {
		stat := statdao.GetWeeklyLoginStatByWeek(week)
		stat.AddLoginCount(1)
	}
}

func recordMonthlyLogin(now time.Time, userId uint64) {
	month := entity.FormatMonthlyLoginStatKey(now)
	if statdao.TryRecordMonthlyLogin(month, userId) {
		stat := statdao.GetMonthlyLoginStatByMonth(month)
		stat.AddLoginCount(1)
	}
}
