package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbUploadResourceCfg db.TbName = "upload_resource_cfgs"
)

// UploadResourceCfg 静态资源域名与 App 图片审核(CMS 管理,通常仅一条)
type UploadResourceCfg struct {
	migrate.OneModel
	ResourceDomain  string `gorm:"size:256;default:'';comment:资源访问域名" json:"resourceDomain"`
	StoragePath     string `gorm:"size:512;default:'/home/ec2-user/cdn/images';comment:统一文件存储路径(头像/图片/短视频/CMS导出)" json:"storagePath"`
	CmsExportTtlMinutes int `gorm:"default:30;comment:CMS文件导出过期清理(分钟)" json:"cmsExportTtlMinutes"`
	AppImageMaxSize uint64 `gorm:"default:1048576;comment:App端图片上传大小上限(字节)" json:"appImageMaxSize"`
	// App 端图片审核(阿里云 ImageModeration)
	ImageModerationEnabled         bool   `gorm:"default:0;comment:是否开启App图片审核" json:"imageModerationEnabled"`
	ImageModerationAccessKeyId     string `gorm:"size:128;default:'';comment:图片审核AccessKeyId" json:"imageModerationAccessKeyId"`
	ImageModerationAccessKeySecret string `gorm:"size:256;default:'';comment:图片审核AccessKeySecret" json:"imageModerationAccessKeySecret"`
	ImageModerationRegionId        string `gorm:"size:32;default:'cn-shanghai';comment:图片审核地域" json:"imageModerationRegionId"`
	ImageModerationEndpoint        string `gorm:"size:128;default:'green-cip.cn-shanghai.aliyuncs.com';comment:图片审核接入点" json:"imageModerationEndpoint"`
	ImageModerationService         string `gorm:"size:64;default:'profilePhotoCheck';comment:图片审核Service" json:"imageModerationService"`
	// S3 兼容对象存储(Cloudflare R2 / AWS S3 等);开启后新上传与 CMS 导出走 PutObject
	S3Enabled         bool   `gorm:"default:0;comment:是否开启S3兼容云桶" json:"s3Enabled"`
	S3PublicDomain    string `gorm:"size:256;default:'';comment:云桶公开访问域名(R2自定义域)" json:"s3PublicDomain"`
	S3Endpoint        string `gorm:"size:256;default:'';comment:S3 API Endpoint(R2必填)" json:"s3Endpoint"`
	S3Region          string `gorm:"size:64;default:'auto';comment:内部固定auto,CMS不再配置" json:"s3Region"`
	S3Bucket          string `gorm:"size:128;default:'';comment:桶名" json:"s3Bucket"`
	S3AccessKeyId     string `gorm:"size:128;default:'';comment:S3 AccessKeyId" json:"s3AccessKeyId"`
	S3SecretAccessKey string `gorm:"size:256;default:'';comment:S3 SecretAccessKey" json:"s3SecretAccessKey"`
	S3KeyPrefix       string `gorm:"size:128;default:'';comment:对象Key环境前缀如test/或prod/" json:"s3KeyPrefix"`
}

func initUploadResourceCfg() {
	migrate.AutoMigrate(&UploadResourceCfg{})
}
