package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// ShowcaseLiveRoomJoinTrackingEvent 观众成功 joinRoom 后的秀场进房次数埋点(每次 +1,由 tracking 模块消费)
	ShowcaseLiveRoomJoinTrackingEvent event.Type = "ShowcaseLiveRoomJoinTrackingEvent"
)

// ShowcaseLiveRoomJoinTrackingEventData joinRoom 成功埋点载荷
type ShowcaseLiveRoomJoinTrackingEventData struct {
	UserId uint64
	RoomId uint64
	StatAt time.Time
}

func NewShowcaseLiveRoomJoinTrackingEventData(userId, roomId uint64, statAt time.Time) *ShowcaseLiveRoomJoinTrackingEventData {
	return &ShowcaseLiveRoomJoinTrackingEventData{
		UserId: userId,
		RoomId: roomId,
		StatAt: statAt,
	}
}
