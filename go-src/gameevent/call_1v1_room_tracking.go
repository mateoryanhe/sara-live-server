package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// Call1v1RoomTrackingEvent 1v1 房间 source=3 视频通话发起成功(每次 +1,由 tracking 模块消费)
	Call1v1RoomTrackingEvent event.Type = "Call1v1RoomTrackingEvent"
)

// Call1v1RoomTrackingEventData 1v1 房间视频通话发起埋点(call_order.source=3)
type Call1v1RoomTrackingEventData struct {
	CallerId uint64
	AnchorId uint64
	StatAt   time.Time
}

func NewCall1v1RoomTrackingEventData(callerId, anchorId uint64, statAt time.Time) *Call1v1RoomTrackingEventData {
	return &Call1v1RoomTrackingEventData{
		CallerId: callerId,
		AnchorId: anchorId,
		StatAt:   statAt,
	}
}
