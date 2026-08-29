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
	CreatedAt                      string `json:"createdAt"`
	UpdatedAt                      string `json:"updatedAt"`
}

type GetUploadResourceCfgRes struct {
	Cfg *UploadResourceCfgItem `json:"cfg"`
}

type SaveUploadResourceCfgReq struct {
	g.Meta                         `path:"/saveUploadResourceCfg" method:"post" summary:"保存上传资源配置" tags:"上传配置"`
	ID                             uint64 `json:"id"`
	ResourceDomain                 string `json:"resourceDomain"`
	StoragePath                    string `json:"storagePath" v:"required|length:1,512#存储路径不能为空|存储路径过长" dc:"统一文件存储路径"`
	CmsExportTtlMinutes            int    `json:"cmsExportTtlMinutes" v:"min:0|max:10080#TTL不能为负|TTL不能超过10080分钟" dc:"CMS文件导出过期清理(分钟),0表示使用默认30分钟"`
	AppImageMaxSizeMB              uint32 `json:"appImageMaxSizeMB" dc:"App端图片上传大小上限(MB)"`
	ImageModerationEnabled         bool   `json:"imageModerationEnabled"`
	ImageModerationAccessKeyId     string `json:"imageModerationAccessKeyId"`
	ImageModerationAccessKeySecret string `json:"imageModerationAccessKeySecret"`
	ImageModerationRegionId        string `json:"imageModerationRegionId"`
	ImageModerationEndpoint        string `json:"imageModerationEndpoint"`
	ImageModerationService         string `json:"imageModerationService"`
}

type SaveUploadResourceCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
