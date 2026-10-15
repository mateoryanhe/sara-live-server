package userwithdrawdeploydto

import (
	"github.com/gogf/gf/v2/frame/g"
)

const (
	UserWithdrawSiteKey           = "user-withdraw"
	UserWithdrawStaticPrefix      = "/user-withdraw"
	DefaultUserWithdrawDeployPath = "/home/ec2-user/cdn/user-withdraw"
)

// GetUserWithdrawDeployInfoReq 获取提现 H5 静态部署信息
type GetUserWithdrawDeployInfoReq struct {
	g.Meta `path:"/getUserWithdrawDeployInfo" method:"post" summary:"获取提现部署信息" tags:"提现部署"`
}

type UserWithdrawDeployInfoItem struct {
	Domain       string `json:"domain" dc:"提现站点域名,不含协议与路径"`
	UrlPrefix    string `json:"urlPrefix" dc:"静态访问前缀,如 /user-withdraw"`
	DeployPath   string `json:"deployPath" dc:"解压目标物理目录"`
	AcceptExt    string `json:"acceptExt" dc:"允许上传的扩展名"`
	LastUploadAt string `json:"lastUploadAt" dc:"最后一次成功上传并解压时间"`
}

type GetUserWithdrawDeployInfoRes struct {
	Info *UserWithdrawDeployInfoItem `json:"info"`
}

type SaveUserWithdrawDeployCfgReq struct {
	g.Meta     `path:"/saveUserWithdrawDeployCfg" method:"post" summary:"保存提现部署配置" tags:"提现部署"`
	Domain     string `json:"domain" v:"required#域名不能为空" dc:"提现站点域名,不含协议与路径"`
	DeployPath string `json:"deployPath" v:"required#部署目录不能为空" dc:"提现静态文件绝对目录"`
}

type SaveUserWithdrawDeployCfgRes struct {
	Success bool `json:"success"`
}

// DeployUserWithdrawZipReq CMS 上传 zip 并解压; 路由由 RegCMSHandler 注册
type DeployUserWithdrawZipReq struct {
	g.Meta `path:"/deployZip" method:"post" mime:"multipart/form-data" summary:"上传zip并部署提现静态资源" tags:"提现部署"`
}

type DeployUserWithdrawZipRes struct {
	FileCount  int    `json:"fileCount" dc:"解压写入的文件数"`
	DirCount   int    `json:"dirCount" dc:"创建的目录数"`
	DeployPath string `json:"deployPath" dc:"解压目标物理目录"`
	UrlPrefix  string `json:"urlPrefix" dc:"静态访问前缀"`
}
