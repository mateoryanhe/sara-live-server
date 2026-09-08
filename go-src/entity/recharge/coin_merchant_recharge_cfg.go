package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbCoinMerchantRechargeCfg db.TbName = "coin_merchant_recharge_cfgs"
)

// 币商充值档位上下架状态(与 recharge_cfgs 一致)
const (
	CoinMerchantRechargeCfgStatusOffShelf uint8 = 0
	CoinMerchantRechargeCfgStatusOnShelf  uint8 = 1
)

// CoinMerchantRechargeCfg 币商充值档位配置(CMS 管理,上架项缓存)
type CoinMerchantRechargeCfg struct {
	migrate.OneModel
	Name   string  `gorm:"size:64;uniqueIndex;comment:档位名称" json:"name"`
	Price  float64 `gorm:"type:decimal(10,4);default:0;comment:USD价格" json:"price"`
	Gold   uint64  `gorm:"default:0;comment:到账金币数" json:"gold"`
	Status uint8   `gorm:"default:0;comment:状态(0下架,1上架)" json:"status"`
}

func initCoinMerchantRechargeCfg() {
	migrate.AutoMigrate(&CoinMerchantRechargeCfg{})
}
