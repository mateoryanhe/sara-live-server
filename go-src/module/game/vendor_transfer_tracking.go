package game

import (
	"time"

	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

func pubMiniGameRoundResultTracking(userId uint64, gameCode, transactionId string, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	event.Pub(
		gameevent.MiniGameRoundResultTrackingEvent,
		gameevent.NewMiniGameRoundResultTrackingEventData(userId, gameCode, transactionId, at),
	)
}

func statTimeFromVendorUpdatedTime(updatedTime int64) time.Time {
	if updatedTime <= 0 {
		return time.Now()
	}
	sec := updatedTime
	if updatedTime > 1_000_000_000_000 {
		sec = updatedTime / 1000
	}
	return time.Unix(sec, 0)
}
