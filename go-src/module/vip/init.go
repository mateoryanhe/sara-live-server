package vip

import (
	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

// Init 服务启动时加载VIP配置到内存
func Init() {
	reloadVipCfgMemory()
	event.Sub(gameevent.RechargeGoldArrivedEvent, onRechargeGoldArrived)
	event.Sub(gameevent.UsdIncomeArrivedEvent, onUsdIncomeArrived)
}
