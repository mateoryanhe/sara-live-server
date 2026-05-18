package currencylog

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/entity"
	"xr-game-server/gameevent"
)

// Init 订阅货币变动事件,落流水
func Init() {
	event.Sub(gameevent.CurrencyChangeEvent, onCurrencyChange)
}

func onCurrencyChange(val any) {
	data, ok := val.(*gameevent.CurrencyChangeEventData)
	if !ok || data == nil {
		g.Log().Errorf(gctx.New(), "CurrencyChangeEvent payload type error: %T", val)
		return
	}
	_ = entity.NewCurrencyLog(
		data.UserId,
		data.Type,
		data.Action,
		data.Amount,
		data.Before,
		data.After,
		data.Reason,
	)
	g.Log().Debugf(context.Background(),
		"currency log userId=%d type=%d action=%d amount=%v before=%v after=%v reason=%s",
		data.UserId, data.Type, data.Action, data.Amount, data.Before, data.After, data.Reason,
	)
}
