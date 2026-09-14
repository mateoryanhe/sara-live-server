package recharge

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/constants/currency"
	"xr-game-server/core/xrlog"
	"xr-game-server/dao/rechargeorderdao"
	"xr-game-server/entity/recharge"
	"xr-game-server/errercode"
)

// CompleteChannelPayOrder 渠道/币商支付成功后发币（HaiPay 以查单确认成功后调用）
func CompleteChannelPayOrder(ctx context.Context, orderIdStr, thirdOrderId string) error {
	orderId, err := strconv.ParseUint(strings.TrimSpace(orderIdStr), 10, 64)
	if err != nil || orderId == 0 {
		return errercode.CreateCode(errercode.RechargeOrderNonExist)
	}
	thirdOrderId = strings.TrimSpace(thirdOrderId)
	if thirdOrderId != "" {
		gmlock.Lock(rechargeThirdOrderLockKey(thirdOrderId))
		defer gmlock.Unlock(rechargeThirdOrderLockKey(thirdOrderId))
		if existing := rechargeorderdao.GetByThirdOrderId(thirdOrderId); existing != nil {
			if existing.ID == orderId && existing.Status == entity.RechargeOrderStatusCompleted {
				return nil
			}
			if existing.ID != orderId {
				xrlog.DetailLog.Warningf(ctx, "channelPay thirdOrder reused third=%s orderId=%d existing=%d", thirdOrderId, orderId, existing.ID)
				return errercode.CreateCode(errercode.RechargeOrderStateInvalid)
			}
		}
	}

	order := rechargeorderdao.GetById(orderId)
	if order == nil {
		return errercode.CreateCode(errercode.RechargeOrderNonExist)
	}
	if order.Status == entity.RechargeOrderStatusCompleted {
		return nil
	}
	if order.PayChannel != entity.RechargeCfgTypeChannel && order.PayChannel != entity.RechargeCfgTypeCoinMerchant {
		return errercode.CreateCode(errercode.RechargeOrderStateInvalid)
	}
	if thirdOrderId != "" && order.ThirdOrderId != thirdOrderId {
		order.SetThirdOrderId(thirdOrderId)
	}
	_, err = completeOrder(order, currency.ReasonRecharge)
	return err
}

// FailChannelPayOrder HaiPay 查单非成功时，将待支付渠道单标为失败。已完成/已取消/已失败幂等。
func FailChannelPayOrder(ctx context.Context, orderIdStr, thirdOrderId string, platformStatus int) error {
	orderId, err := strconv.ParseUint(strings.TrimSpace(orderIdStr), 10, 64)
	if err != nil || orderId == 0 {
		return errercode.CreateCode(errercode.RechargeOrderNonExist)
	}
	thirdOrderId = strings.TrimSpace(thirdOrderId)

	lockKey := rechargeOrderLockKey(orderId)
	gmlock.Lock(lockKey)
	defer gmlock.Unlock(lockKey)

	order := rechargeorderdao.GetById(orderId)
	if order == nil {
		return errercode.CreateCode(errercode.RechargeOrderNonExist)
	}
	if order.Status == entity.RechargeOrderStatusCompleted {
		return nil
	}
	if order.Status == entity.RechargeOrderStatusCancelled || order.Status == entity.RechargeOrderStatusFailed {
		return nil
	}
	if order.PayChannel != entity.RechargeCfgTypeChannel && order.PayChannel != entity.RechargeCfgTypeCoinMerchant {
		return errercode.CreateCode(errercode.RechargeOrderStateInvalid)
	}
	if thirdOrderId != "" && order.ThirdOrderId != thirdOrderId {
		order.SetThirdOrderId(thirdOrderId)
	}
	now := time.Now()
	order.SetStatus(entity.RechargeOrderStatusFailed)
	order.SetUpdatedAt(now)
	if remark := strings.TrimSpace(order.Remark); remark == "" {
		order.SetRemark(fmt.Sprintf("haipay status=%d", platformStatus))
	}
	CancelRechargeOrderTimeout(order.ID)
	rechargeorderdao.FlushOrderCache(order)
	xrlog.DetailLog.Infof(ctx, "channelPay mark failed orderId=%d platformStatus=%d third=%s", orderId, platformStatus, thirdOrderId)
	return nil
}
