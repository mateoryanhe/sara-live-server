package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const ChargeResultServerTrackingEvent event.Type = "ChargeResultServerTrackingEvent"

type ChargeResultServerTrackingEventData struct {
	UserId uint64
	StatAt time.Time
}

func NewChargeResultServerTrackingEventData(userId uint64, statAt time.Time) *ChargeResultServerTrackingEventData {
	return &ChargeResultServerTrackingEventData{UserId: userId, StatAt: statAt}
}
