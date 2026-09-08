package uploaddto

import "github.com/gogf/gf/v2/frame/g"

type GetUploadResourceCfgReq struct {
	g.Meta `path:"/getUploadResourceCfg" method:"post" summary:"查询上传资源配置" tags:"上传配置"`
}

type UploadResourceCfgItem struct {
	ID                             string `json:"id"`
	ResourceDomain                 string `json:"resourceDomain"`
	StoragePath                    string `json:"storagePath" dc:"统一文件存储路径"`
	CmsExportTtlMinutes            int    `json:"cmsExportTtlMinutes" dc:"CMS文件导出过期清理(分钟)"`
	AppImageMaxSizeMB              uint32 `json:"appImageMaxSizeMB" dc:"App端图片上传大小上限(MB)"`
	ImageModerationEnabled         bool   `json:"imageModerationEnabled"`
	ImageModerationAccessKeyId     string `json:"imageModerationAccessKeyId"`
	ImageModerationAccessKeySecret string `json:"imageModerationAccessKeySecret"`
	ImageModerationRegionId        string `json:"imageModerationRegionId"`
	ImageModerationEndpoint        string `json:"imageModerationEndpoint"`
	ImageModerationService         string `json:"imageModerationService"`
	S3Enabled                      bool   `json:"s3Enabled"`
	S3PublicDomain                 string `json:"s3PublicDomain"`
	S3Endpoint                     string `json:"s3Endpoint"`
	S3Bucket                       string `json:"s3Bucket"`
	S3AccessKeyId                  string `json:"s3AccessKeyId"`
	S3SecretAccessKey              string `json:"s3SecretAccessKey"`
	S3KeyPrefix                    string `json:"s3KeyPrefix"`
	CreatedAt                      string `json:"createdAt"`
	UpdatedAt                      string `json:"updatedAt"`
}

type GetUploadResourceCfgRes struct {
	Cfg *UploadResourceCfgItem `json:"cfg"`
}

type SaveUploadResourceCfgReq struct {
	g.Meta                         `path:"/saveUploadResourceCfg" method:"post" summary:"保存上传资源配置" tags:"上传配置"`
	ID                             uint64 `json:"id" p:"id"`
	ResourceDomain                 string `json:"resourceDomain" p:"resourceDomain"`
	StoragePath                    string `json:"storagePath" p:"storagePath" v:"required|length:1,512#存储路径不能为空|存储路径过长" dc:"统一文件存储路径"`
	CmsExportTtlMinutes            int    `json:"cmsExportTtlMinutes" p:"cmsExportTtlMinutes" v:"min:0|max:10080#TTL不能为负|TTL不能超过10080分钟" dc:"CMS文件导出过期清理(分钟),0表示使用默认30分钟"`
	AppImageMaxSizeMB              uint32 `json:"appImageMaxSizeMB" p:"appImageMaxSizeMB" dc:"App端图片上传大小上限(MB)"`
	ImageModerationEnabled         bool   `json:"imageModerationEnabled" p:"imageModerationEnabled"`
	ImageModerationAccessKeyId     string `json:"imageModerationAccessKeyId" p:"imageModerationAccessKeyId"`
	ImageModerationAccessKeySecret string `json:"imageModerationAccessKeySecret" p:"imageModerationAccessKeySecret"`
	ImageModerationRegionId        string `json:"imageModerationRegionId" p:"imageModerationRegionId"`
	ImageModerationEndpoint        string `json:"imageModerationEndpoint" p:"imageModerationEndpoint"`
	ImageModerationService         string `json:"imageModerationService" p:"imageModerationService"`
	S3Enabled                      bool   `json:"s3Enabled" p:"s3Enabled"`
	S3PublicDomain                 string `json:"s3PublicDomain" p:"s3PublicDomain"`
	S3Endpoint                     string `json:"s3Endpoint" p:"s3Endpoint"`
	S3Bucket                       string `json:"s3Bucket" p:"s3Bucket"`
	S3AccessKeyId                  string `json:"s3AccessKeyId" p:"s3AccessKeyId"`
	S3SecretAccessKey              string `json:"s3SecretAccessKey" p:"s3SecretAccessKey"`
	S3KeyPrefix                    string `json:"s3KeyPrefix" p:"s3KeyPrefix"`
}

type SaveUploadResourceCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
