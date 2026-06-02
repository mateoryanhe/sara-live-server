package stat

import (
	"github.com/gogf/gf/v2/os/gmlock"
	"time"
	"xr-game-server/core/event"
	"xr-game-server/dao/dailyloginstatdao"
	"xr-game-server/dao/dailyuserlogindao"
	"xr-game-server/dao/monthlyloginstatdao"
	"xr-game-server/dao/monthlyuserlogindao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dao/weeklyloginstatdao"
	"xr-game-server/dao/weeklyuserlogindao"
	"xr-game-server/entity"
)

func initLoginEvent() {
	event.Sub(event.AppToken, onLoginEvent)
}

func onLoginEvent(data any) {
	val := data.(*event.AppTokenData)
	userInfo := userinfodao.GetUserInfoByUserId(val.Id)
	now := time.Now()
	userInfo.SetLastLoginTime(&now)

	lockName := "stat_login"
	gmlock.Lock(lockName)
	defer gmlock.Unlock(lockName)

	recordDailyLogin(now, val.Id)
	recordWeeklyLogin(now, val.Id)
	recordMonthlyLogin(now, val.Id)
}

func recordDailyLogin(now time.Time, userId uint64) {
	date := entity.FormatDailyLoginStatDate(now)
	if dailyuserlogindao.TryRecordLogin(date, userId) {
		stat := dailyloginstatdao.GetByDate(date)
		stat.AddLoginCount(1)
	}
}

func recordWeeklyLogin(now time.Time, userId uint64) {
	week := entity.FormatWeeklyLoginStatKey(now)
	if weeklyuserlogindao.TryRecordLogin(week, userId) {
		stat := weeklyloginstatdao.GetByWeek(week)
		stat.AddLoginCount(1)
	}
}

func recordMonthlyLogin(now time.Time, userId uint64) {
	month := entity.FormatMonthlyLoginStatKey(now)
	if monthlyuserlogindao.TryRecordLogin(month, userId) {
		stat := monthlyloginstatdao.GetByMonth(month)
		stat.AddLoginCount(1)
	}
}
