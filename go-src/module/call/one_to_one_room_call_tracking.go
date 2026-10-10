package call

import (
	"time"

	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

func pubCall1v1RoomTracking(callerId, anchorId uint64, at time.Time) {
	if callerId == 0 || anchorId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	event.Pub(gameevent.Call1v1RoomTrackingEvent, gameevent.NewCall1v1RoomTrackingEventData(callerId, anchorId, at))
}
