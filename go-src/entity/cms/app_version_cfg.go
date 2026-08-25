package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbAppVersionCfg db.TbName = "app_version_cfgs"
)

// AppVersionCfg App版本查询配置(CMS 管理,通常仅一条)
type AppVersionCfg struct {
	migrate.OneModel
	VersionQueryEnabled bool   `gorm:"default:0;comment:App版本查询开关(仅透传App端)" json:"versionQueryEnabled"`
	Version             string `gorm:"size:32;default:'';comment:版本号" json:"version"`
	BuildVersion        string `gorm:"size:32;default:'';comment:构建版本号" json:"buildVersion"`
	DownloadUrl         string `gorm:"size:512;default:'';comment:下载地址" json:"downloadUrl"`
}

func initAppVersionCfg() {
	migrate.AutoMigrate(&AppVersionCfg{})
}
