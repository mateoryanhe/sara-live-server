package liverevenuedto

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// CMSLiveRevenueLogListReq CMS分页查询直播收益流水
type CMSLiveRevenueLogListReq struct {
	g.Meta `path:"/cmsLiveRevenueLogList" method:"post" summary:"CMS查询直播收益流水" tags:"直播收益流水"`
	httpserver.CMSQueryReq
	ReceiverId  string `json:"receiverId"  dc:"收益用户ID(可选,留空查全部)"`
	RevenueType uint8  `json:"revenueType" dc:"收益类型(0=全部,1礼物,2付费弹幕,3游戏下注)"`
	StartTime   int64  `json:"startTime"   dc:"创建时间起(秒, 0=不过滤)"`
	EndTime     int64  `json:"endTime"     dc:"创建时间止(秒, 0=不过滤)"`
}

// CMSLiveRevenueLogItem CMS直播收益流水列表项
type CMSLiveRevenueLogItem struct {
	Id               uint64     `json:"id,string"`
	RevenueType      uint8      `json:"revenueType"`
	RevenueTypeText  string     `json:"revenueTypeText"`
	RoomId           uint64     `json:"roomId,string"`
	LiveRecordId     uint64     `json:"liveRecordId,string"`
	SenderId         uint64     `json:"senderId,string"`
	SenderNickname   string     `json:"senderNickname"`
	ReceiverId       uint64     `json:"receiverId,string"`
	ReceiverNickname string     `json:"receiverNickname"`
	BizId            uint64     `json:"bizId,string"`
	BizName          string     `json:"bizName"`
	Count            int        `json:"count"`
	UnitPrice        float64    `json:"unitPrice"`
	TotalAmount      float64    `json:"totalAmount"`
	CreatedAt        *time.Time `json:"createdAt"`
}
