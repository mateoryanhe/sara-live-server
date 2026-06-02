package stat

import (
	"time"
	"xr-game-server/dao/dailyloginstatdao"
	"xr-game-server/dao/dailyuserdiamondconsumdao"
	"xr-game-server/dao/dailyusergoldconsumdao"
	"xr-game-server/dao/monthlyloginstatdao"
	"xr-game-server/dao/monthlyuserdiamondconsumdao"
	"xr-game-server/dao/monthlyusergoldconsumdao"
	"xr-game-server/dao/weeklyloginstatdao"
	"xr-game-server/dao/weeklyuserdiamondconsumdao"
	"xr-game-server/dao/weeklyusergoldconsumdao"
	"xr-game-server/entity"
)

func recordPeriodGoldConsumeUser(statAt time.Time, userId uint64) {
	if userId == 0 {
		return
	}
	date := entity.FormatDailyLoginStatDate(statAt)
	if dailyusergoldconsumdao.TryRecordConsume(date, userId) {
		dailyloginstatdao.GetByDate(date).AddGoldConsumeUserCount(1)
	}
	week := entity.FormatWeeklyLoginStatKey(statAt)
	if weeklyusergoldconsumdao.TryRecordConsume(week, userId) {
		weeklyloginstatdao.GetByWeek(week).AddGoldConsumeUserCount(1)
	}
	month := entity.FormatMonthlyLoginStatKey(statAt)
	if monthlyusergoldconsumdao.TryRecordConsume(month, userId) {
		monthlyloginstatdao.GetByMonth(month).AddGoldConsumeUserCount(1)
	}
}

func recordPeriodDiamondConsumeUser(statAt time.Time, userId uint64) {
	if userId == 0 {
		return
	}
	date := entity.FormatDailyLoginStatDate(statAt)
	if dailyuserdiamondconsumdao.TryRecordConsume(date, userId) {
		dailyloginstatdao.GetByDate(date).AddDiamondConsumeUserCount(1)
	}
	week := entity.FormatWeeklyLoginStatKey(statAt)
	if weeklyuserdiamondconsumdao.TryRecordConsume(week, userId) {
		weeklyloginstatdao.GetByWeek(week).AddDiamondConsumeUserCount(1)
	}
	month := entity.FormatMonthlyLoginStatKey(statAt)
	if monthlyuserdiamondconsumdao.TryRecordConsume(month, userId) {
		monthlyloginstatdao.GetByMonth(month).AddDiamondConsumeUserCount(1)
	}
}
