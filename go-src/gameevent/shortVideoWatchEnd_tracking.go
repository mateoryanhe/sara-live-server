package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// ShortVideoWatchEndTrackingEvent 半屏游戏窗口曝光(每次 +1)
	ShortVideoWatchEndTrackingEvent event.Type = "ShortVideoWatchEndTrackingEvent"
)

type ShortVideoWatchEndTrackingEventData struct {
	UserId uint64
	StatAt time.Time
}

func NewShortVideoWatchEndTrackingEventData(userId uint64, statAt time.Time) *ShortVideoWatchEndTrackingEventData {
	return &ShortVideoWatchEndTrackingEventData{
		UserId: userId,
		StatAt: statAt,
	}
}
