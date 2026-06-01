package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbShortVideo db.TbName = "short_videos"
)

const (
	ShortVideoStatusOffShelf uint8 = 0
	ShortVideoStatusOnShelf  uint8 = 1
)

const (
	ShortVideoPaidNo  uint8 = 0 // 免费
	ShortVideoPaidYes uint8 = 1 // 付费
)

// ShortVideo 短视频(CMS 管理)
type ShortVideo struct {
	migrate.OneModel
	Title            string `gorm:"size:64;comment:标题" json:"title"`
	Video            string `gorm:"size:255;default:'';comment:视频资源名" json:"video"`
	Cover            string `gorm:"size:255;default:'';comment:封面资源名" json:"cover"`
	Sort             int    `gorm:"default:0;comment:排序值(越大越靠前)" json:"sort"`
	Status           uint8  `gorm:"default:0;comment:状态(0-下架,1-上架)" json:"status"`
	IsPaid           uint8  `gorm:"default:0;comment:是否付费(0免费,1付费)" json:"isPaid"`
	DiamondPerSecond uint64 `gorm:"default:0;comment:每秒钻石数(付费时有效)" json:"diamondPerSecond"`
	Description      string `gorm:"size:255;default:'';comment:描述" json:"description"`
}

func initShortVideo() {
	migrate.AutoMigrate(&ShortVideo{})
}
