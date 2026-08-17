package guilddto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// CMSMyGuildAnchorIncomeSettlementLogListReq CMS分页查询当前用户名下工会的主播结算流水
type CMSMyGuildAnchorIncomeSettlementLogListReq struct {
	g.Meta `path:"/cmsMyGuildAnchorIncomeSettlementLogList" method:"post" summary:"CMS查询当前用户工会主播结算流水" tags:"直播工会"`
	httpserver.CMSQueryReq
	GuildId   string `json:"guildId"   dc:"工会ID(可选,须为当前用户管理的工会)"`
	RoomId    string `json:"roomId"    dc:"直播间ID(可选)"`
	StartTime int64  `json:"startTime" dc:"创建时间起(秒, 0=不过滤)"`
	EndTime   int64  `json:"endTime"   dc:"创建时间止(秒, 0=不过滤)"`
}
