package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// ShortVideoPaywallClickTrackingEvent 短视频确认付费解锁(payShortVideo 未购时每次 +1)
	ShortVideoPaywallClickTrackingEvent event.Type = "ShortVideoPaywallClickTrackingEvent"
)

type ShortVideoPaywallClickTrackingEventData struct {
	UserId uint64
	StatAt time.Time
}

func NewShortVideoPaywallClickTrackingEventData(userId uint64, statAt time.Time) *ShortVideoPaywallClickTrackingEventData {
	return &ShortVideoPaywallClickTrackingEventData{
		UserId: userId,
		StatAt: statAt,
	}
}
