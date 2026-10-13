package metapixel

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/core/xrpool"
	"xr-game-server/dao/rechargeorderdao"
	"xr-game-server/dao/userinfodao"
	rechargeentity "xr-game-server/entity/recharge"
	"xr-game-server/gameevent"
)

func initEvents() {
	event.Sub(gameevent.RegisterEvent, onRegisterEvent)
	event.Sub(gameevent.RechargeOrderCreatedEvent, onRechargeOrderCreatedEvent)
}

func onRegisterEvent(data any) {
	val, ok := data.(*gameevent.RegisterEventData)
	if !ok || val == nil || val.UserId == 0 {
		return
	}
	registeredAt := val.RegisteredAt
	userId := val.UserId
	attr := val.Attribution
	xrpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		trackCompleteRegistration(ctx, userId, registeredAt, attr)
	})
}

func onRechargeOrderCreatedEvent(data any) {
	val, ok := data.(*gameevent.RechargeOrderCreatedEventData)
	if !ok || val == nil || val.OrderId == 0 {
		return
	}
	orderId := val.OrderId
	xrpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		order := rechargeorderdao.GetById(orderId)
		if !shouldTrackMetaPurchaseOnOrderCreated(order) {
			return
		}
		at := order.CreatedAt
		if at.IsZero() {
			at = time.Now()
		}
		trackPurchase(ctx, order.UserId, order.ID, order.Price, at, val.Attribution, val.Checkout)
	})
}

// shouldTrackMetaPurchaseOnOrderCreated 用户 App/H5/渠道下单即上报 Purchase;排除币商/后台单/白名单测试号.
func shouldTrackMetaPurchaseOnOrderCreated(order *rechargeentity.RechargeOrder) bool {
	if order == nil || order.UserId == 0 {
		return false
	}
	if order.Source == rechargeentity.RechargeOrderSourceManual ||
		order.Source == rechargeentity.RechargeOrderSourceCoinMerchant {
		return false
	}
	if order.PayChannel == rechargeentity.RechargeCfgTypeCoinMerchant {
		return false
	}
	if userinfodao.IsRechargeWhitelist(order.UserId) {
		return false
	}
	return order.Price > 0
}
