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
	GameCode     string   `json:"gameCode" dc:"游戏编码(可选,模糊匹配)"`
	GameCodes    []string `json:"gameCodes" dc:"游戏编码(可选,多选精确匹配,优先于 gameCode)"`
	OrderId      string `json:"orderId" dc:"订单ID(可选,模糊匹配)"`
	PlatformType string `json:"platformType" dc:"平台类型(可选,精确匹配,忽略大小写)"`
}

// CMSGameWinLogListRes CMS 游戏派彩记录列表响应
type CMSGameWinLogListRes struct {
	Total     int                  `json:"total" dc:"总数"`
	Data      []*CMSGameWinLogItem `json:"data" dc:"列表"`
	Platforms []string             `json:"platforms" dc:"平台下拉(来自游戏库,大写)"`
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
	BetAmount    float64    `json:"betAmount" dc:"对应下注金额(按同 orderId 匹配下注记录)"`
	PlatformType string     `json:"platformType"`
	OrderId      string     `json:"orderId"`
	CreatedAt    *time.Time `json:"createdAt"`
}
