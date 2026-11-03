package call

import (
	"context"
	"time"

	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

func pubCall1v1InitiateTracking(ctx context.Context, callerId, anchorId uint64, at time.Time) {
	if callerId == 0 || anchorId == 0 || callerId == anchorId {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	ev := gameevent.NewCall1v1InitiateTrackingEventData(callerId, anchorId, at)
	ev.ClientPlatform = gameevent.StatClientPlatformFromCtx(ctx)
	event.Pub(gameevent.Call1v1InitiateTrackingEvent, ev)
}
