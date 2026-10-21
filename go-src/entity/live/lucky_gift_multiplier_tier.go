package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbLuckyGiftMultiplierTier db.TbName = "lucky_gift_multiplier_tiers"
)

const (
	LuckyGiftMultiplierTypeLow    = "low"
	LuckyGiftMultiplierTypeMedium = "medium"
	LuckyGiftMultiplierTypeHigh   = "high"
)

func IsValidLuckyGiftMultiplierType(v string) bool {
	return v == LuckyGiftMultiplierTypeLow || v == LuckyGiftMultiplierTypeMedium || v == LuckyGiftMultiplierTypeHigh
}

// LuckyGiftMultiplierTier 幸运礼物返奖倍率挡位(按礼物 gift_id 全量覆盖)
type LuckyGiftMultiplierTier struct {
	migrate.OneModel
	GiftID     uint64 `gorm:"index;default:0;comment:礼物ID(live_gifts.id)" json:"giftId"`
	Multiplier uint32 `gorm:"default:1;comment:返奖倍率" json:"multiplier"`
	Weight     uint32 `gorm:"default:1;comment:抽奖权重" json:"weight"`
	Type       string `gorm:"size:16;default:'low';comment:挡位类型(low/medium/high)" json:"type"`
	Sort       int    `gorm:"default:0;comment:排序(越小越靠前)" json:"sort"`
}

func initLuckyGiftMultiplierTier() {
	migrate.AutoMigrate(&LuckyGiftMultiplierTier{})
}
