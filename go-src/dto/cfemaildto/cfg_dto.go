package cfemaildto

import "github.com/gogf/gf/v2/frame/g"

type GetCfEmailCfgReq struct {
	g.Meta `path:"/getCfEmailCfg" method:"post" summary:"查询邮件SMTP发信配置" tags:"邮件发信配置"`
}

type CfEmailCfgItem struct {
	ID                    string `json:"id"`
	Enabled               bool   `json:"enabled"`
	SmtpHost              string `json:"smtpHost"`
	SmtpPort              int    `json:"smtpPort"`
	SmtpUsername          string `json:"smtpUsername"`
	SmtpPassword          string `json:"smtpPassword"`
	SmtpPasswordConfigured bool   `json:"smtpPasswordConfigured"`
	FromEmail             string `json:"fromEmail"`
	CreatedAt             string `json:"createdAt"`
	UpdatedAt             string `json:"updatedAt"`
}

type GetCfEmailCfgRes struct {
	Cfg *CfEmailCfgItem `json:"cfg"`
}

type SaveCfEmailCfgReq struct {
	g.Meta       `path:"/saveCfEmailCfg" method:"post" summary:"保存邮件SMTP发信配置" tags:"邮件发信配置"`
	ID           uint64 `json:"id" dc:"配置ID,新建传0"`
	Enabled      bool   `json:"enabled" dc:"是否启用"`
	SmtpHost     string `json:"smtpHost" v:"required|length:1,256#SMTP主机不能为空|SMTP主机长度需在1到256之间" dc:"SMTP服务器地址"`
	SmtpPort     int    `json:"smtpPort" v:"required|between:1,65535#SMTP端口无效" dc:"SMTP端口(SES通常为587)"`
	SmtpUsername string `json:"smtpUsername" v:"required|length:1,128" dc:"SMTP用户名(可填明文或Base64)"`
	SmtpPassword string `json:"smtpPassword" dc:"SMTP密码(新建必填;修改时空表示不变,可填明文或Base64)"`
	FromEmail    string `json:"fromEmail" v:"required|length:1,256" dc:"发件地址(须已在SES验证)"`
}

type SaveCfEmailCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type SendCfEmailTestReq struct {
	g.Meta    `path:"/sendCfEmailTest" method:"post" summary:"发送测试验证码邮件" tags:"邮件发信配置"`
	TestEmail string `json:"testEmail" v:"required|email#测试邮箱不能为空|测试邮箱格式不正确" dc:"测试收件邮箱"`
	Lang      string `json:"lang" dc:"邮件语言(en/es/hi/pt/id/zh);空则 en"`
}

type SendCfEmailTestRes struct {
	Success bool `json:"success"`
}
