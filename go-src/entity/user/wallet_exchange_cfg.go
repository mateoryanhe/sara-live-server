package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbWalletExchangeCfg db.TbName = "wallet_exchange_cfgs"
)

const (
	WalletExchangeCfgGoldToDiamondRate  db.TbCol = "gold_to_diamond_rate"
	WalletExchangeCfgExchangeFeePercent db.TbCol = "exchange_fee_percent"
)

// WalletExchangeCfg 金币兑换钻石配置(CMS 管理,通常仅一条)
type WalletExchangeCfg struct {
	migrate.OneModel
	GoldToDiamondRate  int     `gorm:"default:100;comment:1金币兑换钻石数" json:"goldToDiamondRate"`
	ExchangeFeePercent float64 `gorm:"type:decimal(6,2);default:3;comment:App手动兑换手续费(%)，从兑换钻石中扣除" json:"exchangeFeePercent"`
}

func initWalletExchangeCfg() {
	migrate.AutoMigrate(&WalletExchangeCfg{})
}
