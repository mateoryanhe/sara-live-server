package activity

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/gameevent"
)

// onRechargeGoldArrivedForFirstRecharge 账号首充时推送 App 隐藏首充入口
func onRechargeGoldArrivedForFirstRecharge(val any) {
	data, ok := val.(*gameevent.RechargeGoldArrivedEventData)
	if !ok || data == nil || data.Order == nil {
		g.Log().Errorf(gctx.New(), "RechargeGoldArrivedEvent(first recharge) payload type error: %T", val)
		return
	}
	if !data.IsAccountFirst {
		return
	}
	pushFirstRechargeSuccessToApp(data.Order.UserId, data.Order)
}
