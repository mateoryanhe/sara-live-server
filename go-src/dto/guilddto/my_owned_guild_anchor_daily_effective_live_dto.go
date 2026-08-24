package guilddto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type CMSMyGuildAnchorDailyEffectiveLiveListReq struct {
	g.Meta `path:"/cmsMyGuildAnchorDailyEffectiveLiveList" method:"post" summary:"CMS查询当前用户管理的工会名下主播每日流水" tags:"直播工会"`
	httpserver.CMSQueryReq
	RoomId        string `json:"roomId" dc:"主播ID/直播间ID(可选)"`
	LiveDateStart string `json:"liveDateStart" dc:"日期起(YYYY-MM-DD,可选)"`
	LiveDateEnd   string `json:"liveDateEnd" dc:"日期止(YYYY-MM-DD,可选)"`
	Settled       int8   `json:"settled" dc:"结算状态(-1全部,0未结算,1已结算)"`
}
