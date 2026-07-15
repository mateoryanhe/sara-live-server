package stat

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/constants/currency"
	"xr-game-server/core/event"
	"xr-game-server/dao/monthlyloginstatdao"
	"xr-game-server/dao/statdao"
	"xr-game-server/dao/weeklyloginstatdao"
	"xr-game-server/entity"
	"xr-game-server/gameevent"
)

func initDiamondConsumeEvent() {
	event.Sub(gameevent.CurrencyChangeEvent, onDiamondConsumeEvent)
}

func onDiamondConsumeEvent(data any) {
	ev, ok := data.(*gameevent.CurrencyChangeEventData)
	if !ok || ev == nil {
		g.Log().Errorf(gctx.New(), "CurrencyChangeEvent payload type error: %T", data)
		return
	}
	if ev.Type != gameevent.CurrencyTypeDiamond || ev.Amount <= 0 {
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

	statAt := time.Now()

	lockName := "stat_diamond_consume"
	gmlock.Lock(lockName)
	defer gmlock.Unlock(lockName)

	if stat := statdao.GetSysStat(); stat != nil {
		stat.AddTotalDiamondConsume(delta)
	}
	recordPeriodDiamondConsume(statAt, delta)
	if delta > 0 {
		recordPeriodDiamondConsumeUser(statAt, ev.UserId)
	}
}

func recordPeriodDiamondConsume(statAt time.Time, amount float64) {
	daily := statdao.GetDailyLoginStatByDate(entity.FormatDailyLoginStatDate(statAt))
	daily.AddDiamondConsumeAmount(amount)

	weekly := weeklyloginstatdao.GetByWeek(entity.FormatWeeklyLoginStatKey(statAt))
	weekly.AddDiamondConsumeAmount(amount)

	monthly := monthlyloginstatdao.GetByMonth(entity.FormatMonthlyLoginStatKey(statAt))
	monthly.AddDiamondConsumeAmount(amount)
}
