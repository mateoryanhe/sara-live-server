package liveroom

import (
	"context"
	"time"

	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

// pubLiveRoomJoinTracking joinRoom 成功返回前发布进房埋点(秀场/游戏/语聊由 tracking 按 category 过滤累计).
func pubLiveRoomJoinTracking(ctx context.Context, userId, roomId uint64, at time.Time) {
	if userId == 0 || roomId == 0 || userId == roomId {
		return
	}
	if at.IsZero() {
		at = time.Now()
	}
	platform := gameevent.StatClientPlatformFromCtx(ctx)
	showcase := gameevent.NewShowcaseLiveRoomJoinTrackingEventData(userId, roomId, at)
	showcase.ClientPlatform = platform
	event.Pub(gameevent.ShowcaseLiveRoomJoinTrackingEvent, showcase)
	if !gameevent.GameLiveRoomTrackingPersistDisabled {
		game := gameevent.NewGameLiveRoomJoinTrackingEventData(userId, roomId, at)
		game.ClientPlatform = platform
		event.Pub(gameevent.GameLiveRoomJoinTrackingEvent, game)
	}
	voice := gameevent.NewVoiceRoomEnterTrackingEventData(userId, roomId, at)
	voice.ClientPlatform = platform
	event.Pub(gameevent.VoiceRoomEnterTrackingEvent, voice)
}
