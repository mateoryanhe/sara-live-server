package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// MiniGameControllerClickTrackingEvent 点击游戏手柄(每次 +1)
	MiniGameControllerClickTrackingEvent event.Type = "MiniGameControllerClickTrackingEvent"
)

type MiniGameControllerClickTrackingEventData struct {
	UserId uint64
	StatAt time.Time
}

func NewMiniGameControllerClickTrackingEventData(userId uint64, statAt time.Time) *MiniGameControllerClickTrackingEventData {
	return &MiniGameControllerClickTrackingEventData{
		UserId: userId,
		StatAt: statAt,
	}
}
