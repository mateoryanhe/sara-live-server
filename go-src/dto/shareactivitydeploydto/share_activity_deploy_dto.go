package shareactivitydeploydto

import (
	"github.com/gogf/gf/v2/frame/g"
)

const (
	ShareActivitySiteKey           = "share-activity"
	ShareActivityStaticPrefix      = "/share-activity"
	DefaultShareActivityDeployPath = "/home/ec2-user/cdn/share-activity"
)

type GetShareActivityDeployInfoReq struct {
	g.Meta `path:"/getShareActivityDeployInfo" method:"post" summary:"获取分享活动页部署信息" tags:"分享活动页部署"`
}

type ShareActivityDeployInfoItem struct {
	ID           string `json:"id" dc:"配置ID"`
	Domain       string `json:"domain" dc:"站点域名,不含协议与路径"`
	UrlPrefix    string `json:"urlPrefix" dc:"静态访问前缀,如 /share-activity"`
	DeployPath   string `json:"deployPath" dc:"解压目标物理目录"`
	AcceptExt    string `json:"acceptExt" dc:"允许上传的扩展名"`
	UpdatedAt    string `json:"updatedAt" dc:"最近更新时间"`
	LastUploadAt string `json:"lastUploadAt" dc:"最后一次成功上传并解压时间"`
}

type GetShareActivityDeployInfoRes struct {
	Info *ShareActivityDeployInfoItem `json:"info"`
}

type SaveShareActivityDeployCfgReq struct {
	g.Meta     `path:"/saveShareActivityDeployCfg" method:"post" summary:"保存分享活动页部署配置" tags:"分享活动页部署"`
	ID         uint64 `json:"id" dc:"配置ID,首次保存可为0"`
	Domain     string `json:"domain" v:"required#域名不能为空" dc:"站点域名,不含协议与路径"`
	DeployPath string `json:"deployPath" v:"required#部署目录不能为空" dc:"静态文件绝对目录"`
}

type SaveShareActivityDeployCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type DeployShareActivityZipReq struct {
	g.Meta `path:"/deployZip" method:"post" mime:"multipart/form-data" summary:"上传zip并部署分享活动页静态资源" tags:"分享活动页部署"`
}

type DeployShareActivityZipRes struct {
	FileCount  int    `json:"fileCount" dc:"解压写入的文件数"`
	DirCount   int    `json:"dirCount" dc:"创建的目录数"`
	DeployPath string `json:"deployPath" dc:"解压目标物理目录"`
	UrlPrefix  string `json:"urlPrefix" dc:"静态访问前缀"`
}
