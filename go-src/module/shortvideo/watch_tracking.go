package shortvideo

import (
	"time"

	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

func pubShortVideoPlayStartTracking(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	event.Pub(gameevent.ShortVideoPlayStartTrackingEvent, gameevent.NewShortVideoPlayStartTrackingEventData(userId, at))
}

func pubShortVideoPaywallClickTracking(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	event.Pub(gameevent.ShortVideoPaywallClickTrackingEvent, gameevent.NewShortVideoPaywallClickTrackingEventData(userId, at))
}

func pubShortVideoPaySuccessTracking(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	event.Pub(gameevent.ShortVideoPaySuccessTrackingEvent, gameevent.NewShortVideoPaySuccessTrackingEventData(userId, at))
}

func pubShortVideoWatchEndTracking(userId uint64, at time.Time) {
	if userId == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	event.Pub(gameevent.ShortVideoWatchEndTrackingEvent, gameevent.NewShortVideoWatchEndTrackingEventData(userId, at))
}
