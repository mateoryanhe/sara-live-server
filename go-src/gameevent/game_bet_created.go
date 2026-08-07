package gameevent

import "xr-game-server/core/event"

const (
	// GameBetCreatedEvent 游戏下注记录创建事件
	GameBetCreatedEvent event.Type = "GameBetCreatedEvent"
)

// GameBetCreatedEventData 游戏下注事件数据
type GameBetCreatedEventData struct {
	UserId       uint64
	LiveRecordId uint64
	Amount       float64
}

func NewGameBetCreatedEventData(userId, liveRecordId uint64, amount float64) *GameBetCreatedEventData {
	return &GameBetCreatedEventData{
		UserId:       userId,
		LiveRecordId: liveRecordId,
		Amount:       amount,
	}
}
