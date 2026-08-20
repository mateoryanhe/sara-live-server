package activity

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity/recharge"
)

func onFirstRechargeCompleted(val any) {
	order, ok := val.(*entity.RechargeOrder)
	if !ok || order == nil {
		g.Log().Errorf(gctx.New(), "FirstRechargeCompletedEvent payload type error: %T", val)
		return
	}
	pushFirstRechargeSuccessToApp(order.UserId, order)
}
