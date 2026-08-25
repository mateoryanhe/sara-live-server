package guilddto

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/accountdto"
)

// CMSGuildAnchorDailyEffectiveLiveListReq CMS分页查询工会名下主播每日流水
type CMSGuildAnchorDailyEffectiveLiveListReq struct {
	g.Meta `path:"/cmsGuildAnchorDailyEffectiveLiveList" method:"post" summary:"CMS查询工会名下主播每日流水" tags:"直播工会"`
	httpserver.CMSQueryReq
	GuildId       string `json:"guildId"       v:"required#工会ID不能为空" dc:"工会ID"`
	RoomId        string `json:"roomId"        dc:"主播ID/直播间ID(可选)"`
	LiveDateStart string `json:"liveDateStart" dc:"日期起(YYYY-MM-DD,可选)"`
	LiveDateEnd   string `json:"liveDateEnd"   dc:"日期止(YYYY-MM-DD,可选)"`
	Settled       int8   `json:"settled"       dc:"结算状态(-1全部,0未结算,1已结算)"`
}

// GuildAnchorDailyEffectiveLiveItem 工会名下主播每日流水
type GuildAnchorDailyEffectiveLiveItem struct {
	ID           string  `json:"id"`
	RoomId       uint64  `json:"roomId,string"`
	RoomNickname         string  `json:"roomNickname"`
	RoomAvatar           string  `json:"roomAvatar"`
	UnsettledTotalIncome float64 `json:"unsettledTotalIncome"`
	LiveDate             string  `json:"liveDate"`
	LiveDuration float64 `json:"liveDuration"`
	Settled      bool    `json:"settled"`
	accountdto.LiveRoomIncomeAmountsItem
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}
