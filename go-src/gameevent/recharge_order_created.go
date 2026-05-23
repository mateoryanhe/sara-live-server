package gameevent

import "xr-game-server/core/event"

const (
	// RechargeOrderCreatedEvent 充值订单创建成功事件
	RechargeOrderCreatedEvent event.Type = "RechargeOrderCreatedEvent"
)

// RechargeOrderCreatedEventData 充值订单创建成功事件数据
type RechargeOrderCreatedEventData struct {
	OrderId uint64 // 充值订单ID
}

func NewRechargeOrderCreatedEventData(orderId uint64) *RechargeOrderCreatedEventData {
	return &RechargeOrderCreatedEventData{OrderId: orderId}
}
