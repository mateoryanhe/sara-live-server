package guilddto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type GetMyGuildAnchorDailyEffectiveLiveListReq struct {
	g.Meta `path:"/getMyGuildAnchorDailyEffectiveLiveList" method:"post" summary:"获取当前CMS用户管理的工会名下主播每日流水" tags:"直播工会"`
	httpserver.CMSQueryReq
	GuildId       uint64 `json:"guildId" v:"required#工会ID不能为空" dc:"工会ID"`
	AnchorId      uint64 `json:"anchorId" v:"required#主播ID不能为空" dc:"主播用户ID(==roomId)"`
	LiveDateStart string `json:"liveDateStart" dc:"日期起(YYYY-MM-DD,可选)"`
	LiveDateEnd   string `json:"liveDateEnd" dc:"日期止(YYYY-MM-DD,可选)"`
	Settled       int8   `json:"settled" dc:"结算状态(-1全部,0未结算,1已结算)"`
}
