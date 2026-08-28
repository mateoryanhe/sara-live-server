package liverecorddto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// CMSWeeklyUnsettledLiveListReq CMS分页查询本周未结算流水
type CMSWeeklyUnsettledLiveListReq struct {
	g.Meta `path:"/cmsWeeklyUnsettledLiveList" method:"post" summary:"CMS查询本周未结算流水" tags:"直播记录"`
	httpserver.CMSQueryReq
	AnchorId         string   `json:"anchorId"  dc:"主播ID(可选,留空查全部)"`
	PlatformAnchorId string   `json:"platformAnchorId" dc:"平台主播ID(可选,兼容旧参数)"`
	GuildAnchorId    string   `json:"guildAnchorId" dc:"工会主播ID(可选,兼容旧参数)"`
	AnchorIds        []string `json:"anchorIds" dc:"主播ID列表(可选,多选)"`
	Keyword          string   `json:"keyword" dc:"关键字(可选,模糊匹配流水ID/主播ID/昵称)"`
}

// CMSWeeklyUnsettledLiveItem CMS本周未结算流水列表项
type CMSWeeklyUnsettledLiveItem struct {
	ID                         string  `json:"id"`
	RoomId                     uint64  `json:"roomId,string"`
	RoomNickname               string  `json:"roomNickname"`
	RoomAvatar                 string  `json:"roomAvatar"`
	LiveDate                   string  `json:"liveDate"`
	WeeklyUnsettledTotalIncome float64 `json:"weeklyUnsettledTotalIncome" dc:"本周未结算总收益"`
	LiveDuration               float64 `json:"liveDuration" dc:"有效直播时长(秒)"`
	TotalIncome                float64 `json:"totalIncome" dc:"当天直播收益"`
}
