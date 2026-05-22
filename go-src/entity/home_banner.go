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

// Banner 展示位置
const (
	HomeBannerPositionHomeTop    uint8 = 1 // 首页顶部
	HomeBannerPositionHomeMiddle uint8 = 2 // 首页中部
	HomeBannerPositionHomeBottom uint8 = 3 // 首页底部
	HomeBannerPositionLiveHall   uint8 = 4 // 直播大厅
)

// HomeBanner App 首页 Banner 配置(CMS 管理)
type HomeBanner struct {
	migrate.OneModel
	Title     string `gorm:"size:64;comment:标题" json:"title"`
	Image     string `gorm:"size:255;default:'';comment:图片资源名" json:"image"`
	Link      string `gorm:"size:512;default:'';comment:跳转链接" json:"link"`
	Direction uint8  `gorm:"default:1;comment:展示位置(1首页顶部,2首页中部,3首页底部,4直播大厅)" json:"direction"`
	Sort      int    `gorm:"default:0;comment:排序值(越大越靠前)" json:"sort"`
	Status    uint8  `gorm:"default:0;comment:状态(0-下架,1-上架)" json:"status"`
}

func initHomeBanner() {
	migrate.AutoMigrate(&HomeBanner{})
}
