package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbShareAppDownloadCfg db.TbName = "share_app_download_cfgs"
)

// ShareAppDownloadCfg 分享 App 下载链接文案(CMS 管理,通常仅一条)
type ShareAppDownloadCfg struct {
	migrate.OneModel
	DownloadUrl string `gorm:"size:512;default:'';comment:App下载链接" json:"downloadUrl"`
	ShareText   string `gorm:"type:text;comment:分享文案" json:"shareText"`
}

func initShareAppDownloadCfg() {
	migrate.AutoMigrate(&ShareAppDownloadCfg{})
}
