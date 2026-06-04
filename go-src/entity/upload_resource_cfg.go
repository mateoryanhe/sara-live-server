package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbUploadResourceCfg db.TbName = "upload_resource_cfgs"
)

// UploadResourceCfg 静态资源域名与默认头像(CMS 管理,通常仅一条)
type UploadResourceCfg struct {
	migrate.OneModel
	ResourceDomain   string `gorm:"size:256;default:'';comment:资源访问域名" json:"resourceDomain"`
	DefaultAvatarUrl string `gorm:"size:512;default:'';comment:默认头像完整URL" json:"defaultAvatarUrl"`
}

func initUploadResourceCfg() {
	migrate.AutoMigrate(&UploadResourceCfg{})
}
