package liveroom

import (
	"time"

	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

// pubShowcaseLiveRoomJoinTracking joinRoom 成功返回前发布秀场进房埋点(由 tracking 模块异步累计).
func pubShowcaseLiveRoomJoinTracking(userId, roomId uint64, at time.Time) {
	if userId == 0 || roomId == 0 || userId == roomId {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	event.Pub(gameevent.ShowcaseLiveRoomJoinTrackingEvent, gameevent.NewShowcaseLiveRoomJoinTrackingEventData(userId, roomId, at))
}
