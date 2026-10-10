package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// MiniGameRoundResultTrackingEvent 第三方 vendor transfer 单局结算成功(每次 +1)
	MiniGameRoundResultTrackingEvent event.Type = "MiniGameRoundResultTrackingEvent"
)

type MiniGameRoundResultTrackingEventData struct {
	UserId        uint64
	GameCode      string
	TransactionId string
	StatAt        time.Time
}

func NewMiniGameRoundResultTrackingEventData(userId uint64, gameCode, transactionId string, statAt time.Time) *MiniGameRoundResultTrackingEventData {
	return &MiniGameRoundResultTrackingEventData{
		UserId:        userId,
		GameCode:      gameCode,
		TransactionId: transactionId,
		StatAt:        statAt,
	}
}
