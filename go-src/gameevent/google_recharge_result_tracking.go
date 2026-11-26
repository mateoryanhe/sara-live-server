package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const GoogleRechargeResultTrackingEvent event.Type = "GoogleRechargeResultTrackingEvent"

type GoogleRechargeResultTrackingEventData struct {
	UserId uint64
	StatAt time.Time
}

func NewGoogleRechargeResultTrackingEventData(userId uint64, statAt time.Time) *GoogleRechargeResultTrackingEventData {
	return &GoogleRechargeResultTrackingEventData{UserId: userId, StatAt: statAt}
}