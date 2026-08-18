package guilddto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// CMSGuildAnchorIncomeSettlementLogListReq CMS分页查询指定工会名下主播结算流水
type CMSGuildAnchorIncomeSettlementLogListReq struct {
	g.Meta `path:"/cmsGuildAnchorIncomeSettlementLogList" method:"post" summary:"CMS查询工会名下主播结算流水" tags:"直播工会"`
	httpserver.CMSQueryReq
	GuildId   string `json:"guildId"   dc:"工会ID(必填)"`
	RoomId    string `json:"roomId"    dc:"直播间ID(可选)"`
	StartTime int64  `json:"startTime" dc:"创建时间起(秒, 0=不过滤)"`
	EndTime   int64  `json:"endTime"   dc:"创建时间止(秒, 0=不过滤)"`
}
