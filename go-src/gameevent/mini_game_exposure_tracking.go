package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// MiniGameExposureTrackingEvent 半屏游戏窗口曝光(每次 +1)
	MiniGameExposureTrackingEvent event.Type = "MiniGameExposureTrackingEvent"
)

type MiniGameExposureTrackingEventData struct {
	UserId uint64
	StatAt time.Time
}

func NewMiniGameExposureTrackingEventData(userId uint64, statAt time.Time) *MiniGameExposureTrackingEventData {
	return &MiniGameExposureTrackingEventData{
		UserId: userId,
		StatAt: statAt,
	}
}
