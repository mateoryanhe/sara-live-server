package recharge

import (
	"time"
	"xr-game-server/constants/currency"
	"xr-game-server/core/event"
	"xr-game-server/dao/rechargeorderdao"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
	"xr-game-server/module/wallet"
)

// completeOrder 内部统一的"订单完成 → 发放金币"逻辑
// 幂等:已经是已完成状态的订单不会重复发放
// 返回(发放后玩家金币余额, 错误)
func completeOrder(o *entity.RechargeOrder, reason currency.Reason) (float64, error) {
	if o.Status == entity.RechargeOrderStatusCompleted {
		return 0, errercode.CreateCode(errercode.RechargeOrderStateInvalid)
	}
	if o.Gold <= 0 {
		return 0, errercode.CreateCode(errercode.RechargeGoldInvalid)
	}
	after, err := wallet.GoldAdd(o.UserId, o.Gold, reason)
	if err != nil {
		return 0, err
	}
	paidAt := time.Now()
	o.SetStatus(entity.RechargeOrderStatusCompleted)
	o.SetPaidAt(paidAt)
	o.SetUpdatedAt(paidAt)
	event.Pub(gameevent.RechargeArrivedEvent, gameevent.NewRechargeArrivedEventData(
		o.ID, o.UserId, o.CfgId, o.Price,
		o.Currency, o.Gold, o.Source,
		o.PayChannel, o.ThirdOrderId,
		after, paidAt,
	))
	return after, nil
}

// CompleteOrder 对外:支付回调成功时调用此函数,完成订单并发放金币
// (本次需求不开发回调路由,该函数保留以便后续接入第三方支付回调)
func CompleteOrder(orderId uint64) (float64, error) {
	o := rechargeorderdao.GetById(orderId)
	if o == nil {
		return 0, errercode.CreateCode(errercode.RechargeOrderNonExist)
	}
	return completeOrder(o, currency.ReasonRecharge)
}
