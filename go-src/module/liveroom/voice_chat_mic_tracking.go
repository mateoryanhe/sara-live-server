package liveroom

import (
	"time"

	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

func pubVoiceRoomMicRequestTracking(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	event.Pub(gameevent.VoiceRoomMicRequestTrackingEvent, gameevent.NewVoiceRoomMicRequestTrackingEventData(userId, at))
}
