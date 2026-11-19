package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const EmailAccountTrackingEvent event.Type = "EmailAccountTrackingEvent"

type EmailAccountTrackingEventData struct {
	StatAt time.Time
}

func NewEmailAccountTrackingEventData(statAt time.Time) *EmailAccountTrackingEventData {
	if statAt.IsZero() {
		statAt = time.Now()
	}
	return &EmailAccountTrackingEventData{StatAt: statAt}
}

// PubEmailAccountTracking user_ext.email 对该 user_id 首次写入时 +1(异步埋点; 注销同邮箱再注册=new user_id 再计)
func PubEmailAccountTracking(statAt time.Time) {
	event.Pub(EmailAccountTrackingEvent, NewEmailAccountTrackingEventData(statAt))
}