package uploaddto

import "github.com/gogf/gf/v2/frame/g"

type GetUploadResourceCfgReq struct {
	g.Meta `path:"/getUploadResourceCfg" method:"post" summary:"查询上传资源配置" tags:"上传配置"`
}

type UploadResourceCfgItem struct {
	ID               string `json:"id"`
	ResourceDomain   string `json:"resourceDomain"`
	DefaultAvatarUrl string `json:"defaultAvatarUrl"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

type GetUploadResourceCfgRes struct {
	Cfg *UploadResourceCfgItem `json:"cfg"`
}

type SaveUploadResourceCfgReq struct {
	g.Meta           `path:"/saveUploadResourceCfg" method:"post" summary:"保存上传资源配置" tags:"上传配置"`
	ID               uint64 `json:"id"`
	ResourceDomain   string `json:"resourceDomain"`
	DefaultAvatarUrl string `json:"defaultAvatarUrl"`
}

type SaveUploadResourceCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
