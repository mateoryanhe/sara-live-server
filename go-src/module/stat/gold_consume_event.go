package stat

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/core/event"
	"xr-game-server/dao/statdao"
	"xr-game-server/entity"
	"xr-game-server/gameevent"
)

func initGoldConsumeEvent() {
	event.Sub(gameevent.CurrencyChangeEvent, onGoldConsumeEvent)
}

func onGoldConsumeEvent(data any) {
	ev, ok := data.(*gameevent.CurrencyChangeEventData)
	if !ok || ev == nil {
		g.Log().Errorf(gctx.New(), "CurrencyChangeEvent payload type error: %T", data)
		return
	}
	if ev.Type != gameevent.CurrencyTypeGold || ev.Action != gameevent.CurrencyActionSub || ev.Amount <= 0 {
		return
	}
	if !shouldCountUserStat(ev.UserId) {
		return
	}

	statAt := time.Now()

	lockName := "stat_gold_consume"
	gmlock.Lock(lockName)
	defer gmlock.Unlock(lockName)

	if stat := statdao.GetSysStat(); stat != nil {
		stat.AddTotalGoldConsume(ev.Amount)
	}
	recordPeriodGoldConsume(statAt, ev.Amount)
	recordPeriodGoldConsumeUser(statAt, ev.UserId)
}

func recordPeriodGoldConsume(statAt time.Time, amount float64) {
	daily := statdao.GetDailyLoginStatByDate(entity.FormatDailyLoginStatDate(statAt))
	daily.AddGoldConsumeAmount(amount)

	weekly := statdao.GetWeeklyLoginStatByWeek(entity.FormatWeeklyLoginStatKey(statAt))
	weekly.AddGoldConsumeAmount(amount)

	monthly := statdao.GetMonthlyLoginStatByMonth(entity.FormatMonthlyLoginStatKey(statAt))
	monthly.AddGoldConsumeAmount(amount)
}
