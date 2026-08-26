package stat

import (
	"time"
	"xr-game-server/dao/statdao"
	"xr-game-server/dao/userinfodao"
	statentity "xr-game-server/entity/stat"
	userentity "xr-game-server/entity/user"
)

func shouldCountUserStat(userId uint64) bool {
	if userId == 0 {
		return false
	}
	user := userinfodao.GetUserInfoByUserId(userId)
	return !userentity.UserTypeExcludedFromStat(user.UserType)
}

func recordPeriodGoldConsumeUser(statAt time.Time, userId uint64) {
	if !shouldCountUserStat(userId) {
		return
	}
	date := statentity.FormatDailyLoginStatDate(statAt)
	if statdao.TryRecordDailyGoldConsume(date, userId) {
		dailyStat := statdao.GetDailyLoginStatByDate(date)
		dailyStat.AddGoldConsumeUserCount(1)
		statdao.PublishDailyLoginStat(dailyStat)
	}
	week := statentity.FormatWeeklyLoginStatKey(statAt)
	if statdao.TryRecordWeeklyGoldConsume(week, userId) {
		weeklyStat := statdao.GetWeeklyLoginStatByWeek(week)
		weeklyStat.AddGoldConsumeUserCount(1)
		statdao.PublishWeeklyLoginStat(weeklyStat)
	}
	month := statentity.FormatMonthlyLoginStatKey(statAt)
	if statdao.TryRecordMonthlyGoldConsume(month, userId) {
		monthlyStat := statdao.GetMonthlyLoginStatByMonth(month)
		monthlyStat.AddGoldConsumeUserCount(1)
		statdao.PublishMonthlyLoginStat(monthlyStat)
	}
}

func recordPeriodDiamondConsumeUser(statAt time.Time, userId uint64) {
	if !shouldCountUserStat(userId) {
		return
	}
	date := statentity.FormatDailyLoginStatDate(statAt)
	if statdao.TryRecordDailyDiamondConsume(date, userId) {
		dailyStat := statdao.GetDailyLoginStatByDate(date)
		dailyStat.AddDiamondConsumeUserCount(1)
		statdao.PublishDailyLoginStat(dailyStat)
	}
	week := statentity.FormatWeeklyLoginStatKey(statAt)
	if statdao.TryRecordWeeklyDiamondConsume(week, userId) {
		weeklyStat := statdao.GetWeeklyLoginStatByWeek(week)
		weeklyStat.AddDiamondConsumeUserCount(1)
		statdao.PublishWeeklyLoginStat(weeklyStat)
	}
	month := statentity.FormatMonthlyLoginStatKey(statAt)
	if statdao.TryRecordMonthlyDiamondConsume(month, userId) {
		monthlyStat := statdao.GetMonthlyLoginStatByMonth(month)
		monthlyStat.AddDiamondConsumeUserCount(1)
		statdao.PublishMonthlyLoginStat(monthlyStat)
	}
}
