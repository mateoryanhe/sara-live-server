package game

import (
	"time"

	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

func pubMiniGameRoundStartTracking(userId uint64, gameCode string, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	event.Pub(gameevent.MiniGameRoundStartTrackingEvent, gameevent.NewMiniGameRoundStartTrackingEventData(userId, gameCode, at))
}
