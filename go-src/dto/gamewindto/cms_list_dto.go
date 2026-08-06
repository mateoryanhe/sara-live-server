package gamewindto

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// CMSGameWinLogListReq CMS 分页查询游戏派彩记录
type CMSGameWinLogListReq struct {
	g.Meta `path:"/cmsGameWinLogList" method:"post" summary:"CMS查询游戏派彩记录" tags:"游戏派彩"`
	httpserver.CMSQueryReq
	UserId       string `json:"userId" dc:"用户ID(可选)"`
	GameCode     string `json:"gameCode" dc:"游戏编码(可选,模糊匹配)"`
	OrderId      string `json:"orderId" dc:"订单ID(可选,模糊匹配)"`
	PlatformType string `json:"platformType" dc:"平台类型(可选,模糊匹配)"`
}

// CMSGameWinLogItem CMS 游戏派彩记录项
type CMSGameWinLogItem struct {
	Id           uint64     `json:"id,string"`
	UserId       uint64     `json:"userId,string"`
	Nickname     string     `json:"nickname"`
	GameCode     string     `json:"gameCode"`
	NameEn       string     `json:"nameEn"`
	Cover        string     `json:"cover"`
	Amount       float64    `json:"amount"`
	PlatformType string     `json:"platformType"`
	OrderId      string     `json:"orderId"`
	CreatedAt    *time.Time `json:"createdAt"`
}
