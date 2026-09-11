package vip

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/gameevent"
)

// onUsdIncomeArrived 用户累计充值(USD):订阅美金入账,含虚拟美金(白名单);须在 wallet 发币成功后
func onUsdIncomeArrived(val any) {
	data, ok := val.(*gameevent.UsdIncomeArrivedEventData)
	if !ok || data == nil || data.Order == nil {
		g.Log().Errorf(gctx.New(), "UsdIncomeArrivedEvent payload type error: %T", val)
		return
	}
	if !data.Granted {
		return
	}
	amount := data.UsdAmount
	if amount <= 0 {
		amount = data.Order.Price
	}
	if amount <= 0 || data.Order.UserId == 0 {
		return
	}
	stat := userinfodao.GetUserCumulativeStatByUserId(data.Order.UserId)
	stat.AddTotalRecharge(amount)
	userinfodao.PublishUserCumulativeStat(stat)
}
