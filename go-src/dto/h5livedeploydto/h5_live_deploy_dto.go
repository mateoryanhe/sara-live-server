package h5livedeploydto

import (
	"github.com/gogf/gf/v2/frame/g"
)

const (
	H5LiveSiteKey       = "h5-live"
	H5LiveStaticPrefix  = "/h5-live"
	DefaultH5DeployPath = "/home/ec2-user/cdn/h5-live"
)

// GetH5LiveDeployInfoReq 获取 H5 直播静态部署信息
type GetH5LiveDeployInfoReq struct {
	g.Meta `path:"/getH5LiveDeployInfo" method:"post" summary:"获取H5直播部署信息" tags:"H5直播部署"`
}

type H5LiveDeployInfoItem struct {
	ID           string `json:"id" dc:"配置ID"`
	Domain       string `json:"domain" dc:"H5直播站点域名,不含协议与路径"`
	UrlPrefix    string `json:"urlPrefix" dc:"静态访问前缀,如 /h5-live"`
	DeployPath   string `json:"deployPath" dc:"解压目标物理目录"`
	AcceptExt    string `json:"acceptExt" dc:"允许上传的扩展名"`
	DeploySecret string `json:"deploySecret" dc:"H5部署密钥"`
	UpdatedAt    string `json:"updatedAt" dc:"最近更新时间"`
	LastUploadAt string `json:"lastUploadAt" dc:"最后一次成功上传并解压时间"`
}

type GetH5LiveDeployInfoRes struct {
	Info *H5LiveDeployInfoItem `json:"info"`
}

type SaveH5LiveDeployCfgReq struct {
	g.Meta       `path:"/saveH5LiveDeployCfg" method:"post" summary:"保存H5直播部署配置" tags:"H5直播部署"`
	ID           uint64 `json:"id" dc:"配置ID,首次保存可为0"`
	Domain       string `json:"domain" v:"required#域名不能为空" dc:"H5直播站点域名,不含协议与路径"`
	DeployPath   string `json:"deployPath" v:"required#部署目录不能为空" dc:"H5直播静态文件绝对目录"`
	DeploySecret string `json:"deploySecret" v:"required#部署密钥不能为空" dc:"H5部署密钥"`
}

type SaveH5LiveDeployCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

// DeployH5LiveZipReq CMS 上传 zip 并解压到 H5 直播目录;路由由 RegCMSHandler 注册
type DeployH5LiveZipReq struct {
	g.Meta `path:"/deployZip" method:"post" mime:"multipart/form-data" summary:"上传zip并部署H5直播静态资源" tags:"H5直播部署"`
}

type DeployH5LiveZipRes struct {
	FileCount  int    `json:"fileCount" dc:"解压写入的文件数"`
	DirCount   int    `json:"dirCount" dc:"创建的目录数"`
	DeployPath string `json:"deployPath" dc:"解压目标物理目录"`
	UrlPrefix  string `json:"urlPrefix" dc:"静态访问前缀"`
}
