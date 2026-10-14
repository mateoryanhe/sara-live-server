package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// OneToOneRoomListTrackingEvent App 查询 1v1 房间列表成功(每次 +1,由 tracking 模块消费)
	OneToOneRoomListTrackingEvent event.Type = "OneToOneRoomListTrackingEvent"
)

type OneToOneRoomListTrackingEventData struct {
	UserId uint64
	StatAt time.Time
}

func NewOneToOneRoomListTrackingEventData(userId uint64, statAt time.Time) *OneToOneRoomListTrackingEventData {
	return &OneToOneRoomListTrackingEventData{
		UserId: userId,
		StatAt: statAt,
	}
}
