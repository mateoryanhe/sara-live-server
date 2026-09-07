package cfemaildto

import "github.com/gogf/gf/v2/frame/g"

type GetCfEmailCfgReq struct {
	g.Meta `path:"/getCfEmailCfg" method:"post" summary:"查询Cloudflare邮件发信配置" tags:"Cloudflare邮件配置"`
}

type CfEmailCfgItem struct {
	ID        string `json:"id"`
	Enabled   bool   `json:"enabled"`
	AccountId string `json:"accountId"`
	ApiToken  string `json:"apiToken"`
	FromEmail string `json:"fromEmail"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type GetCfEmailCfgRes struct {
	Cfg *CfEmailCfgItem `json:"cfg"`
}

type SaveCfEmailCfgReq struct {
	g.Meta    `path:"/saveCfEmailCfg" method:"post" summary:"保存Cloudflare邮件发信配置" tags:"Cloudflare邮件配置"`
	ID        uint64 `json:"id" dc:"配置ID,新建传0"`
	Enabled   bool   `json:"enabled" dc:"是否启用"`
	AccountId string `json:"accountId" v:"required|length:1,64#Account ID不能为空|Account ID长度需在1到64之间" dc:"Cloudflare Account ID"`
	ApiToken  string `json:"apiToken" v:"required|length:1,256#API Token不能为空|API Token长度需在1到256之间" dc:"API Token(Email Sending Edit)"`
	FromEmail string `json:"fromEmail" v:"required|length:1,256#发件地址不能为空|发件地址长度需在1到256之间" dc:"发件地址(如noreply@mail.saralive.net)"`
}

type SaveCfEmailCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
