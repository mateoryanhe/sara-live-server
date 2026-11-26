package recharge

import (
	"time"

	"xr-game-server/core/event"
	"xr-game-server/entity/recharge"
	"xr-game-server/gameevent"
)

func pubGoogleRechargeListTracking(userId uint64, at time.Time) {
	if at.IsZero() {
		at = time.Now()
	}
	event.Pub(gameevent.GoogleRechargeListTrackingEvent, gameevent.NewGoogleRechargeListTrackingEventData(userId, at))
}

func pubGoogleRechargeOrderTracking(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	event.Pub(gameevent.GoogleRechargeOrderTrackingEvent, gameevent.NewGoogleRechargeOrderTrackingEventData(userId, at))
}

func pubGoogleRechargeResultTracking(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	event.Pub(gameevent.GoogleRechargeResultTrackingEvent, gameevent.NewGoogleRechargeResultTrackingEventData(userId, at))
}

func isGooglePlayPayChannel(ch uint8) bool {
	return ch == entity.RechargeCfgTypeGoogle
}