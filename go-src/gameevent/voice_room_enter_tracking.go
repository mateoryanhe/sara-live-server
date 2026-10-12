package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// VoiceRoomEnterTrackingEvent 观众成功 joinRoom 后的语聊房进房次数埋点(每次 +1,由 tracking 模块消费)
	VoiceRoomEnterTrackingEvent event.Type = "VoiceRoomEnterTrackingEvent"
)

// VoiceRoomEnterTrackingEventData joinRoom 成功埋点载荷
type VoiceRoomEnterTrackingEventData struct {
	UserId uint64
	RoomId uint64
	StatAt time.Time
}

func NewVoiceRoomEnterTrackingEventData(userId, roomId uint64, statAt time.Time) *VoiceRoomEnterTrackingEventData {
	return &VoiceRoomEnterTrackingEventData{
		UserId: userId,
		RoomId: roomId,
		StatAt: statAt,
	}
}
