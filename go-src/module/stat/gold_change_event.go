package stat

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/core/event"
	"xr-game-server/dao/statdao"
	"xr-game-server/gameevent"
)

func initGoldChangeEvent() {
	event.Sub(gameevent.CurrencyChangeEvent, onGoldChangeEvent)
}

// onGoldChangeEvent 监听金币变动,加减均累计到系统总金币
func onGoldChangeEvent(data any) {
	ev, ok := data.(*gameevent.CurrencyChangeEventData)
	if !ok || ev == nil {
		g.Log().Errorf(gctx.New(), "CurrencyChangeEvent payload type error: %T", data)
		return
	}
	if ev.Type != gameevent.CurrencyTypeGold || ev.Amount <= 0 {
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

	lockName := "stat_gold_change"
	gmlock.Lock(lockName)
	defer gmlock.Unlock(lockName)

	stat := statdao.GetSysStat()
	if stat == nil {
		return
	}
	stat.AddTotalGold(delta)
}
