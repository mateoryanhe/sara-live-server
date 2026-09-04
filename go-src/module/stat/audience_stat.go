package stat

import (
	"time"

	"xr-game-server/core/event"
	"xr-game-server/dao/statdao"
	"xr-game-server/entity/stat"
	"xr-game-server/gameevent"
)

func initAudienceEvent() {
	event.Sub(gameevent.ValidAudienceEvent, onValidAudienceEvent)
}

func onValidAudienceEvent(data any) {
	enqueue(&statJob{Kind: jobAudience, Payload: data})
}

func consumeAudienceJob(job *statJob) {
	ev, ok := job.Payload.(*gameevent.ValidAudienceEventData)
	if !ok || ev == nil || ev.UserId == 0 {
		return
	}
	statAt := ev.StatAt
	if statAt.IsZero() {
		statAt = time.Now()
	}
	recordPeriodAudienceUser(statAt, ev.UserId)
}

func recordPeriodAudienceUser(statAt time.Time, userId uint64) {
	date := entity.FormatDailyLoginStatDate(statAt)
	if statdao.TryRecordDailyAudience(date, userId) {
		dailyStat := statdao.GetDailyLoginStatByDate(date)
		dailyStat.AddAudienceUserCount(1)
		statdao.PublishDailyLoginStat(dailyStat)
	}
	week := entity.FormatWeeklyLoginStatKey(statAt)
	if statdao.TryRecordWeeklyAudience(week, userId) {
		weeklyStat := statdao.GetWeeklyLoginStatByWeek(week)
		weeklyStat.AddAudienceUserCount(1)
		statdao.PublishWeeklyLoginStat(weeklyStat)
	}
	month := entity.FormatMonthlyLoginStatKey(statAt)
	if statdao.TryRecordMonthlyAudience(month, userId) {
		monthlyStat := statdao.GetMonthlyLoginStatByMonth(month)
		monthlyStat.AddAudienceUserCount(1)
		statdao.PublishMonthlyLoginStat(monthlyStat)
	}
}
