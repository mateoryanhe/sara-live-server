package cfemaildto

import "github.com/gogf/gf/v2/frame/g"

type GetCfEmailCfgReq struct {
	g.Meta `path:"/getCfEmailCfg" method:"post" summary:"查询AWS SES邮件发信配置" tags:"邮件发信配置"`
}

type CfEmailCfgItem struct {
	ID              string `json:"id"`
	Enabled         bool   `json:"enabled"`
	Region          string `json:"region"`
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	FromEmail       string `json:"fromEmail"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

type GetCfEmailCfgRes struct {
	Cfg *CfEmailCfgItem `json:"cfg"`
}

type SaveCfEmailCfgReq struct {
	g.Meta          `path:"/saveCfEmailCfg" method:"post" summary:"保存AWS SES邮件发信配置" tags:"邮件发信配置"`
	ID              uint64 `json:"id" dc:"配置ID,新建传0"`
	Enabled         bool   `json:"enabled" dc:"是否启用"`
	Region          string `json:"region" v:"required|length:1,64#Region不能为空|Region长度需在1到64之间" dc:"AWS Region(如us-west-1)"`
	AccessKeyId     string `json:"accessKeyId" v:"required|length:1,128#Access Key不能为空|Access Key长度需在1到128之间" dc:"AWS Access Key ID"`
	SecretAccessKey string `json:"secretAccessKey" v:"required|length:1,256#Secret Key不能为空|Secret Key长度需在1到256之间" dc:"AWS Secret Access Key"`
	FromEmail       string `json:"fromEmail" v:"required|length:1,256#发件地址不能为空|发件地址长度需在1到256之间" dc:"发件地址(须已在SES验证)"`
}

type SaveCfEmailCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
