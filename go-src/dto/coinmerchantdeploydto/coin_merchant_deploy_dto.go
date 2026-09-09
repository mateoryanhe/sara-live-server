package coinmerchantdeploydto

import (
	"github.com/gogf/gf/v2/frame/g"
)

const CoinMerchantStaticPrefix = "/coin-merchant"

// GetCoinMerchantDeployInfoReq 获取币商 H5 静态部署信息
type GetCoinMerchantDeployInfoReq struct {
	g.Meta `path:"/getCoinMerchantDeployInfo" method:"post" summary:"获取币商H5部署信息" tags:"币商H5部署"`
}

type CoinMerchantDeployInfoItem struct {
	ID           string `json:"id" dc:"配置ID"`
	UrlPrefix    string `json:"urlPrefix" dc:"静态访问前缀,如 /coin-merchant"`
	DeployPath   string `json:"deployPath" dc:"解压目标物理目录"`
	AcceptExt    string `json:"acceptExt" dc:"允许上传的扩展名"`
	DeploySecret string `json:"deploySecret" dc:"币商H5部署密钥"`
	UpdatedAt    string `json:"updatedAt" dc:"最近更新时间"`
}

type GetCoinMerchantDeployInfoRes struct {
	Info *CoinMerchantDeployInfoItem `json:"info"`
}

type SaveCoinMerchantDeployCfgReq struct {
	g.Meta       `path:"/saveCoinMerchantDeployCfg" method:"post" summary:"保存币商H5部署配置" tags:"币商H5部署"`
	ID           uint64 `json:"id" dc:"配置ID,首次保存可为0"`
	DeploySecret string `json:"deploySecret" v:"required#部署密钥不能为空" dc:"币商H5部署密钥"`
}

type SaveCoinMerchantDeployCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

// DeployCoinMerchantZipReq CMS 上传 zip 并解压到币商 H5 目录;路由由 RegCMSHandler 注册
type DeployCoinMerchantZipReq struct {
	g.Meta `path:"/deployZip" method:"post" mime:"multipart/form-data" summary:"上传zip并部署币商H5静态资源" tags:"币商H5部署"`
}

type DeployCoinMerchantZipRes struct {
	FileCount  int    `json:"fileCount" dc:"解压写入的文件数"`
	DirCount   int    `json:"dirCount" dc:"创建的目录数"`
	DeployPath string `json:"deployPath" dc:"解压目标物理目录"`
	UrlPrefix  string `json:"urlPrefix" dc:"静态访问前缀"`
}
