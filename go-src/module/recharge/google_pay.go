package recharge

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

func initGooglePay() {
	event.Sub(gameevent.RechargeOrderCreatedEvent, onGooglePayCreateOrder)
}

func onGooglePayCreateOrder(data any) {
	g.Log().Infof(gctx.New(), "google pay create order%v", data)
}
