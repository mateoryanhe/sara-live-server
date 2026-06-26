package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbLiveGift db.TbName = "live_gifts"
)

// 礼物上下架状态
const (
	LiveGiftStatusOffShelf uint8 = 0 // 下架
	LiveGiftStatusOnShelf  uint8 = 1 // 上架
)

// LiveGift 直播礼物配置(由 CMS 后台管理)
type LiveGift struct {
	migrate.OneModel
	Name        string     `gorm:"size:64;comment:礼物名称" json:"name"`
	NameEn      string     `gorm:"size:64;default:'';comment:礼物名称(英文)" json:"nameEn"`
	NameEs      string     `gorm:"size:64;default:'';comment:礼物名称(西班牙语)" json:"nameEs"`
	NamePt      string     `gorm:"size:64;default:'';comment:礼物名称(葡萄牙语)" json:"namePt"`
	NameHi      string     `gorm:"size:64;default:'';comment:礼物名称(印地语)" json:"nameHi"`
	Icon        string     `gorm:"size:255;default:'';comment:图标URL" json:"icon"`
	Animation   string     `gorm:"size:255;default:'';comment:动画资源URL" json:"animation"`
	Price       float64    `gorm:"default:0;comment:钻石单价" json:"price"`
	Category    string     `gorm:"size:32;default:'';comment:分类" json:"category"`
	Sort        int        `gorm:"default:0;comment:排序值(越大越靠前)" json:"sort"`
	Status      uint8      `gorm:"default:0;comment:状态(0-下架,1-上架)" json:"status"`
	PublishedAt *time.Time `gorm:"comment:发布时间" json:"publishedAt"`
	Description string     `gorm:"size:255;default:'';comment:描述" json:"description"`
}

func InitLiveGift() {
	migrate.AutoMigrate(&LiveGift{})
}
