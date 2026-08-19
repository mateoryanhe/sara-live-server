package guilddto

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/accountdto"
)

// GetGuildDailyEffectiveLiveListReq CMS分页查询工会每日流水
type GetGuildDailyEffectiveLiveListReq struct {
	g.Meta `path:"/getGuildDailyEffectiveLiveList" method:"post" summary:"CMS查询工会每日流水" tags:"工会"`
	httpserver.CMSQueryReq
	GuildId       uint64 `json:"guildId"       v:"required#工会ID不能为空" dc:"工会ID"`
	LiveDateStart string `json:"liveDateStart" dc:"日期起(YYYY-MM-DD,可选)"`
	LiveDateEnd   string `json:"liveDateEnd"   dc:"日期止(YYYY-MM-DD,可选)"`
	Settled       int8   `json:"settled"       dc:"结算状态(-1全部,0未结算,1已结算)"`
}

// GuildDailyEffectiveLiveItem 工会每日流水
type GuildDailyEffectiveLiveItem struct {
	ID           string  `json:"id"`
	GuildId      uint64  `json:"guildId,string"`
	LiveDate     string  `json:"liveDate"`
	LiveDuration float64 `json:"liveDuration"`
	Settled      bool    `json:"settled"`
	accountdto.LiveRoomIncomeAmountsItem
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}
