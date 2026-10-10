package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// MiniGameRoundStartTrackingEvent App 点击开始游戏并成功获取启动链接(每次 +1)
	MiniGameRoundStartTrackingEvent event.Type = "MiniGameRoundStartTrackingEvent"
)

type MiniGameRoundStartTrackingEventData struct {
	UserId   uint64
	GameCode string
	StatAt   time.Time
}

func NewMiniGameRoundStartTrackingEventData(userId uint64, gameCode string, statAt time.Time) *MiniGameRoundStartTrackingEventData {
	return &MiniGameRoundStartTrackingEventData{
		UserId:   userId,
		GameCode: gameCode,
		StatAt:   statAt,
	}
}
