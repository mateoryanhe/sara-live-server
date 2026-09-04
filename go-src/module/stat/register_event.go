package stat

import (
	"time"

	"xr-game-server/core/event"
	"xr-game-server/dao/statdao"
	"xr-game-server/entity/stat"
	"xr-game-server/gameevent"
)

func initRegisterEvent() {
	event.Sub(gameevent.RegisterEvent, onRegisterEvent)
}

func onRegisterEvent(data any) {
	enqueue(&statJob{Kind: jobRegister, Payload: data})
}

func consumeRegisterJob(job *statJob) {
	val, ok := job.Payload.(*gameevent.RegisterEventData)
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
	stat := statdao.GetDailyLoginStatByDate(date)
	stat.AddRegisterCount(1)
	statdao.PublishDailyLoginStat(stat)
}

func recordWeeklyRegister(now time.Time) {
	week := entity.FormatWeeklyLoginStatKey(now)
	stat := statdao.GetWeeklyLoginStatByWeek(week)
	stat.AddRegisterCount(1)
	statdao.PublishWeeklyLoginStat(stat)
}

func recordMonthlyRegister(now time.Time) {
	month := entity.FormatMonthlyLoginStatKey(now)
	stat := statdao.GetMonthlyLoginStatByMonth(month)
	stat.AddRegisterCount(1)
	statdao.PublishMonthlyLoginStat(stat)
}
