package h5livedeploydto

import (
	"github.com/gogf/gf/v2/frame/g"
)

const H5LiveStaticPrefix = "/h5-live"

// GetH5LiveDeployInfoReq 获取 H5 直播静态部署信息
type GetH5LiveDeployInfoReq struct {
	g.Meta `path:"/getH5LiveDeployInfo" method:"post" summary:"获取H5直播部署信息" tags:"H5直播部署"`
}

type H5LiveDeployInfoItem struct {
	UrlPrefix  string `json:"urlPrefix" dc:"静态访问前缀,如 /h5-live"`
	DeployPath string `json:"deployPath" dc:"解压目标物理目录"`
	AcceptExt  string `json:"acceptExt" dc:"允许上传的扩展名"`
}

type GetH5LiveDeployInfoRes struct {
	Info *H5LiveDeployInfoItem `json:"info"`
}

// DeployH5LiveZipReq CMS 上传 zip 并解压到 H5 直播目录;路由由 RegCMSHandler 注册
type DeployH5LiveZipReq struct {
	g.Meta `path:"/deployZip" method:"post" mime:"multipart/form-data" summary:"上传zip并部署H5直播静态资源" tags:"H5直播部署"`
}

type DeployH5LiveZipRes struct {
	FileCount   int    `json:"fileCount" dc:"解压写入的文件数"`
	DirCount    int    `json:"dirCount" dc:"创建的目录数"`
	DeployPath  string `json:"deployPath" dc:"解压目标物理目录"`
	UrlPrefix   string `json:"urlPrefix" dc:"静态访问前缀"`
}
