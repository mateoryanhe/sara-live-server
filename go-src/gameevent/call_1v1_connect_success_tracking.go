package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// Call1v1ConnectSuccessTrackingEvent 双方 confirmCall 完成且首次扣费成功(进入通话中),每次 +1
	Call1v1ConnectSuccessTrackingEvent event.Type = "Call1v1ConnectSuccessTrackingEvent"
)

// Call1v1ConnectSuccessTrackingEventData 1v1 视频通话双方成功接通埋点
type Call1v1ConnectSuccessTrackingEventData struct {
	Source     uint8
	CallType   uint8
	CallerId   uint64
	ReceiverId uint64
	StatAt     time.Time
}

func NewCall1v1ConnectSuccessTrackingEventData(source, callType uint8, callerId, receiverId uint64, statAt time.Time) *Call1v1ConnectSuccessTrackingEventData {
	return &Call1v1ConnectSuccessTrackingEventData{
		Source:     source,
		CallType:   callType,
		CallerId:   callerId,
		ReceiverId: receiverId,
		StatAt:     statAt,
	}
}
