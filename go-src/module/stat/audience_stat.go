package stat

import (
	"time"

	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/dao/statdao"
	"xr-game-server/entity/stat"
)

// RecordValidAudience 记录有效观众(跨直播间去重,日/周/月各计一次)
func RecordValidAudience(userId uint64, statAt time.Time) {
	if userId == 0 {
		return
	}
	if statAt.IsZero() {
		statAt = time.Now()
	}

	lockName := "stat_audience"
	gmlock.Lock(lockName)
	defer gmlock.Unlock(lockName)

	recordPeriodAudienceUser(statAt, userId)
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
