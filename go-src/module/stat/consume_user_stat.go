package stat

import (
	"time"
	"xr-game-server/dao/statdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/entity"
)

func shouldCountUserStat(userId uint64) bool {
	if userId == 0 {
		return false
	}
	user := userinfodao.GetUserInfoByUserId(userId)
	return !entity.UserTypeExcludedFromStat(user.UserType)
}

func recordPeriodGoldConsumeUser(statAt time.Time, userId uint64) {
	if !shouldCountUserStat(userId) {
		return
	}
	date := entity.FormatDailyLoginStatDate(statAt)
	if statdao.TryRecordDailyGoldConsume(date, userId) {
		statdao.GetDailyLoginStatByDate(date).AddGoldConsumeUserCount(1)
	}
	week := entity.FormatWeeklyLoginStatKey(statAt)
	if statdao.TryRecordWeeklyGoldConsume(week, userId) {
		statdao.GetWeeklyLoginStatByWeek(week).AddGoldConsumeUserCount(1)
	}
	month := entity.FormatMonthlyLoginStatKey(statAt)
	if statdao.TryRecordMonthlyGoldConsume(month, userId) {
		statdao.GetMonthlyLoginStatByMonth(month).AddGoldConsumeUserCount(1)
	}
}

func recordPeriodDiamondConsumeUser(statAt time.Time, userId uint64) {
	if !shouldCountUserStat(userId) {
		return
	}
	date := entity.FormatDailyLoginStatDate(statAt)
	if statdao.TryRecordDailyDiamondConsume(date, userId) {
		statdao.GetDailyLoginStatByDate(date).AddDiamondConsumeUserCount(1)
	}
	week := entity.FormatWeeklyLoginStatKey(statAt)
	if statdao.TryRecordWeeklyDiamondConsume(week, userId) {
		statdao.GetWeeklyLoginStatByWeek(week).AddDiamondConsumeUserCount(1)
	}
	month := entity.FormatMonthlyLoginStatKey(statAt)
	if statdao.TryRecordMonthlyDiamondConsume(month, userId) {
		statdao.GetMonthlyLoginStatByMonth(month).AddDiamondConsumeUserCount(1)
	}
}
