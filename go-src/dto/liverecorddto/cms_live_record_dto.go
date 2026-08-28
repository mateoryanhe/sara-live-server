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
	AnchorId         string   `json:"anchorId"  dc:"主播ID(可选,留空查全部)"`
	PlatformAnchorId string   `json:"platformAnchorId" dc:"平台主播ID(可选,兼容旧参数)"`
	GuildAnchorId    string   `json:"guildAnchorId" dc:"工会主播ID(可选,兼容旧参数)"`
	AnchorIds        []string `json:"anchorIds" dc:"主播ID列表(可选,多选)"`
	LiveRecordId     string   `json:"liveRecordId" dc:"直播记录ID(可选)"`
	Keyword          string   `json:"keyword" dc:"关键字(可选,模糊匹配记录ID/主播ID/昵称)"`
	StartTime        int64    `json:"startTime" dc:"直播开始时间起(秒, 0=不过滤)"`
	EndTime          int64    `json:"endTime"   dc:"直播开始时间止(秒, 0=不过滤)"`
}

// CMSLiveRecordItem CMS直播记录列表项
type CMSLiveRecordItem struct {
	Id                           uint64     `json:"id,string"`
	AnchorId                     uint64     `json:"anchorId,string"`
	Nickname                     string     `json:"nickname"`
	Avatar                       string     `json:"avatar"`
	StartTime                    *time.Time `json:"startTime"`
	EndTime                      *time.Time `json:"endTime"`
	TotalAudience                uint64     `json:"totalAudience"`
	TotalLiveDuration            float64    `json:"totalLiveDuration"`
	TotalIncome                  float64    `json:"totalIncome"`
	TotalGiftIncome              float64    `json:"totalGiftIncome"`
	TotalPaidDanmakuIncome       float64    `json:"totalPaidDanmakuIncome"`
	TotalPrivateRoomIncome       float64    `json:"totalPrivateRoomIncome"`
	TotalPrivateRoomTicketIncome float64    `json:"totalPrivateRoomTicketIncome"`
	TotalPrivateRoomWatchIncome  float64    `json:"totalPrivateRoomWatchIncome"`
	TotalVideoCallIncome         float64    `json:"totalVideoCallIncome"`
	TotalVideoCallTicketIncome   float64    `json:"totalVideoCallTicketIncome"`
	TotalVideoCallBillingIncome  float64    `json:"totalVideoCallBillingIncome"`
	TotalGameBet                 float64    `json:"totalGameBet"`
	TotalGiftSender              uint64     `json:"totalGiftSender"`
	TotalNewFollower             uint64     `json:"totalNewFollower"`
	CreatedAt                    *time.Time `json:"createdAt"`
}
