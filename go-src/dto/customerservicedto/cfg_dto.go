package customerservicedto

import "github.com/gogf/gf/v2/frame/g"

type GetCustomerServiceCfgReq struct {
	g.Meta `path:"/getCustomerServiceCfg" method:"post" summary:"查询客服联系配置" tags:"客服联系配置"`
}

type CustomerServiceCfgItem struct {
	ID          string `json:"id"`
	TelegramUrl string `json:"telegramUrl"`
	FacebookUrl string `json:"facebookUrl"`
	WhatsappUrl string `json:"whatsappUrl"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type GetCustomerServiceCfgRes struct {
	Cfg *CustomerServiceCfgItem `json:"cfg"`
}

type SaveCustomerServiceCfgReq struct {
	g.Meta      `path:"/saveCustomerServiceCfg" method:"post" summary:"保存客服联系配置" tags:"客服联系配置"`
	ID          uint64 `json:"id" dc:"配置ID,首次保存可为0"`
	TelegramUrl string `json:"telegramUrl" v:"max-length:512#Telegram联系方式长度不能超过512" dc:"Telegram联系方式"`
	FacebookUrl string `json:"facebookUrl" v:"max-length:512#Facebook联系方式长度不能超过512" dc:"Facebook联系方式"`
	WhatsappUrl string `json:"whatsappUrl" v:"max-length:512#WhatsApp联系方式长度不能超过512" dc:"WhatsApp联系方式"`
}

type SaveCustomerServiceCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type AppCustomerServiceCfgReq struct {
	g.Meta `path:"/cfg" method:"post" summary:"App查询客服联系配置" tags:"客服联系"`
}

type AppCustomerServiceCfgRes struct {
	TelegramUrl string `json:"telegramUrl" dc:"Telegram联系方式"`
	FacebookUrl string `json:"facebookUrl" dc:"Facebook联系方式"`
	WhatsappUrl string `json:"whatsappUrl" dc:"WhatsApp联系方式"`
}
