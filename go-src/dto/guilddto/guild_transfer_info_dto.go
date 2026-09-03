package guilddto

import (
	"github.com/gogf/gf/v2/frame/g"
)

// GetGuildTransferInfoReq CMS 获取工会收款/转账信息
type GetGuildTransferInfoReq struct {
	g.Meta  `path:"/getGuildTransferInfo" method:"post" summary:"获取工会转账信息" tags:"直播工会"`
	GuildId uint64 `json:"guildId" v:"required#工会ID不能为空" dc:"工会ID"`
}

type GuildTransferInfoItem struct {
	GuildId   string `json:"guildId" dc:"工会ID"`
	Currency  string `json:"currency" dc:"收款币种"`
	PayeeName string `json:"payeeName" dc:"收款人姓名"`
	BankName  string `json:"bankName" dc:"银行名称"`
	AccountNo string `json:"accountNo" dc:"收款账号"`
	BankCode  string `json:"bankCode" dc:"银行代码"`
	Remark    string `json:"remark" dc:"备注"`
	UpdatedAt string `json:"updatedAt" dc:"最近更新时间"`
}

type GetGuildTransferInfoRes struct {
	Info *GuildTransferInfoItem `json:"info"`
}

// SaveGuildTransferInfoReq CMS 保存工会收款/转账信息(直写DB,币种必填)
type SaveGuildTransferInfoReq struct {
	g.Meta    `path:"/saveGuildTransferInfo" method:"post" summary:"保存工会转账信息" tags:"直播工会"`
	GuildId   uint64 `json:"guildId" v:"required#工会ID不能为空" dc:"工会ID"`
	Currency  string `json:"currency" v:"required|length:3,8#收款币种不能为空|收款币种长度需在3到8之间" dc:"收款币种(如IDR)"`
	PayeeName string `json:"payeeName" v:"max-length:128#收款人姓名最长128字符" dc:"收款人姓名"`
	BankName  string `json:"bankName" v:"max-length:128#银行名称最长128字符" dc:"银行名称"`
	AccountNo string `json:"accountNo" v:"max-length:128#收款账号最长128字符" dc:"收款账号"`
	BankCode  string `json:"bankCode" v:"max-length:64#银行代码最长64字符" dc:"银行代码"`
	Remark    string `json:"remark" v:"max-length:255#备注最长255字符" dc:"备注"`
}

type SaveGuildTransferInfoRes struct {
	Success bool `json:"success"`
}
