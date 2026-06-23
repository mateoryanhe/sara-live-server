package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbShortVideoCategory db.TbName = "short_video_categories"
)

// ShortVideoCategory 短视频分类(CMS 管理)
type ShortVideoCategory struct {
	migrate.OneModel
	Name string `gorm:"size:64;comment:分类名称" json:"name"`
	Sort int    `gorm:"default:0;comment:排序值(越大越靠前)" json:"sort"`
}

func initShortVideoCategory() {
	migrate.AutoMigrate(&ShortVideoCategory{})
}
