package countryflagdeploydto

import "github.com/gogf/gf/v2/frame/g"

// CountryFlagStaticDir 图片静态根下的国旗资源目录名
const CountryFlagStaticDir = "country-flags"

// GetCountryFlagDeployInfoReq 获取国旗部署信息
type GetCountryFlagDeployInfoReq struct {
	g.Meta `path:"/getCountryFlagDeployInfo" method:"post" summary:"获取国旗资源部署信息" tags:"国旗资源部署"`
}

type CountryFlagDeployInfoItem struct {
	ID         string                    `json:"id" dc:"配置ID"`
	Version    string                    `json:"version" dc:"当前生效版本目录名"`
	UrlPrefix  string                    `json:"urlPrefix" dc:"当前版本完整URL前缀(与头像同一资源域名/S3域)"`
	DeployPath string                    `json:"deployPath" dc:"国旗资源根物理目录(其下按 version 分子目录)"`
	AcceptExt  string                    `json:"acceptExt" dc:"允许上传扩展名"`
	UpdatedAt  string                    `json:"updatedAt" dc:"最近更新时间"`
	FileCount  int                       `json:"fileCount" dc:"当前版本已部署国旗数量"`
	Flags      []*CountryFlagPreviewItem `json:"flags" dc:"当前版本国旗列表(便于CMS查询)"`
}

// CountryFlagPreviewItem 当前版本单面国旗预览
type CountryFlagPreviewItem struct {
	Code   string `json:"code" dc:"国家简码"`
	NameEn string `json:"nameEn" dc:"英文名"`
	NameZh string `json:"nameZh" dc:"中文名"`
	File   string `json:"file" dc:"文件名,如 id.png"`
	Icon   string `json:"icon" dc:"完整URL"`
}

type GetCountryFlagDeployInfoRes struct {
	Info *CountryFlagDeployInfoItem `json:"info"`
}

// DeployCountryFlagZipReq CMS 上传 zip;路由由 RegCMSHandler 注册
type DeployCountryFlagZipReq struct {
	g.Meta `path:"/deployZip" method:"post" mime:"multipart/form-data" summary:"上传zip部署国旗资源(新version)" tags:"国旗资源部署"`
}

type DeployCountryFlagZipRes struct {
	Version    string `json:"version" dc:"新生成的版本号"`
	FileCount  int    `json:"fileCount" dc:"解压写入的png数"`
	DeployPath string `json:"deployPath" dc:"本版本物理目录"`
	UrlPrefix  string `json:"urlPrefix" dc:"本版本URL前缀"`
	Removed    int    `json:"removed" dc:"清理的旧version目录数"`
}
