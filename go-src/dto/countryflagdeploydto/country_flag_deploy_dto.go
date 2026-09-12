package countryflagdeploydto

import "github.com/gogf/gf/v2/frame/g"

// CountryFlagStaticDir 图片静态根下的国旗资源目录名
const CountryFlagStaticDir = "country-flags"

// GetCountryFlagDeployInfoReq 获取国旗部署信息
type GetCountryFlagDeployInfoReq struct {
	g.Meta `path:"/getCountryFlagDeployInfo" method:"post" summary:"获取国旗资源部署信息" tags:"国旗资源部署"`
}

type CountryFlagDeployInfoItem struct {
	ID         string `json:"id" dc:"配置ID"`
	Version    string `json:"version" dc:"当前生效版本目录名"`
	UrlPrefix  string `json:"urlPrefix" dc:"当前版本URL前缀,如 /images/country-flags/20260102150405"`
	DeployPath string `json:"deployPath" dc:"国旗资源根物理目录(其下按 version 分子目录)"`
	AcceptExt  string `json:"acceptExt" dc:"允许上传扩展名"`
	UpdatedAt  string `json:"updatedAt" dc:"最近更新时间"`
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
