package gameevent

import "xr-game-server/core/event"

const (
	// LiveRecordStoppedEvent 单场直播结束事件
	LiveRecordStoppedEvent event.Type = "LiveRecordStoppedEvent"
)

// LiveRecordStoppedEventData 单场直播结束事件数据
type LiveRecordStoppedEventData struct {
	LiveRecordId uint64
}

func NewLiveRecordStoppedEventData(liveRecordId uint64) *LiveRecordStoppedEventData {
	return &LiveRecordStoppedEventData{LiveRecordId: liveRecordId}
}
