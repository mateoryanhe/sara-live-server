package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// LiveFirstFrameRenderedTrackingEvent App 上报直播首帧画面渲染完成(每次 +1,由 tracking 模块消费)
	LiveFirstFrameRenderedTrackingEvent event.Type = "LiveFirstFrameRenderedTrackingEvent"
)

// LiveFirstFrameRenderedTrackingEventData 首帧渲染埋点载荷
type LiveFirstFrameRenderedTrackingEventData struct {
	UserId uint64
	RoomId uint64
	StatAt time.Time
}

func NewLiveFirstFrameRenderedTrackingEventData(userId, roomId uint64, statAt time.Time) *LiveFirstFrameRenderedTrackingEventData {
	return &LiveFirstFrameRenderedTrackingEventData{
		UserId: userId,
		RoomId: roomId,
		StatAt: statAt,
	}
}
