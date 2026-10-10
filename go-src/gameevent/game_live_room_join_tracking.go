package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// GameLiveRoomJoinTrackingEvent 观众成功 joinRoom 后的游戏直播类进房次数埋点(每次 +1,由 tracking 模块消费)
	GameLiveRoomJoinTrackingEvent event.Type = "GameLiveRoomJoinTrackingEvent"
)

// GameLiveRoomJoinTrackingEventData joinRoom 成功埋点载荷
type GameLiveRoomJoinTrackingEventData struct {
	UserId uint64
	RoomId uint64
	StatAt time.Time
}

func NewGameLiveRoomJoinTrackingEventData(userId, roomId uint64, statAt time.Time) *GameLiveRoomJoinTrackingEventData {
	return &GameLiveRoomJoinTrackingEventData{
		UserId: userId,
		RoomId: roomId,
		StatAt: statAt,
	}
}
