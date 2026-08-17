package accountdto

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// GetAnchorDailyEffectiveLiveListReq CMS分页查询主播每日直播时长
type GetAnchorDailyEffectiveLiveListReq struct {
	g.Meta `path:"/getAnchorDailyEffectiveLiveList" method:"post" summary:"CMS查询主播每日直播时长" tags:"账号"`
	httpserver.CMSQueryReq
	AnchorId      uint64 `json:"anchorId"      v:"required#主播ID不能为空" dc:"主播用户ID(==roomId)"`
	LiveDateStart string `json:"liveDateStart" dc:"日期起(YYYY-MM-DD,可选)"`
	LiveDateEnd   string `json:"liveDateEnd"   dc:"日期止(YYYY-MM-DD,可选)"`
	Settled       int8   `json:"settled"       dc:"结算状态(-1全部,0未结算,1已结算)"`
}

// AnchorDailyEffectiveLiveItem 主播每日直播时长
type AnchorDailyEffectiveLiveItem struct {
	ID           string     `json:"id"`
	RoomId       uint64     `json:"roomId,string"`
	LiveDate     string     `json:"liveDate"`
	LiveDuration float64    `json:"liveDuration"`
	Settled      bool       `json:"settled"`
	CreatedAt    *time.Time `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
}
