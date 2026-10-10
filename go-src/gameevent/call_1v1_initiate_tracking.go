package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// Call1v1InitiateTrackingEvent 秀场直播间 source=1 视频通话发起成功(每次 +1,由 tracking 模块消费)
	Call1v1InitiateTrackingEvent event.Type = "Call1v1InitiateTrackingEvent"
)

// Call1v1InitiateTrackingEventData 1v1 视频通话发起埋点(call_order.source=1)
type Call1v1InitiateTrackingEventData struct {
	CallerId uint64
	AnchorId uint64
	StatAt   time.Time
}

func NewCall1v1InitiateTrackingEventData(callerId, anchorId uint64, statAt time.Time) *Call1v1InitiateTrackingEventData {
	return &Call1v1InitiateTrackingEventData{
		CallerId: callerId,
		AnchorId: anchorId,
		StatAt:   statAt,
	}
}
