package gamebetdto

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// CMSGameBetLogListReq CMS 分页查询游戏下注记录
type CMSGameBetLogListReq struct {
	g.Meta `path:"/cmsGameBetLogList" method:"post" summary:"CMS查询游戏下注记录" tags:"游戏下注"`
	httpserver.CMSQueryReq
	UserId       string `json:"userId" dc:"用户ID(可选)"`
	GameCode     string `json:"gameCode" dc:"游戏编码(可选,模糊匹配)"`
	OrderId      string `json:"orderId" dc:"订单ID(可选,模糊匹配)"`
	PlatformType string `json:"platformType" dc:"平台类型(可选,模糊匹配)"`
}

// CMSGameBetLogItem CMS 游戏下注记录项
type CMSGameBetLogItem struct {
	Id             uint64     `json:"id,string"`
	UserId         uint64     `json:"userId,string"`
	Nickname       string     `json:"nickname"`
	GameCode       string     `json:"gameCode"`
	NameEn         string     `json:"nameEn"`
	Cover          string     `json:"cover"`
	Amount         float64    `json:"amount"`
	PlatformType   string     `json:"platformType"`
	OrderId        string     `json:"orderId"`
	LiveRoomId     uint64     `json:"liveRoomId,string"`
	LiveRecordId   uint64     `json:"liveRecordId,string"`
	LiveRoomTitle  string     `json:"liveRoomTitle"`
	AnchorNickname string     `json:"anchorNickname"`
	CreatedAt      *time.Time `json:"createdAt"`
}
