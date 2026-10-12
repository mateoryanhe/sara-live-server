package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// ShortVideoPaySuccessTrackingEvent 半屏游戏窗口曝光(每次 +1)
	ShortVideoPaySuccessTrackingEvent event.Type = "ShortVideoPaySuccessTrackingEvent"
)

type ShortVideoPaySuccessTrackingEventData struct {
	UserId uint64
	StatAt time.Time
}

func NewShortVideoPaySuccessTrackingEventData(userId uint64, statAt time.Time) *ShortVideoPaySuccessTrackingEventData {
	return &ShortVideoPaySuccessTrackingEventData{
		UserId: userId,
		StatAt: statAt,
	}
}
