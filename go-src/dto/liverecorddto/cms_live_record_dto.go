package liverecorddto

import (
	"github.com/gogf/gf/v2/frame/g"
	"time"
	"xr-game-server/core/httpserver"
)

// CMSLiveRecordListReq CMS分页查询直播记录
type CMSLiveRecordListReq struct {
	g.Meta `path:"/cmsLiveRecordList" method:"post" summary:"CMS查询直播记录" tags:"直播记录"`
	httpserver.CMSQueryReq
	AnchorId  string `json:"anchorId"  dc:"主播ID(可选,留空查全部)"`
	StartTime int64  `json:"startTime" dc:"直播开始时间起(秒, 0=不过滤)"`
	EndTime   int64  `json:"endTime"   dc:"直播开始时间止(秒, 0=不过滤)"`
}

// CMSLiveRecordItem CMS直播记录列表项
type CMSLiveRecordItem struct {
	Id                     uint64     `json:"id,string"`
	AnchorId               uint64     `json:"anchorId,string"`
	Nickname               string     `json:"nickname"`
	StartTime              *time.Time `json:"startTime"`
	EndTime                *time.Time `json:"endTime"`
	TotalAudience          uint64     `json:"totalAudience"`
	TotalLiveDuration      float64    `json:"totalLiveDuration"`
	TotalIncome            float64    `json:"totalIncome"`
	TotalGiftIncome        float64    `json:"totalGiftIncome"`
	TotalPrivateRoomIncome float64    `json:"totalPrivateRoomIncome"`
	TotalGameBet           float64    `json:"totalGameBet"`
	TotalGiftSender        uint64     `json:"totalGiftSender"`
	CreatedAt              *time.Time `json:"createdAt"`
}
