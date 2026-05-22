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

// 充值配置类型(支付渠道)
const (
	RechargeCfgTypeIOS     uint8 = 1 // iOS
	RechargeCfgTypeGoogle  uint8 = 2 // Google
	RechargeCfgTypeChannel uint8 = 3 // 渠道
)

const RechargeCfgCurrencyUSD = "USD"

// RechargeCfg 充值档位配置(由 CMS 后台管理,无缓存,直接读写 DB)
type RechargeCfg struct {
	migrate.OneModel
	Name         string `gorm:"size:64;comment:档位名称" json:"name"`
	CfgType      uint8  `gorm:"column:cfg_type;default:1;comment:类型(1iOS,2Google,3渠道)" json:"cfgType"`
	Icon         string `gorm:"size:255;default:'';comment:图标URL" json:"icon"`
	Diamond      uint64 `gorm:"default:0;comment:基础到账金币数" json:"diamond"`
	ExtraDiamond uint64 `gorm:"default:0;comment:额外赠送金币数" json:"extraDiamond"`
	Price        uint64 `gorm:"default:0;comment:现实货币价格(单位:美分,USD)" json:"price"`
	Currency     string `gorm:"size:8;default:'USD';comment:货币(固定USD)" json:"currency"`
	ProductId    string `gorm:"size:64;default:'';comment:第三方商品SKU" json:"productId"`
	Sort         int    `gorm:"default:0;comment:排序值(越大越靠前)" json:"sort"`
	Status       uint8  `gorm:"default:0;comment:状态(0-下架,1-上架)" json:"status"`
	Description  string `gorm:"size:255;default:'';comment:描述" json:"description"`
}

func initRechargeCfg() {
	migrate.AutoMigrate(&RechargeCfg{})
}
