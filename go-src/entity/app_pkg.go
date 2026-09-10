package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbAppPkg db.TbName = "app_pkgs"
)

// AppPkg App包配置(CMS管理)
type AppPkg struct {
	migrate.OneModel
	PackageName          string `gorm:"uniqueIndex;size:128;comment:包名" json:"packageName"`
	Remark               string `gorm:"size:512;default:'';comment:备注" json:"remark"`
	AttributionEnabled   bool   `gorm:"default:0;comment:是否启用归因" json:"attributionEnabled"`
	AttributionProvider  string `gorm:"size:64;default:'';comment:归因渠道(如 appsFlyer)" json:"attributionProvider"`
	AppsFlyerDevKey      string `gorm:"size:128;default:'';comment:AppsFlyer Dev Key" json:"appsFlyerDevKey"`
	AppsFlyerAppId       string `gorm:"size:128;default:'';comment:AppsFlyer App ID(iOS)" json:"appsFlyerAppId"`
}

func initAppPkg() {
	migrate.AutoMigrate(&AppPkg{})
}
