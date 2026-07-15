package stat

import (
	"time"

	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/dao/statdao"
	"xr-game-server/entity"
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
		statdao.GetDailyLoginStatByDate(date).AddAudienceUserCount(1)
	}
	week := entity.FormatWeeklyLoginStatKey(statAt)
	if statdao.TryRecordWeeklyAudience(week, userId) {
		statdao.GetWeeklyLoginStatByWeek(week).AddAudienceUserCount(1)
	}
	month := entity.FormatMonthlyLoginStatKey(statAt)
	if statdao.TryRecordMonthlyAudience(month, userId) {
		statdao.GetMonthlyLoginStatByMonth(month).AddAudienceUserCount(1)
	}
}
