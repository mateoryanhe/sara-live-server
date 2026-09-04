package stat

import (
	"time"

	"xr-game-server/core/event"
	"xr-game-server/dao/statdao"
	"xr-game-server/entity/stat"
	"xr-game-server/gameevent"
)

func initLoginEvent() {
	event.Sub(gameevent.LoginEvent, onLoginEvent)
}

func onLoginEvent(data any) {
	enqueue(&statJob{Kind: jobLogin, At: time.Now(), Payload: data})
}

func consumeLoginJob(job *statJob) {
	userId, ok := job.Payload.(uint64)
	if !ok || userId == 0 {
		return
	}
	now := job.At
	if now.IsZero() {
		now = time.Now()
	}
	recordDailyLogin(now, userId)
	recordWeeklyLogin(now, userId)
	recordMonthlyLogin(now, userId)
}

func recordDailyLogin(now time.Time, userId uint64) {
	date := entity.FormatDailyLoginStatDate(now)
	if statdao.TryRecordDailyLogin(date, userId) {
		stat := statdao.GetDailyLoginStatByDate(date)
		stat.AddLoginCount(1)
		statdao.PublishDailyLoginStat(stat)
	}
}

func recordWeeklyLogin(now time.Time, userId uint64) {
	week := entity.FormatWeeklyLoginStatKey(now)
	if statdao.TryRecordWeeklyLogin(week, userId) {
		stat := statdao.GetWeeklyLoginStatByWeek(week)
		stat.AddLoginCount(1)
		statdao.PublishWeeklyLoginStat(stat)
	}
}

func recordMonthlyLogin(now time.Time, userId uint64) {
	month := entity.FormatMonthlyLoginStatKey(now)
	if statdao.TryRecordMonthlyLogin(month, userId) {
		stat := statdao.GetMonthlyLoginStatByMonth(month)
		stat.AddLoginCount(1)
		statdao.PublishMonthlyLoginStat(stat)
	}
}
