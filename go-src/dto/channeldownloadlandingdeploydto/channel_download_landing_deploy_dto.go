package channeldownloadlandingdeploydto

import (
	"github.com/gogf/gf/v2/frame/g"
)

const (
	ChannelDownloadLandingSiteKey      = "channel-download-landing"
	ChannelDownloadLandingStaticPrefix = "/channel-download"
	DefaultChannelDownloadDeployPath   = "/home/ec2-user/cdn/channel-download-landing"
)

type GetChannelDownloadLandingDeployInfoReq struct {
	g.Meta `path:"/getChannelDownloadLandingDeployInfo" method:"post" summary:"获取渠道下载落地页部署信息" tags:"渠道下载落地页部署"`
}

type ChannelDownloadLandingDeployInfoItem struct {
	ID           string `json:"id" dc:"配置ID"`
	Domain       string `json:"domain" dc:"站点域名,不含协议与路径"`
	UrlPrefix    string `json:"urlPrefix" dc:"静态访问前缀,如 /channel-download"`
	DeployPath   string `json:"deployPath" dc:"解压目标物理目录"`
	AcceptExt    string `json:"acceptExt" dc:"允许上传的扩展名"`
	UpdatedAt    string `json:"updatedAt" dc:"最近更新时间"`
	LastUploadAt string `json:"lastUploadAt" dc:"最后一次成功上传并解压时间"`
}

type GetChannelDownloadLandingDeployInfoRes struct {
	Info *ChannelDownloadLandingDeployInfoItem `json:"info"`
}

type SaveChannelDownloadLandingDeployCfgReq struct {
	g.Meta       `path:"/saveChannelDownloadLandingDeployCfg" method:"post" summary:"保存渠道下载落地页部署配置" tags:"渠道下载落地页部署"`
	ID           uint64 `json:"id" dc:"配置ID,首次保存可为0"`
	Domain       string `json:"domain" v:"required#域名不能为空" dc:"站点域名,不含协议与路径"`
	DeployPath   string `json:"deployPath" v:"required#部署目录不能为空" dc:"静态文件绝对目录"`
}

type SaveChannelDownloadLandingDeployCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type DeployChannelDownloadLandingZipReq struct {
	g.Meta `path:"/deployZip" method:"post" mime:"multipart/form-data" summary:"上传zip并部署渠道下载落地页静态资源" tags:"渠道下载落地页部署"`
}

type DeployChannelDownloadLandingZipRes struct {
	FileCount  int    `json:"fileCount" dc:"解压写入的文件数"`
	DirCount   int    `json:"dirCount" dc:"创建的目录数"`
	DeployPath string `json:"deployPath" dc:"解压目标物理目录"`
	UrlPrefix  string `json:"urlPrefix" dc:"静态访问前缀"`
}
