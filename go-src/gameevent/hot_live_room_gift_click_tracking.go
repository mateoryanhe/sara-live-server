package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// HotLiveRoomGiftClickTrackingEvent 秀场直播间礼物面板点击(每次 +1)
	HotLiveRoomGiftClickTrackingEvent event.Type = "HotLiveRoomGiftClickTrackingEvent"
)

type HotLiveRoomGiftClickTrackingEventData struct {
	UserId uint64
	StatAt time.Time
}

func NewHotLiveRoomGiftClickTrackingEventData(userId uint64, statAt time.Time) *HotLiveRoomGiftClickTrackingEventData {
	return &HotLiveRoomGiftClickTrackingEventData{
		UserId: userId,
		StatAt: statAt,
	}
}
