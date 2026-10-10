package liveroom

import (
	"time"

	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

// pubShowcaseLiveRoomLeaveTracking leaveRoom 成功返回前发布秀场退房埋点(由 tracking 模块异步累计).
func pubShowcaseLiveRoomLeaveTracking(userId, roomId uint64, at time.Time) {
	if userId == 0 || roomId == 0 || userId == roomId {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	event.Pub(gameevent.ShowcaseLiveRoomLeaveTrackingEvent, gameevent.NewShowcaseLiveRoomLeaveTrackingEventData(userId, roomId, at))
}
