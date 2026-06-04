package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbUploadResourceCfg db.TbName = "upload_resource_cfgs"
)

// UploadResourceCfg 静态资源域名、默认头像与 App 图片审核(CMS 管理,通常仅一条)
type UploadResourceCfg struct {
	migrate.OneModel
	ResourceDomain   string `gorm:"size:256;default:'';comment:资源访问域名" json:"resourceDomain"`
	DefaultAvatarUrl string `gorm:"size:512;default:'';comment:默认头像完整URL" json:"defaultAvatarUrl"`
	// App 端图片审核(阿里云 ImageModeration)
	ImageModerationEnabled         bool   `gorm:"default:0;comment:是否开启App图片审核" json:"imageModerationEnabled"`
	ImageModerationAccessKeyId     string `gorm:"size:128;default:'';comment:图片审核AccessKeyId" json:"imageModerationAccessKeyId"`
	ImageModerationAccessKeySecret string `gorm:"size:256;default:'';comment:图片审核AccessKeySecret" json:"imageModerationAccessKeySecret"`
	ImageModerationRegionId        string `gorm:"size:32;default:'cn-shanghai';comment:图片审核地域" json:"imageModerationRegionId"`
	ImageModerationEndpoint        string `gorm:"size:128;default:'green-cip.cn-shanghai.aliyuncs.com';comment:图片审核接入点" json:"imageModerationEndpoint"`
	ImageModerationService         string `gorm:"size:64;default:'profilePhotoCheck';comment:图片审核Service" json:"imageModerationService"`
}

func initUploadResourceCfg() {
	migrate.AutoMigrate(&UploadResourceCfg{})
}
