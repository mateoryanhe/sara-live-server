package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// MiniGameSpinClickTrackingEvent Spin 点击(每次 +1)
	MiniGameSpinClickTrackingEvent event.Type = "MiniGameSpinClickTrackingEvent"
)

type MiniGameSpinClickTrackingEventData struct {
	UserId uint64
	StatAt time.Time
}

func NewMiniGameSpinClickTrackingEventData(userId uint64, statAt time.Time) *MiniGameSpinClickTrackingEventData {
	return &MiniGameSpinClickTrackingEventData{
		UserId: userId,
		StatAt: statAt,
	}
}
