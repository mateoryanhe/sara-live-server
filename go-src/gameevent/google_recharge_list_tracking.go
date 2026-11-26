package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const GoogleRechargeListTrackingEvent event.Type = "GoogleRechargeListTrackingEvent"

type GoogleRechargeListTrackingEventData struct {
	UserId uint64
	StatAt time.Time
}

func NewGoogleRechargeListTrackingEventData(userId uint64, statAt time.Time) *GoogleRechargeListTrackingEventData {
	return &GoogleRechargeListTrackingEventData{UserId: userId, StatAt: statAt}
}