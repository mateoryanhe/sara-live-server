package stat

import (
	"time"

	"xr-game-server/core/event"
	"xr-game-server/dao/statdao"
	statentity "xr-game-server/entity/stat"
	"xr-game-server/gameevent"
)

func initGameBetUserStatEvent() {
	event.Sub(gameevent.GameBetCreatedEvent, onGameBetCreatedForUserStat)
}

func onGameBetCreatedForUserStat(data any) {
	ev, ok := data.(*gameevent.GameBetCreatedEventData)
	if !ok || ev == nil || ev.UserId == 0 || ev.Amount <= 0 {
		return
	}
	recordPeriodGameBetUser(time.Now(), ev.UserId)
}

func recordPeriodGameBetUser(statAt time.Time, userId uint64) {
	if !shouldCountUserStat(userId) {
		return
	}
	date := statentity.FormatDailyLoginStatDate(statAt)
	if statdao.TryRecordDailyGameBet(date, userId) {
		daily := statdao.GetDailyGameBetUserStatByDate(date)
		daily.AddGameBetUserCount(1)
		statdao.PublishDailyGameBetUserStat(daily)
	}
	week := statentity.FormatWeeklyLoginStatKey(statAt)
	if statdao.TryRecordWeeklyGameBet(week, userId) {
		weekly := statdao.GetWeeklyGameBetUserStatByWeek(week)
		weekly.AddGameBetUserCount(1)
		statdao.PublishWeeklyGameBetUserStat(weekly)
	}
	month := statentity.FormatMonthlyGameBetUserStatKey(statAt)
	if statdao.TryRecordMonthlyGameBet(month, userId) {
		monthly := statdao.GetMonthlyGameBetUserStatByMonth(month)
		monthly.AddGameBetUserCount(1)
		statdao.PublishMonthlyGameBetUserStat(monthly)
	}
}
