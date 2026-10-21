package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// MiniGamePlusClickTrackingEvent 游戏+(每次 +1)
	MiniGamePlusClickTrackingEvent event.Type = "MiniGamePlusClickTrackingEvent"
)

type MiniGamePlusClickTrackingEventData struct {
	UserId uint64
	StatAt time.Time
}

func NewMiniGamePlusClickTrackingEventData(userId uint64, statAt time.Time) *MiniGamePlusClickTrackingEventData {
	return &MiniGamePlusClickTrackingEventData{
		UserId: userId,
		StatAt: statAt,
	}
}
