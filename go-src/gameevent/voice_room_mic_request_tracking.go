package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// VoiceRoomMicRequestTrackingEvent 语聊房申请/上麦成功(每次 +1)
	VoiceRoomMicRequestTrackingEvent event.Type = "VoiceRoomMicRequestTrackingEvent"
)

type VoiceRoomMicRequestTrackingEventData struct {
	UserId uint64
	StatAt time.Time
}

func NewVoiceRoomMicRequestTrackingEventData(userId uint64, statAt time.Time) *VoiceRoomMicRequestTrackingEventData {
	return &VoiceRoomMicRequestTrackingEventData{
		UserId: userId,
		StatAt: statAt,
	}
}
