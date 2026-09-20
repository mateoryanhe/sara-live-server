package guilddto

import (
	"github.com/gogf/gf/v2/frame/g"
)

// GetGuildTransferInfoReq CMS 获取工会收款/转账信息
type GetGuildTransferInfoReq struct {
	g.Meta  `path:"/getGuildTransferInfo" method:"post" summary:"获取工会转账信息" tags:"直播工会"`
	GuildId uint64 `json:"guildId" v:"required#工会ID不能为空" dc:"工会ID"`
}

// GuildTransferWalletOption HaiPay 代付钱包选项。
type GuildTransferWalletOption struct {
	Code string `json:"code" dc:"HaiPay钱包支付编码"`
	Name string `json:"name" dc:"钱包名称"`
}

// GuildTransferMethodOption 是一条可直接用于 HaiPay 代付的支付方式。
type GuildTransferMethodOption struct {
	AccountType          string   `json:"accountType" dc:"HaiPay账户类型"`
	BankCode             string   `json:"bankCode" dc:"HaiPay支付编码"`
	Limit                string   `json:"limit" dc:"官网标注的单笔限额"`
	Description          string   `json:"description" dc:"官网支付方式说明"`
	IdentifyTypeRequired bool     `json:"identifyTypeRequired" dc:"是否必须填写identifyType"`
	IdentifyTypeOptions  []string `json:"identifyTypeOptions" dc:"identifyType可选值;空数组表示自由输入"`
	CountryRequired      bool     `json:"countryRequired" dc:"是否必须上报收款国家"`
	AddressRequired      bool     `json:"addressRequired" dc:"是否必须填写收款地址"`
}

// GuildTransferCountryOption 代付可选国家(统一 country 包)。
type GuildTransferCountryOption struct {
	CountryCode  string                      `json:"countryCode" dc:"国家简码"`
	NameEn       string                      `json:"nameEn"`
	NameZh       string                      `json:"nameZh"`
	Currency     string                      `json:"currency" dc:"代付币种,如IDR"`
	Region       string                      `json:"region" dc:"HaiPay代付地区分类"`
	AccountTypes []string                    `json:"accountTypes" dc:"该币种支持的HaiPay账户类型"`
	Wallets      []GuildTransferWalletOption `json:"wallets" dc:"当前已接入的电子钱包"`
	Methods      []GuildTransferMethodOption `json:"methods" dc:"当前可用的具体代付方式"`
	Icon         string                      `json:"icon" dc:"国旗URL"`
}

type GuildTransferInfoItem struct {
	GuildId      string `json:"guildId" dc:"工会ID"`
	CountryCode  string `json:"countryCode" dc:"国家简码"`
	Currency     string `json:"currency" dc:"代付币种(由国家推导)"`
	AccountType  string `json:"accountType" dc:"HaiPay accountType,默认BANK_ACCOUNT"`
	PayeeName    string `json:"payeeName" dc:"收款人姓名"`
	Phone        string `json:"phone" dc:"收款人手机"`
	Email        string `json:"email" dc:"收款人邮箱"`
	BankName     string `json:"bankName" dc:"HaiPay支付方式说明(由支付编码自动填充)"`
	AccountNo    string `json:"accountNo" dc:"收款账号"`
	BankCode     string `json:"bankCode" dc:"银行代码"`
	IdentifyType string `json:"identifyType" dc:"特殊代付方式的账号类型或银行路由编码"`
	Address1     string `json:"address1" dc:"收款地址-详细街道"`
	Address2     string `json:"address2" dc:"收款地址-城市"`
	Address3     string `json:"address3" dc:"收款地址-省/州"`
	PostalCode   string `json:"postalCode" dc:"收款邮编"`
	Remark       string `json:"remark" dc:"备注"`
	UpdatedAt    string `json:"updatedAt" dc:"最近更新时间"`
}

type GetGuildTransferInfoRes struct {
	Info      *GuildTransferInfoItem        `json:"info"`
	Countries []*GuildTransferCountryOption `json:"countries" dc:"可选代付国家列表"`
}

// SaveGuildTransferInfoReq CMS 保存工会收款/转账信息(直写DB;选国家推导币种)
type SaveGuildTransferInfoReq struct {
	g.Meta       `path:"/saveGuildTransferInfo" method:"post" summary:"保存工会转账信息" tags:"直播工会"`
	GuildId      uint64 `json:"guildId" v:"required#工会ID不能为空" dc:"工会ID"`
	CountryCode  string `json:"countryCode" v:"required|length:2,8#国家不能为空|国家简码无效" dc:"国家/地区ISO简码(如ID)"`
	Currency     string `json:"currency" v:"max-length:16#币种最长16字符" dc:"代付币种(优先使用,兼容旧客户端可留空)"`
	AccountType  string `json:"accountType" v:"required|max-length:32#accountType不能为空|accountType最长32" dc:"HaiPay账户类型"`
	PayeeName    string `json:"payeeName" v:"max-length:128#收款人姓名最长128字符" dc:"收款人姓名"`
	Phone        string `json:"phone" v:"max-length:32#手机号最长32字符" dc:"收款人手机"`
	Email        string `json:"email" v:"max-length:128#邮箱最长128字符" dc:"收款人邮箱"`
	BankName     string `json:"bankName" v:"max-length:128#支付方式说明最长128字符" dc:"兼容字段,服务端按支付编码自动填充"`
	AccountNo    string `json:"accountNo" v:"max-length:128#收款账号最长128字符" dc:"收款账号"`
	BankCode     string `json:"bankCode" v:"required|max-length:64#支付编码不能为空|支付编码最长64字符" dc:"与accountType配套的HaiPay支付编码"`
	IdentifyType string `json:"identifyType" v:"max-length:64#identifyType最长64字符" dc:"PIX/Papara账号类型或ACH银行路由编码"`
	Address1     string `json:"address1" v:"max-length:255#详细街道最长255字符" dc:"收款地址-详细街道"`
	Address2     string `json:"address2" v:"max-length:128#城市最长128字符" dc:"收款地址-城市"`
	Address3     string `json:"address3" v:"max-length:128#省州最长128字符" dc:"收款地址-省/州"`
	PostalCode   string `json:"postalCode" v:"max-length:32#邮编最长32字符" dc:"收款邮编"`
	Remark       string `json:"remark" v:"max-length:255#备注最长255字符" dc:"备注"`
}

type SaveGuildTransferInfoRes struct {
	Success bool `json:"success"`
}
