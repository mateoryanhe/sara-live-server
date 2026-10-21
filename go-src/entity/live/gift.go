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

// 礼物分类(新字段 category_type; 旧字段 category 字符串保留兼容)
const (
	LiveGiftCategoryTypeNormal uint8 = 0 // 普通
	LiveGiftCategoryTypeLucky  uint8 = 1 // 幸运礼物
)

func IsValidLiveGiftCategoryType(v uint8) bool {
	return v == LiveGiftCategoryTypeNormal || v == LiveGiftCategoryTypeLucky
}

// LiveGift 直播礼物配置(由 CMS 后台管理)
type LiveGift struct {
	migrate.OneModel
	Name        string     `gorm:"size:64;comment:礼物名称" json:"name"`
	NameEn      string     `gorm:"size:64;default:'';comment:礼物名称(英文)" json:"nameEn"`
	NameEs      string     `gorm:"size:64;default:'';comment:礼物名称(西班牙语)" json:"nameEs"`
	NamePt      string     `gorm:"size:64;default:'';comment:礼物名称(葡萄牙语)" json:"namePt"`
	NameHi      string     `gorm:"size:64;default:'';comment:礼物名称(印地语)" json:"nameHi"`
	NameId      string     `gorm:"size:64;default:'';comment:礼物名称(印尼语)" json:"nameId"`
	Icon        string     `gorm:"size:255;default:'';comment:图标URL" json:"icon"`
	Animation   string     `gorm:"size:255;default:'';comment:动画资源URL" json:"animation"`
	Price       float64    `gorm:"default:0;comment:钻石单价" json:"price"`
	Category     string `gorm:"size:32;default:'';comment:分类(历史字段,不再维护)" json:"category"`
	CategoryType uint8  `gorm:"default:0;comment:礼物分类(0普通,1幸运礼物)" json:"categoryType"`
	Sort         int    `gorm:"default:0;comment:排序值(越大越靠前)" json:"sort"`
	Status      uint8      `gorm:"default:0;comment:状态(0-下架,1-上架)" json:"status"`
	PublishedAt *time.Time `gorm:"comment:发布时间" json:"publishedAt"`
	Description  string  `gorm:"size:255;default:'';comment:描述" json:"description"`
	LuckyGiftRtp float64 `gorm:"type:decimal(10,4);default:0;comment:幸运礼物RTP(百分比,0-100,仅category_type=1)" json:"luckyGiftRtp"`
}

func InitLiveGift() {
	migrate.AutoMigrate(&LiveGift{})
	initLuckyGiftMultiplierTier()
}
