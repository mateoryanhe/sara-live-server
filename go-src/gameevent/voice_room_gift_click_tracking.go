package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// VoiceRoomGiftClickTrackingEvent 语聊房礼物面板点击(每次 +1)
	VoiceRoomGiftClickTrackingEvent event.Type = "VoiceRoomGiftClickTrackingEvent"
)

type VoiceRoomGiftClickTrackingEventData struct {
	UserId uint64
	StatAt time.Time
}

func NewVoiceRoomGiftClickTrackingEventData(userId uint64, statAt time.Time) *VoiceRoomGiftClickTrackingEventData {
	return &VoiceRoomGiftClickTrackingEventData{
		UserId: userId,
		StatAt: statAt,
	}
}
