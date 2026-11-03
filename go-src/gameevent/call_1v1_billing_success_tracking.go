package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// Call1v1BillingSuccessTrackingEvent 1v1 视频通话分钟计费扣款成功(含接通首分钟与续费),每次 +1
	Call1v1BillingSuccessTrackingEvent event.Type = "Call1v1BillingSuccessTrackingEvent"
)

// Call1v1BillingSuccessTrackingEventData 1v1 视频通话成功消费埋点
type Call1v1BillingSuccessTrackingEventData struct {
	Source     uint8
	CallType   uint8
	CallerId   uint64
	ReceiverId uint64
	OrderId    uint64
	StatAt     time.Time
}

func NewCall1v1BillingSuccessTrackingEventData(source, callType uint8, callerId, receiverId, orderId uint64, statAt time.Time) *Call1v1BillingSuccessTrackingEventData {
	return &Call1v1BillingSuccessTrackingEventData{
		Source:     source,
		CallType:   callType,
		CallerId:   callerId,
		ReceiverId: receiverId,
		OrderId:    orderId,
		StatAt:     statAt,
	}
}
