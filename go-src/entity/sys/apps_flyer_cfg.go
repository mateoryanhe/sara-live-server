package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbAppsFlyerCfg db.TbName = "apps_flyer_cfgs"
)

// AppsFlyerCfg AppsFlyer 平台密钥配置(CMS 管理,通常仅一条).
type AppsFlyerCfg struct {
	migrate.OneModel
	AndroidDevKey string `gorm:"size:128;default:'';comment:Android Dev Key" json:"androidDevKey"`
	IosDevKey     string `gorm:"size:128;default:'';comment:iOS Dev Key" json:"iosDevKey"`
	IosAppId      string `gorm:"size:128;default:'';comment:iOS App ID" json:"iosAppId"`
}

func (AppsFlyerCfg) TableName() string {
	return string(TbAppsFlyerCfg)
}

func initAppsFlyerCfg() {
	migrate.AutoMigrate(&AppsFlyerCfg{})
}