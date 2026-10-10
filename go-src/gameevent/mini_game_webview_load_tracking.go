package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// MiniGameWebViewLoadTrackingEvent 游戏 WebView 成功加载(每次 +1)
	MiniGameWebViewLoadTrackingEvent event.Type = "MiniGameWebViewLoadTrackingEvent"
)

type MiniGameWebViewLoadTrackingEventData struct {
	UserId uint64
	StatAt time.Time
}

func NewMiniGameWebViewLoadTrackingEventData(userId uint64, statAt time.Time) *MiniGameWebViewLoadTrackingEventData {
	return &MiniGameWebViewLoadTrackingEventData{
		UserId: userId,
		StatAt: statAt,
	}
}
