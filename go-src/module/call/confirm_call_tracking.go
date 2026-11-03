package call

import (
	"context"
	"time"

	"xr-game-server/core/event"
	"xr-game-server/entity/call"
	"xr-game-server/gameevent"
)

func pubCall1v1ConnectSuccessTracking(ctx context.Context, order *entity.CallOrder, at time.Time) {
	if order == nil {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	ev := gameevent.NewCall1v1ConnectSuccessTrackingEventData(
		order.Source,
		order.CallType,
		order.CallerId,
		order.ReceiverId,
		at,
	)
	ev.ClientPlatform = gameevent.StatClientPlatformFromCtx(ctx)
	event.Pub(gameevent.Call1v1ConnectSuccessTrackingEvent, ev)
}
