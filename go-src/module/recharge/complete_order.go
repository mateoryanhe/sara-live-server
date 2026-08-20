package recharge

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/constants/currency"
	"xr-game-server/core/event"
	"xr-game-server/dao/rechargeorderdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/entity/recharge"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
	"xr-game-server/module/wallet"
)

func rechargeOrderLockKey(orderId uint64) string {
	return fmt.Sprintf("recharge_order_%d", orderId)
}

func rechargeThirdOrderLockKey(thirdOrderId string) string {
	return fmt.Sprintf("recharge_third_%s", thirdOrderId)
}

// completeOrder 内部统一的"订单完成 → 发放金币"逻辑
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
	after, err := wallet.GoldAdd(order.UserId, order.Gold, reason)
	if err != nil {
		return 0, err
	}
	paidAt := time.Now()
	order.SetStatus(entity.RechargeOrderStatusCompleted)
	order.SetPaidAt(paidAt)
	order.SetUpdatedAt(paidAt)
	if userinfodao.MarkFirstRechargeDone(order.UserId) {
		event.Pub(gameevent.FirstRechargeCompletedEvent, order)
	}

	CancelRechargeOrderTimeout(order.ID)
	event.Pub(gameevent.RechargeArrivedEvent, order)
	rechargeorderdao.FlushOrderCache(order)
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
