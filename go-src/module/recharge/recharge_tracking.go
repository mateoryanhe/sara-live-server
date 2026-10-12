package recharge

import (
	"time"

	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

func pubChargeClickTracking(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	event.Pub(gameevent.ChargeClickTrackingEvent, gameevent.NewChargeClickTrackingEventData(userId, at))
}

func pubChargeResultServerTracking(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	event.Pub(gameevent.ChargeResultServerTrackingEvent, gameevent.NewChargeResultServerTrackingEventData(userId, at))
}
