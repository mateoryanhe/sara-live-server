package currencylogdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"time"
	"xr-game-server/core/httpserver"
)

// CMSCurrencyLogListReq CMS分页查询货币流水
type CMSCurrencyLogListReq struct {
	g.Meta `path:"/cmsCurrencyLogList" method:"post" summary:"CMS查询货币流水" tags:"货币流水"`
	httpserver.CMSQueryReq
	UserId       string `json:"userId" dc:"用户ID(可选,留空查全部)"`
	CurrencyType uint8  `json:"currencyType" v:"required|in:1,2#货币类型不能为空|货币类型无效" dc:"货币类型 1金币 2钻石"`
}

// CMSCurrencyLogItem CMS货币流水列表项
type CMSCurrencyLogItem struct {
	Id         uint64     `json:"id,string"`
	UserId     uint64     `json:"userId,string"`
	Nickname   string     `json:"nickname"`
	Action     uint8      `json:"action" dc:"1加 2减"`
	Amount     float64    `json:"amount"`
	Before     float64    `json:"before"`
	After      float64    `json:"after"`
	Reason     uint8      `json:"reason"`
	ReasonText string     `json:"reasonText"`
	CreatedAt  *time.Time `json:"createdAt"`
}
