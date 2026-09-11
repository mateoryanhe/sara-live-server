package recharge

import (
	"fmt"

	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/constants/currency"
	"xr-game-server/core/event"
	"xr-game-server/dao/rechargeorderdao"
	"xr-game-server/entity/recharge"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
)

func rechargeOrderLockKey(orderId uint64) string {
	return fmt.Sprintf("recharge_order_%d", orderId)
}

func rechargeThirdOrderLockKey(thirdOrderId string) string {
	return fmt.Sprintf("recharge_third_%s", thirdOrderId)
}

// completeOrder 订单完成 → 发布美金入账 → wallet 发币加赠 → 推送
// 幂等:已经是已完成状态的订单不会重复发放
// 返回(发放后玩家金币余额, 错误)
func completeOrder(o *entity.RechargeOrder, reason currency.Reason) (float64, error) {
	if o == nil || o.ID == 0 {
		return 0, errercode.CreateCode(errercode.RechargeOrderNonExist)
	}
	lockKey := rechargeOrderLockKey(o.ID)
	gmlock.Lock(lockKey)
	defer gmlock.Unlock(lockKey)

	order := rechargeorderdao.GetById(o.ID)
	if order == nil {
		return 0, errercode.CreateCode(errercode.RechargeOrderNonExist)
	}
	if order.Status == entity.RechargeOrderStatusCompleted {
		return 0, errercode.CreateCode(errercode.RechargeOrderStateInvalid)
	}
	if order.Gold <= 0 {
		return 0, errercode.CreateCode(errercode.RechargeGoldInvalid)
	}

	usdData := gameevent.NewUsdIncomeArrivedEventData(order, resolveUsdIncomeKind(order), reason)
	event.Pub(gameevent.UsdIncomeArrivedEvent, usdData)
	if !usdData.Granted {
		return 0, errercode.CreateCode(errercode.RechargeGoldInvalid)
	}

	CancelRechargeOrderTimeout(order.ID)
	rechargeorderdao.FlushOrderCache(order)
	pushRechargeSuccessToApp(order.UserId, order, usdData.GoldBalance)
	return usdData.GoldBalance, nil
}

// CompleteOrder 对外:支付回调成功时调用此函数,完成订单并发放金币
func CompleteOrder(orderId uint64) (float64, error) {
	o := rechargeorderdao.GetById(orderId)
	if o == nil {
		return 0, errercode.CreateCode(errercode.RechargeOrderNonExist)
	}
	return completeOrder(o, currency.ReasonRecharge)
}
