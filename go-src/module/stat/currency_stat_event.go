package stat

import (
	"time"

	"xr-game-server/constants/currency"
	"xr-game-server/core/event"
	"xr-game-server/dao/statdao"
	"xr-game-server/entity/stat"
	"xr-game-server/gameevent"
)

func initCurrencyStatEvent() {
	event.Sub(gameevent.CurrencyChangeEvent, onCurrencyChangeForStat)
}

func onCurrencyChangeForStat(data any) {
	enqueue(&statJob{Kind: jobCurrency, At: time.Now(), Payload: data})
}

func consumeCurrencyJob(job *statJob) {
	ev, ok := job.Payload.(*gameevent.CurrencyChangeEventData)
	if !ok || ev == nil {
		return
	}
	statAt := job.At
	if statAt.IsZero() {
		statAt = time.Now()
	}
	consumeGoldChange(ev)
	consumeGoldConsume(ev, statAt)
	consumeDiamondConsume(ev, statAt)
}

func consumeGoldChange(ev *gameevent.CurrencyChangeEventData) {
	if ev.Type != gameevent.CurrencyTypeGold || ev.Amount <= 0 {
		return
	}
	if !shouldCountUserStat(ev.UserId) {
		return
	}
	var delta float64
	switch ev.Action {
	case gameevent.CurrencyActionAdd:
		delta = ev.Amount
	case gameevent.CurrencyActionSub:
		delta = -ev.Amount
	default:
		return
	}
	stat := statdao.GetSysStat()
	if stat == nil {
		return
	}
	stat.AddTotalGold(delta)
}

func consumeGoldConsume(ev *gameevent.CurrencyChangeEventData, statAt time.Time) {
	if ev.Type != gameevent.CurrencyTypeGold || ev.Action != gameevent.CurrencyActionSub || ev.Amount <= 0 {
		return
	}
	if !shouldCountUserStat(ev.UserId) {
		return
	}
	if stat := statdao.GetSysStat(); stat != nil {
		stat.AddTotalGoldConsume(ev.Amount)
	}
	recordPeriodGoldConsume(statAt, ev.Amount)
	recordPeriodGoldConsumeUser(statAt, ev.UserId)
}

func consumeDiamondConsume(ev *gameevent.CurrencyChangeEventData, statAt time.Time) {
	if ev.Type != gameevent.CurrencyTypeDiamond || ev.Amount <= 0 {
		return
	}
	if !shouldCountUserStat(ev.UserId) {
		return
	}
	var delta float64
	switch ev.Action {
	case gameevent.CurrencyActionSub:
		delta = ev.Amount
	case gameevent.CurrencyActionAdd:
		if ev.Reason != currency.ReasonRefund {
			return
		}
		delta = -ev.Amount
	default:
		return
	}
	if stat := statdao.GetSysStat(); stat != nil {
		stat.AddTotalDiamondConsume(delta)
	}
	recordPeriodDiamondConsume(statAt, delta)
	if delta > 0 {
		recordPeriodDiamondConsumeUser(statAt, ev.UserId)
	}
}

func recordPeriodGoldConsume(statAt time.Time, amount float64) {
	daily := statdao.GetDailyLoginStatByDate(entity.FormatDailyLoginStatDate(statAt))
	daily.AddGoldConsumeAmount(amount)
	statdao.PublishDailyLoginStat(daily)

	weekly := statdao.GetWeeklyLoginStatByWeek(entity.FormatWeeklyLoginStatKey(statAt))
	weekly.AddGoldConsumeAmount(amount)
	statdao.PublishWeeklyLoginStat(weekly)

	monthly := statdao.GetMonthlyLoginStatByMonth(entity.FormatMonthlyLoginStatKey(statAt))
	monthly.AddGoldConsumeAmount(amount)
	statdao.PublishMonthlyLoginStat(monthly)
}

func recordPeriodDiamondConsume(statAt time.Time, amount float64) {
	daily := statdao.GetDailyLoginStatByDate(entity.FormatDailyLoginStatDate(statAt))
	daily.AddDiamondConsumeAmount(amount)
	statdao.PublishDailyLoginStat(daily)

	weekly := statdao.GetWeeklyLoginStatByWeek(entity.FormatWeeklyLoginStatKey(statAt))
	weekly.AddDiamondConsumeAmount(amount)
	statdao.PublishWeeklyLoginStat(weekly)

	monthly := statdao.GetMonthlyLoginStatByMonth(entity.FormatMonthlyLoginStatKey(statAt))
	monthly.AddDiamondConsumeAmount(amount)
	statdao.PublishMonthlyLoginStat(monthly)
}
