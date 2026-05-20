package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbRechargeCfg db.TbName = "recharge_cfgs"
)

// 充值配置上下架状态
const (
	RechargeCfgStatusOffShelf uint8 = 0 // 下架
	RechargeCfgStatusOnShelf  uint8 = 1 // 上架
)

// RechargeCfg 充值档位配置(由 CMS 后台管理,无缓存,直接读写 DB)
type RechargeCfg struct {
	migrate.OneModel
	Name         string `gorm:"size:64;comment:档位名称" json:"name"`
	Icon         string `gorm:"size:255;default:'';comment:图标URL" json:"icon"`
	Diamond      uint64 `gorm:"default:0;comment:基础到账钻石数" json:"diamond"`
	ExtraDiamond uint64 `gorm:"default:0;comment:额外赠送钻石数" json:"extraDiamond"`
	Price        uint64 `gorm:"default:0;comment:现实货币价格(单位:分)" json:"price"`
	Currency     string `gorm:"size:8;default:'CNY';comment:货币(CNY/USD等)" json:"currency"`
	ProductId    string `gorm:"size:64;default:'';comment:第三方商品SKU" json:"productId"`
	Sort         int    `gorm:"default:0;comment:排序值(越大越靠前)" json:"sort"`
	Status       uint8  `gorm:"default:0;comment:状态(0-下架,1-上架)" json:"status"`
	Description  string `gorm:"size:255;default:'';comment:描述" json:"description"`
}

func initRechargeCfg() {
	migrate.AutoMigrate(&RechargeCfg{})
}
