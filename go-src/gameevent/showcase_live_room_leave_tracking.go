package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// ShowcaseLiveRoomLeaveTrackingEvent 观众 leaveRoom 成功后的秀场退房次数埋点(每次 +1,由 tracking 模块消费)
	ShowcaseLiveRoomLeaveTrackingEvent event.Type = "ShowcaseLiveRoomLeaveTrackingEvent"
)

// ShowcaseLiveRoomLeaveTrackingEventData leaveRoom 成功埋点载荷
type ShowcaseLiveRoomLeaveTrackingEventData struct {
	UserId uint64
	RoomId uint64
	StatAt time.Time
}

func NewShowcaseLiveRoomLeaveTrackingEventData(userId, roomId uint64, statAt time.Time) *ShowcaseLiveRoomLeaveTrackingEventData {
	return &ShowcaseLiveRoomLeaveTrackingEventData{
		UserId: userId,
		RoomId: roomId,
		StatAt: statAt,
	}
}
