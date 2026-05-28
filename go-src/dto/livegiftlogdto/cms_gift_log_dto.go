package livegiftlogdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"time"
	"xr-game-server/core/httpserver"
)

// CMSLiveGiftLogListReq CMS分页查询礼物流水
type CMSLiveGiftLogListReq struct {
	g.Meta `path:"/cmsLiveGiftLogList" method:"post" summary:"CMS查询礼物流水" tags:"礼物流水"`
	httpserver.CMSQueryReq
	ReceiverId string `json:"receiverId" dc:"接收者ID(可选,留空查全部)"`
	StartTime  int64  `json:"startTime"  dc:"礼物创建时间起(秒, 0=不过滤)"`
	EndTime    int64  `json:"endTime"    dc:"礼物创建时间止(秒, 0=不过滤)"`
}

// CMSLiveGiftLogItem CMS礼物流水列表项
type CMSLiveGiftLogItem struct {
	Id               uint64     `json:"id,string"`
	RoomId           uint64     `json:"roomId,string"`
	LiveRecordId     uint64     `json:"liveRecordId,string"`
	SenderId         uint64     `json:"senderId,string"`
	SenderNickname   string     `json:"senderNickname"`
	ReceiverId       uint64     `json:"receiverId,string"`
	ReceiverNickname string     `json:"receiverNickname"`
	GiftId           uint64     `json:"giftId,string"`
	GiftName         string     `json:"giftName"`
	Count            int        `json:"count"`
	UnitPrice        uint64     `json:"unitPrice"`
	TotalCost        uint64     `json:"totalCost"`
	CreatedAt        *time.Time `json:"createdAt"`
}
