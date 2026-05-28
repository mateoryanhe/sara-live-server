package stat

import (
	"time"
	"xr-game-server/core/event"
	"xr-game-server/dao/dailyloginstatdao"
	"xr-game-server/dao/monthlyloginstatdao"
	"xr-game-server/dao/weeklyloginstatdao"
	"xr-game-server/entity"
	"xr-game-server/gameevent"
)

func initRegisterEvent() {
	event.Sub(gameevent.RegisterEvent, onRegisterEvent)
}

func onRegisterEvent(data any) {
	val, ok := data.(*gameevent.RegisterEventData)
	if !ok || val == nil || val.UserId == 0 {
		return
	}
	now := val.RegisteredAt
	if now.IsZero() {
		now = time.Now()
	}
	recordDailyRegister(now)
	recordWeeklyRegister(now)
	recordMonthlyRegister(now)
}

func recordDailyRegister(now time.Time) {
	date := entity.FormatDailyLoginStatDate(now)
	stat := dailyloginstatdao.GetByDate(date)
	stat.AddRegisterCount(1)
}

func recordWeeklyRegister(now time.Time) {
	week := entity.FormatWeeklyLoginStatKey(now)
	stat := weeklyloginstatdao.GetByWeek(week)
	stat.AddRegisterCount(1)
}

func recordMonthlyRegister(now time.Time) {
	month := entity.FormatMonthlyLoginStatKey(now)
	stat := monthlyloginstatdao.GetByMonth(month)
	stat.AddRegisterCount(1)
}
