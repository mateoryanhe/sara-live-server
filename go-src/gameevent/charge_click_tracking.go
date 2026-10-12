package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const ChargeClickTrackingEvent event.Type = "ChargeClickTrackingEvent"

type ChargeClickTrackingEventData struct {
	UserId uint64
	StatAt time.Time
}

func NewChargeClickTrackingEventData(userId uint64, statAt time.Time) *ChargeClickTrackingEventData {
	return &ChargeClickTrackingEventData{UserId: userId, StatAt: statAt}
}
