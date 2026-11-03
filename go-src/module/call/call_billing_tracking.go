package call

import (
	"context"
	"time"

	"xr-game-server/core/event"
	"xr-game-server/entity/call"
	"xr-game-server/gameevent"
)

func pubCall1v1BillingSuccessTracking(ctx context.Context, order *entity.CallOrder, at time.Time) {
	if order == nil || order.ID == 0 {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	ev := gameevent.NewCall1v1BillingSuccessTrackingEventData(
		order.Source,
		order.CallType,
		order.CallerId,
		order.ReceiverId,
		order.ID,
		at,
	)
	ev.ClientPlatform = gameevent.StatClientPlatformFromCtx(ctx)
	event.Pub(gameevent.Call1v1BillingSuccessTrackingEvent, ev)
}
