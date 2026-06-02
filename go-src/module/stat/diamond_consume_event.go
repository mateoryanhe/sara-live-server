package stat

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/core/event"
	"xr-game-server/dao/dailyloginstatdao"
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
	if ev.Type != gameevent.CurrencyTypeDiamond || ev.Action != gameevent.CurrencyActionSub || ev.Amount <= 0 {
		return
	}

	statAt := time.Now()

	lockName := "stat_diamond_consume"
	gmlock.Lock(lockName)
	defer gmlock.Unlock(lockName)

	if stat := statdao.GetSysStat(); stat != nil {
		stat.AddTotalDiamondConsume(ev.Amount)
	}
	recordPeriodDiamondConsume(statAt, ev.Amount)
	recordPeriodDiamondConsumeUser(statAt, ev.UserId)
}

func recordPeriodDiamondConsume(statAt time.Time, amount float64) {
	daily := dailyloginstatdao.GetByDate(entity.FormatDailyLoginStatDate(statAt))
	daily.AddDiamondConsumeAmount(amount)

	weekly := weeklyloginstatdao.GetByWeek(entity.FormatWeeklyLoginStatKey(statAt))
	weekly.AddDiamondConsumeAmount(amount)

	monthly := monthlyloginstatdao.GetByMonth(entity.FormatMonthlyLoginStatKey(statAt))
	monthly.AddDiamondConsumeAmount(amount)
}
