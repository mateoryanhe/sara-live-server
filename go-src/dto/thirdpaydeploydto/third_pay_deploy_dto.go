package thirdpaydeploydto

import (
	"github.com/gogf/gf/v2/frame/g"
)

const ThirdPayStaticPrefix = "/third-pay"

// GetThirdPayDeployInfoReq 获取第三方支付静态部署信息
type GetThirdPayDeployInfoReq struct {
	g.Meta `path:"/getThirdPayDeployInfo" method:"post" summary:"获取第三方支付部署信息" tags:"第三方支付部署"`
}

type ThirdPayDeployInfoItem struct {
	UrlPrefix  string `json:"urlPrefix" dc:"静态访问前缀,如 /third-pay"`
	DeployPath string `json:"deployPath" dc:"解压目标物理目录"`
	AcceptExt  string `json:"acceptExt" dc:"允许上传的扩展名"`
}

type GetThirdPayDeployInfoRes struct {
	Info *ThirdPayDeployInfoItem `json:"info"`
}

// DeployThirdPayZipReq CMS 上传 zip 并解压到第三方支付目录;路由由 RegCMSHandler 注册
type DeployThirdPayZipReq struct {
	g.Meta `path:"/deployZip" method:"post" mime:"multipart/form-data" summary:"上传zip并部署第三方支付静态资源" tags:"第三方支付部署"`
}

type DeployThirdPayZipRes struct {
	FileCount  int    `json:"fileCount" dc:"解压写入的文件数"`
	DirCount   int    `json:"dirCount" dc:"创建的目录数"`
	DeployPath string `json:"deployPath" dc:"解压目标物理目录"`
	UrlPrefix  string `json:"urlPrefix" dc:"静态访问前缀"`
}
