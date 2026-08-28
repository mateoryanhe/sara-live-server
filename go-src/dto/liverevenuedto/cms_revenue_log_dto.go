package liverevenuedto

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// CMSLiveRevenueLogListReq CMS分页查询直播间社交流水
type CMSLiveRevenueLogListReq struct {
	g.Meta `path:"/cmsLiveRevenueLogList" method:"post" summary:"CMS查询直播间社交流水" tags:"直播社交流水"`
	httpserver.CMSQueryReq
	ReceiverId       string   `json:"receiverId"  dc:"收益用户ID(可选,留空查全部,兼容旧参数)"`
	PlatformAnchorId string   `json:"platformAnchorId" dc:"平台主播ID(可选,兼容旧参数)"`
	GuildAnchorId    string   `json:"guildAnchorId" dc:"工会主播ID(可选,兼容旧参数)"`
	ReceiverIds      []string `json:"receiverIds" dc:"收益用户ID列表(可选,多选)"`
	LiveRecordId     string   `json:"liveRecordId" dc:"直播记录ID(可选)"`
	Keyword          string   `json:"keyword" dc:"关键字(可选,模糊匹配流水ID/直播记录ID/主播ID/付款用户ID/昵称)"`
	RevenueType      uint8    `json:"revenueType" dc:"流水类型(0=全部,1礼物,2付费弹幕,4私密房计费,5门票,6视频通话门票,7视频通话计费)"`
	StartTime        int64    `json:"startTime"   dc:"创建时间起(秒, 0=不过滤)"`
	EndTime          int64    `json:"endTime"     dc:"创建时间止(秒, 0=不过滤)"`
}

// CMSLiveRevenueLogItem CMS直播间社交流水列表项
type CMSLiveRevenueLogItem struct {
	Id               uint64     `json:"id,string"`
	RevenueType      uint8      `json:"revenueType"`
	RevenueTypeText  string     `json:"revenueTypeText"`
	RoomId           uint64     `json:"roomId,string"`
	LiveRecordId     uint64     `json:"liveRecordId,string"`
	SenderId         uint64     `json:"senderId,string"`
	SenderNickname   string     `json:"senderNickname"`
	SenderAvatar     string     `json:"senderAvatar"`
	ReceiverId       uint64     `json:"receiverId,string"`
	ReceiverNickname string     `json:"receiverNickname"`
	ReceiverAvatar   string     `json:"receiverAvatar"`
	BizId            uint64     `json:"bizId,string"`
	BizName          string     `json:"bizName"`
	Count            int        `json:"count"`
	UnitPrice        float64    `json:"unitPrice"`
	TotalAmount      float64    `json:"totalAmount"`
	Status           uint8      `json:"status"`
	StatusText       string     `json:"statusText"`
	CreatedAt        *time.Time `json:"createdAt"`
}
