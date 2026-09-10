package recharge

import (
	"context"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/constants/currency"
	"xr-game-server/core/xrlog"
	"xr-game-server/dao/rechargeorderdao"
	"xr-game-server/entity/recharge"
	"xr-game-server/errercode"
)

// CompleteChannelPayOrder 渠道/币商支付成功后发币（各 Provider 验签通过后调用）
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
