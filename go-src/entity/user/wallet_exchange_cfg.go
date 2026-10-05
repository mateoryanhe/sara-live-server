package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbWalletExchangeCfg db.TbName = "wallet_exchange_cfgs"
)

const (
	WalletExchangeCfgGoldToDiamondRate                  db.TbCol = "gold_to_diamond_rate"
	WalletExchangeCfgExchangeFeePercent                 db.TbCol = "exchange_fee_percent"
	WalletExchangeCfgUsdToGoldRate                      db.TbCol = "usd_to_gold_rate"
	WalletExchangeCfgEffectiveLiveMinSessionMinutes     db.TbCol = "effective_live_min_session_minutes"
	WalletExchangeCfgPlatformAnchorMinimumSettlementUsd db.TbCol = "platform_anchor_minimum_settlement_usd"
)

// WalletExchangeCfg 结算基础配置(CMS 管理,通常仅一条)。
// 钱包兑换、有效直播门槛和平台主播最低结算金额共用这一行。
type WalletExchangeCfg struct {
	migrate.OneModel
	GoldToDiamondRate                  int     `gorm:"default:100;comment:1金币兑换钻石数" json:"goldToDiamondRate"`
	ExchangeFeePercent                 float64 `gorm:"type:decimal(6,2);default:3;comment:App手动兑换手续费(%)，从兑换钻石中扣除" json:"exchangeFeePercent"`
	UsdToGoldRate                      int     `gorm:"default:100;comment:1美金兑换金币数" json:"usdToGoldRate"`
	EffectiveLiveMinSessionMinutes     int     `gorm:"default:30;comment:单场直播计入有效时长的门槛(分钟)" json:"effectiveLiveMinSessionMinutes"`
	PlatformAnchorMinimumSettlementUsd float64 `gorm:"type:decimal(16,2);default:5;comment:平台主播生成结算单必须严格超过的最低金额(USD)" json:"platformAnchorMinimumSettlementUsd"`
}

func initWalletExchangeCfg() {
	migrate.AutoMigrate(&WalletExchangeCfg{})
}
