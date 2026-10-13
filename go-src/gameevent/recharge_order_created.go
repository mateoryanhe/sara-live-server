package gameevent

import (
	"context"
	"strings"

	"xr-game-server/core/event"
)

const (
	// RechargeOrderCreatedEvent 充值订单创建成功事件
	RechargeOrderCreatedEvent event.Type = "RechargeOrderCreatedEvent"
)

// RechargeCheckoutHint 渠道/H5 建单时用户填写的联系信息(服务端 CAPI 匹配用,可选).
type RechargeCheckoutHint struct {
	Email string
	Phone string
	Name  string
}

// RechargeOrderCreatedEventData 充值订单创建成功事件数据
type RechargeOrderCreatedEventData struct {
	OrderId     uint64
	Attribution RegisterAttribution
	Checkout    RechargeCheckoutHint
}

func NewRechargeOrderCreatedEventData(orderId uint64) *RechargeOrderCreatedEventData {
	return &RechargeOrderCreatedEventData{OrderId: orderId}
}

// NewRechargeOrderCreatedFromCtx 携带下单 HTTP 归因(IP/UA/Cookie fbp/fbc)与可选联系信息.
func NewRechargeOrderCreatedFromCtx(ctx context.Context, orderId uint64, checkout RechargeCheckoutHint) *RechargeOrderCreatedEventData {
	return &RechargeOrderCreatedEventData{
		OrderId:     orderId,
		Attribution: registerAttributionFromCtx(ctx, nil),
		Checkout: RechargeCheckoutHint{
			Email: strings.TrimSpace(checkout.Email),
			Phone: strings.TrimSpace(checkout.Phone),
			Name:  strings.TrimSpace(checkout.Name),
		},
	}
}
