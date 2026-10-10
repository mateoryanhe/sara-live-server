package liveroom

import (
	"time"

	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

// pubLiveRoomJoinTracking joinRoom 成功返回前发布进房埋点(秀场/游戏类由 tracking 按 category 过滤累计).
func pubLiveRoomJoinTracking(userId, roomId uint64, at time.Time) {
	if userId == 0 || roomId == 0 || userId == roomId {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	event.Pub(gameevent.ShowcaseLiveRoomJoinTrackingEvent, gameevent.NewShowcaseLiveRoomJoinTrackingEventData(userId, roomId, at))
	event.Pub(gameevent.GameLiveRoomJoinTrackingEvent, gameevent.NewGameLiveRoomJoinTrackingEventData(userId, roomId, at))
}
