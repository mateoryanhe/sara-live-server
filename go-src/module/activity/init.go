package activity

import (
	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

func Init() {
	ReloadFirstRechargeActivityCache()
	ReloadInviteRechargeRewardCache()
	event.Sub(gameevent.RechargeGoldArrivedEvent, onRechargeGoldArrivedForFirstRecharge)
	event.Sub(gameevent.RechargeGoldArrivedEvent, onInviteRechargeGoldArrived)
}
