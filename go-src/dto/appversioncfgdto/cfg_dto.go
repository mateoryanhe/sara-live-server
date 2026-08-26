package appversioncfgdto

import "github.com/gogf/gf/v2/frame/g"

type GetAppVersionCfgReq struct {
	g.Meta `path:"/getAppVersionCfg" method:"post" summary:"查询App版本配置" tags:"App版本配置"`
}

type AppVersionUpdateDetailItem struct {
	Content string `json:"content"`
	Sort    int    `json:"sort"`
}

type AppVersionCfgItem struct {
	ID                  string                        `json:"id"`
	VersionQueryEnabled bool                          `json:"versionQueryEnabled"`
	Version             string                        `json:"version"`
	BuildVersion        string                        `json:"buildVersion"`
	DownloadUrl         string                        `json:"downloadUrl"`
	DownloadUrlArm      string                        `json:"downloadUrlArm"`
	DownloadUrlAbi      string                        `json:"downloadUrlAbi"`
	UpdateDetails       []*AppVersionUpdateDetailItem `json:"updateDetails"`
	CreatedAt           string                        `json:"createdAt"`
	UpdatedAt           string                        `json:"updatedAt"`
}

type GetAppVersionCfgRes struct {
	Cfg *AppVersionCfgItem `json:"cfg"`
}

type SaveAppVersionCfgReq struct {
	g.Meta              `path:"/saveAppVersionCfg" method:"post" summary:"保存App版本配置" tags:"App版本配置"`
	ID                  uint64                        `json:"id" dc:"配置ID,首次保存可为0"`
	VersionQueryEnabled bool                          `json:"versionQueryEnabled" dc:"App版本查询开关(仅透传App端,服务端不做拦截)"`
	Version             string                        `json:"version" dc:"版本号"`
	BuildVersion        string                        `json:"buildVersion" dc:"构建版本号"`
	DownloadUrl         string                        `json:"downloadUrl" dc:"下载地址"`
	DownloadUrlArm      string                        `json:"downloadUrlArm" dc:"ARM架构下载地址"`
	DownloadUrlAbi      string                        `json:"downloadUrlAbi" dc:"ABI架构下载地址"`
	UpdateDetails       []*AppVersionUpdateDetailItem `json:"updateDetails" dc:"更新明细列表"`
}

type SaveAppVersionCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type AppVersionQueryReq struct {
	g.Meta `path:"/appVersionQuery" method:"post" summary:"App版本查询" tags:"App版本"`
}

type AppVersionQueryRes struct {
	Enabled       bool                          `json:"enabled" dc:"App版本查询开关(仅透传,由App端自行处理)"`
	Version       string                        `json:"version"`
	BuildVersion  string                        `json:"buildVersion"`
	DownloadUrl   string                        `json:"downloadUrl"`
	DownloadUrlArm string                       `json:"downloadUrlArm"`
	DownloadUrlAbi string                       `json:"downloadUrlAbi"`
	UpdateDetails []*AppVersionUpdateDetailItem `json:"updateDetails"`
}
