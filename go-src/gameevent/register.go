package gameevent

import (
	"time"
	"xr-game-server/core/event"
)

const (
	// RegisterEvent 用户注册成功事件
	RegisterEvent event.Type = "RegisterEvent"
)

// RegisterEventData 用户注册成功事件数据
type RegisterEventData struct {
	UserId       uint64
	RegisteredAt time.Time
}

func NewRegisterEventData(userId uint64, registeredAt time.Time) *RegisterEventData {
	return &RegisterEventData{
		UserId:       userId,
		RegisteredAt: registeredAt,
	}
}
