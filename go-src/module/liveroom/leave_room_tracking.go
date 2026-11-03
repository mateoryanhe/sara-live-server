package liveroom

import (
	"context"
	"time"

	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

// pubShowcaseLiveRoomLeaveTracking leaveRoom 成功返回前发布秀场退房埋点(由 tracking 模块异步累计).
func pubShowcaseLiveRoomLeaveTracking(ctx context.Context, userId, roomId uint64, at time.Time) {
	if userId == 0 || roomId == 0 || userId == roomId {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	ev := gameevent.NewShowcaseLiveRoomLeaveTrackingEventData(userId, roomId, at)
	ev.ClientPlatform = gameevent.StatClientPlatformFromCtx(ctx)
	event.Pub(gameevent.ShowcaseLiveRoomLeaveTrackingEvent, ev)
}
