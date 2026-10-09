package metapixel

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/core/xrpool"
	rechargeentity "xr-game-server/entity/recharge"
	"xr-game-server/gameevent"
)

func initEvents() {
	event.Sub(gameevent.RegisterEvent, onRegisterEvent)
	event.Sub(gameevent.RechargeGoldArrivedEvent, onRechargeGoldArrivedEvent)
}

func onRegisterEvent(data any) {
	val, ok := data.(*gameevent.RegisterEventData)
	if !ok || val == nil || val.UserId == 0 {
		return
	}
	registeredAt := val.RegisteredAt
	userId := val.UserId
	xrpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		trackCompleteRegistration(ctx, userId, registeredAt)
	})
}

func onRechargeGoldArrivedEvent(data any) {
	val, ok := data.(*gameevent.RechargeGoldArrivedEventData)
	if !ok || val == nil || val.Order == nil {
		return
	}
	if !isNormalUserRealUsdRecharge(val) {
		return
	}
	order := val.Order
	usd := order.Price
	paidAt := order.PaidAt
	if paidAt.IsZero() {
		paidAt = time.Now()
	}
	userId := order.UserId
	orderId := order.ID
	xrpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		trackPurchase(ctx, userId, orderId, usd, paidAt)
	})
}

func isNormalUserRealUsdRecharge(val *gameevent.RechargeGoldArrivedEventData) bool {
	if val == nil || val.Order == nil || val.Order.UserId == 0 {
		return false
	}
	if val.Kind == gameevent.UsdIncomeKindCoinMerchant || val.Kind == gameevent.UsdIncomeKindWhitelist {
		return false
	}
	if val.Order.PayChannel == rechargeentity.RechargeCfgTypeCoinMerchant {
		return false
	}
	return val.Order.Price > 0 || val.Order.PayAmount > 0
}
