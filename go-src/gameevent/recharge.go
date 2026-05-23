package gameevent

import "xr-game-server/core/event"

const (
	// RechargeArrivedEvent 充值到账事件(订单完成且金币已发放)
	RechargeArrivedEvent event.Type = "RechargeArrivedEvent"
)

// RechargeArrivedEventData 充值到账事件数据
type RechargeArrivedEventData struct {
	OrderId uint64 // 充值订单ID
}

func NewRechargeArrivedEventData(orderId uint64) *RechargeArrivedEventData {
	return &RechargeArrivedEventData{OrderId: orderId}
}
