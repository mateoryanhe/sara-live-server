package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbCoinMerchantDeployCfg db.TbName = "coin_merchant_deploy_cfgs"
)

// CoinMerchantDeployCfg 币商 H5 部署配置(CMS 管理,通常仅一条)
type CoinMerchantDeployCfg struct {
	migrate.OneModel
	DeploySecret string `gorm:"size:128;default:'';comment:币商H5部署密钥" json:"deploySecret"`
}

func initCoinMerchantDeployCfg() {
	migrate.AutoMigrate(&CoinMerchantDeployCfg{})
}
