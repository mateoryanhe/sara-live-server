package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// OfflineGameRoomPlayTrackingEvent 离线游戏房开始玩游戏(每次 +1)
	OfflineGameRoomPlayTrackingEvent event.Type = "OfflineGameRoomPlayTrackingEvent"
)

type OfflineGameRoomPlayTrackingEventData struct {
	UserId   uint64
	RoomId   uint64
	GameCode string
	StatAt   time.Time
}

func NewOfflineGameRoomPlayTrackingEventData(userId, roomId uint64, gameCode string, statAt time.Time) *OfflineGameRoomPlayTrackingEventData {
	return &OfflineGameRoomPlayTrackingEventData{
		UserId:   userId,
		RoomId:   roomId,
		GameCode: gameCode,
		StatAt:   statAt,
	}
}
