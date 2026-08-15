package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbGooglePlayCfg db.TbName = "google_play_cfgs"
)

// GooglePlayCfg Google Play RTDN 充值配置(CMS 管理,通常仅一条)
type GooglePlayCfg struct {
	migrate.OneModel
	Enabled            bool   `gorm:"default:0;comment:是否启用" json:"enabled"`
	PackageName        string `gorm:"size:128;default:'';comment:Android包名" json:"packageName"`
	ServiceAccountJson string `gorm:"type:text;comment:Google服务账号JSON" json:"serviceAccountJson"`
}

func initGooglePlayCfg() {
	migrate.AutoMigrate(&GooglePlayCfg{})
}
