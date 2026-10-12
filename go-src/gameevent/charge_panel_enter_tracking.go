package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const ChargePanelEnterTrackingEvent event.Type = "ChargePanelEnterTrackingEvent"

type ChargePanelEnterTrackingEventData struct {
	UserId uint64
	StatAt time.Time
}

func NewChargePanelEnterTrackingEventData(userId uint64, statAt time.Time) *ChargePanelEnterTrackingEventData {
	return &ChargePanelEnterTrackingEventData{UserId: userId, StatAt: statAt}
}
