package uploaddto

import "github.com/gogf/gf/v2/frame/g"

type GetUploadResourceCfgReq struct {
	g.Meta `path:"/getUploadResourceCfg" method:"post" summary:"查询上传资源配置" tags:"上传配置"`
}

type UploadResourceCfgItem struct {
	ID                             string `json:"id"`
	ResourceDomain                 string `json:"resourceDomain"`
	DefaultAvatarUrl               string `json:"defaultAvatarUrl"`
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
	DefaultAvatarUrl               string `json:"defaultAvatarUrl"`
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
