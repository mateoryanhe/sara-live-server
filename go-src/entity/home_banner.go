package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbHomeBanner db.TbName = "home_banners"
)

const (
	HomeBannerStatusOffShelf uint8 = 0
	HomeBannerStatusOnShelf  uint8 = 1
)

// HomeBanner App 首页 Banner 配置(CMS 管理)
type HomeBanner struct {
	migrate.OneModel
	Title  string `gorm:"size:64;comment:标题" json:"title"`
	Image  string `gorm:"size:255;default:'';comment:图片资源名" json:"image"`
	Link   string `gorm:"size:512;default:'';comment:跳转链接" json:"link"`
	Sort   int    `gorm:"default:0;comment:排序值(越大越靠前)" json:"sort"`
	Status uint8  `gorm:"default:0;comment:状态(0-下架,1-上架)" json:"status"`
}

func initHomeBanner() {
	migrate.AutoMigrate(&HomeBanner{})
}
