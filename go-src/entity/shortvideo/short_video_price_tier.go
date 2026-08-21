package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbShortVideoPriceTier db.TbName = "short_video_price_tiers"
)

const (
	ShortVideoPriceTierStatusOffShelf uint8 = 0
	ShortVideoPriceTierStatusOnShelf  uint8 = 1
)

// ShortVideoPriceTier 短视频价格挡位(CMS 管理)
type ShortVideoPriceTier struct {
	migrate.OneModel
	Price  float64 `gorm:"type:decimal(10,4);default:0;comment:钻石价格" json:"price"`
	Status uint8   `gorm:"default:1;comment:状态(0-下架,1-上架)" json:"status"`
}

func initShortVideoPriceTier() {
	migrate.AutoMigrate(&ShortVideoPriceTier{})
}
