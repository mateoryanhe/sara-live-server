package stat

import (
	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

func initPromotionChannelOrderEvent() {
	event.Sub(gameevent.RechargeOrderCreatedEvent, onRechargeOrderCreatedForPromotionChannel)
}

func onRechargeOrderCreatedForPromotionChannel(data any) {
	enqueue(&statJob{Kind: jobPromotionChannelOrder, Payload: data})
}
