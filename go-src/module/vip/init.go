package vip

import (
	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

// Init 服务启动时加载VIP配置到内存
func Init() {
	reloadVipCfgMemory()
	// Init 订阅充值到账事件,判断VIP升级
	event.Sub(gameevent.RechargeArrivedEvent, onRechargeArrived)
}
