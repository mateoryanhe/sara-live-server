package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// ValidAudienceEvent 有效观众(直播中进房且非主播,日/周/月跨房去重)
	ValidAudienceEvent event.Type = "ValidAudienceEvent"
)

// ValidAudienceEventData 有效观众事件载荷
type ValidAudienceEventData struct {
	UserId uint64
	StatAt time.Time
}

func NewValidAudienceEventData(userId uint64, statAt time.Time) *ValidAudienceEventData {
	return &ValidAudienceEventData{
		UserId: userId,
		StatAt: statAt,
	}
}
