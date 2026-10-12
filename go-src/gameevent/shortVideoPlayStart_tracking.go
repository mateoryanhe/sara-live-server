package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// ShortVideoPlayStartTrackingEvent 半屏游戏窗口曝光(每次 +1)
	ShortVideoPlayStartTrackingEvent event.Type = "ShortVideoPlayStartTrackingEvent"
)

type ShortVideoPlayStartTrackingEventData struct {
	UserId uint64
	StatAt time.Time
}

func NewShortVideoPlayStartTrackingEventData(userId uint64, statAt time.Time) *ShortVideoPlayStartTrackingEventData {
	return &ShortVideoPlayStartTrackingEventData{
		UserId: userId,
		StatAt: statAt,
	}
}
