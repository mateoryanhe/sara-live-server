package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbSofieShareAppDownloadCfg db.TbName = "sofie_share_app_download_cfgs"
)

// SofieShareAppDownloadCfg Sofie 分享下载(Google Play; CMS 通常仅一条).
// PackageName 与 X-App-Package 对齐,短链路径 /dl/p/{packageName}.
type SofieShareAppDownloadCfg struct {
	migrate.OneModel
	PackageName   string `gorm:"size:128;default:'';index;comment:适用App包名(空=默认)" json:"packageName"`
	GooglePlayUrl string `gorm:"size:512;default:'';comment:Google Play 商店链接" json:"googlePlayUrl"`
	ShareText     string `gorm:"type:text;comment:分享文案" json:"shareText"`
}

func initSofieShareAppDownloadCfg() {
	migrate.AutoMigrate(&SofieShareAppDownloadCfg{})
}