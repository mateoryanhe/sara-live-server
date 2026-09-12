package guilddto

import (
	"github.com/gogf/gf/v2/frame/g"
)

// GetGuildTransferInfoReq CMS 获取工会收款/转账信息
type GetGuildTransferInfoReq struct {
	g.Meta  `path:"/getGuildTransferInfo" method:"post" summary:"获取工会转账信息" tags:"直播工会"`
	GuildId uint64 `json:"guildId" v:"required#工会ID不能为空" dc:"工会ID"`
}

// GuildTransferCountryOption 代付可选国家(统一 country 包)
type GuildTransferCountryOption struct {
	CountryCode string `json:"countryCode" dc:"国家简码"`
	NameEn      string `json:"nameEn"`
	NameZh      string `json:"nameZh"`
	Currency    string `json:"currency" dc:"代付币种,如IDR"`
	Icon        string `json:"icon" dc:"国旗URL"`
}

type GuildTransferInfoItem struct {
	GuildId     string `json:"guildId" dc:"工会ID"`
	CountryCode string `json:"countryCode" dc:"国家简码"`
	Currency    string `json:"currency" dc:"代付币种(由国家推导)"`
	AccountType string `json:"accountType" dc:"HaiPay accountType,默认BANK_ACCOUNT"`
	PayeeName   string `json:"payeeName" dc:"收款人姓名"`
	Phone       string `json:"phone" dc:"收款人手机"`
	Email       string `json:"email" dc:"收款人邮箱"`
	BankName    string `json:"bankName" dc:"银行名称(可选)"`
	AccountNo   string `json:"accountNo" dc:"收款账号"`
	BankCode    string `json:"bankCode" dc:"银行代码"`
	Remark      string `json:"remark" dc:"备注"`
	UpdatedAt   string `json:"updatedAt" dc:"最近更新时间"`
}

type GetGuildTransferInfoRes struct {
	Info      *GuildTransferInfoItem        `json:"info"`
	Countries []*GuildTransferCountryOption `json:"countries" dc:"可选代付国家列表"`
}

// SaveGuildTransferInfoReq CMS 保存工会收款/转账信息(直写DB;选国家推导币种)
type SaveGuildTransferInfoReq struct {
	g.Meta      `path:"/saveGuildTransferInfo" method:"post" summary:"保存工会转账信息" tags:"直播工会"`
	GuildId     uint64 `json:"guildId" v:"required#工会ID不能为空" dc:"工会ID"`
	CountryCode string `json:"countryCode" v:"required|length:2,8#国家不能为空|国家简码无效" dc:"国家/地区ISO简码(如ID)"`
	AccountType string `json:"accountType" v:"max-length:32#accountType最长32" dc:"默认BANK_ACCOUNT"`
	PayeeName   string `json:"payeeName" v:"max-length:128#收款人姓名最长128字符" dc:"收款人姓名"`
	Phone       string `json:"phone" v:"max-length:32#手机号最长32字符" dc:"收款人手机"`
	Email       string `json:"email" v:"max-length:128#邮箱最长128字符" dc:"收款人邮箱"`
	BankName    string `json:"bankName" v:"max-length:128#银行名称最长128字符" dc:"银行名称"`
	AccountNo   string `json:"accountNo" v:"max-length:128#收款账号最长128字符" dc:"收款账号"`
	BankCode    string `json:"bankCode" v:"max-length:64#银行代码最长64字符" dc:"银行代码"`
	Remark      string `json:"remark" v:"max-length:255#备注最长255字符" dc:"备注"`
}

type SaveGuildTransferInfoRes struct {
	Success bool `json:"success"`
}
