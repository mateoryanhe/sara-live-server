package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const (
	// GuildPayoutSucceededEvent 工会代付成功事件，金额统一使用 USD。
	GuildPayoutSucceededEvent event.Type = "GuildPayoutSucceededEvent"
)

// GuildPayoutSucceededEventData 工会代付成功事件数据。
type GuildPayoutSucceededEventData struct {
	SettlementID uint64
	GuildID      uint64
	UsdAmount    float64
	PaidAt       time.Time
}

func NewGuildPayoutSucceededEventData(settlementID, guildID uint64, usdAmount float64, paidAt time.Time) *GuildPayoutSucceededEventData {
	return &GuildPayoutSucceededEventData{
		SettlementID: settlementID,
		GuildID:      guildID,
		UsdAmount:    usdAmount,
		PaidAt:       paidAt,
	}
}
