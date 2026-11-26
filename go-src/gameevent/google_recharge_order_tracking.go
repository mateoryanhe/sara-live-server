package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const GoogleRechargeOrderTrackingEvent event.Type = "GoogleRechargeOrderTrackingEvent"

type GoogleRechargeOrderTrackingEventData struct {
	UserId uint64
	StatAt time.Time
}

func NewGoogleRechargeOrderTrackingEventData(userId uint64, statAt time.Time) *GoogleRechargeOrderTrackingEventData {
	return &GoogleRechargeOrderTrackingEventData{UserId: userId, StatAt: statAt}
}