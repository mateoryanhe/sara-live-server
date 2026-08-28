package calldto

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// CMSVideoCallLogListReq CMS分页查询视频通话日志
type CMSVideoCallLogListReq struct {
	g.Meta `path:"/cmsVideoCallLogList" method:"post" summary:"CMS查询视频通话日志" tags:"视频通话日志"`
	httpserver.CMSQueryReq
	CallerId         string   `json:"callerId"   dc:"呼叫者ID(可选,留空查全部)"`
	ReceiverId       string   `json:"receiverId" dc:"接收者ID(可选,留空查全部,兼容旧参数)"`
	PlatformAnchorId string   `json:"platformAnchorId" dc:"平台主播ID(可选,兼容旧参数)"`
	GuildAnchorId    string   `json:"guildAnchorId" dc:"工会主播ID(可选,兼容旧参数)"`
	ReceiverIds      []string `json:"receiverIds" dc:"接收者ID列表(可选,多选)"`
	Source           uint8    `json:"source"     dc:"来源(0=全部,1=直播间,2=私信)"`
	Status     uint8  `json:"status"     dc:"订单状态(0=全部,7=心跳超时,8=钻石不足等)"`
	StartTime  int64  `json:"startTime"  dc:"呼叫开始时间起(秒, 0=不过滤)"`
	EndTime    int64  `json:"endTime"    dc:"呼叫开始时间止(秒, 0=不过滤)"`
}

// CMSVideoCallLogItem CMS视频通话日志列表项
type CMSVideoCallLogItem struct {
	Id                uint64     `json:"id,string"`
	CallerId          uint64     `json:"callerId,string"`
	CallerNickname    string     `json:"callerNickname"`
	CallerAvatar      string     `json:"callerAvatar"`
	ReceiverId         uint64     `json:"receiverId,string"`
	ReceiverNickname   string     `json:"receiverNickname"`
	ReceiverIsAnchor   bool       `json:"receiverIsAnchor"`
	CallType           uint8      `json:"callType"`
	CallTypeText      string     `json:"callTypeText"`
	Source            uint8      `json:"source"`
	SourceText        string     `json:"sourceText"`
	Status            uint8      `json:"status"`
	StatusText        string     `json:"statusText"`
	CallStartTime     *time.Time `json:"callStartTime"`
	AnswerTime        *time.Time `json:"answerTime"`
	CallerHeartTime   *time.Time `json:"callerHeartTime"`
	ReceiverHeartTime *time.Time `json:"receiverHeartTime"`
	OrderEndTime      *time.Time `json:"orderEndTime"`
	CallDuration      uint32     `json:"callDuration"`
	TicketPrice       float64    `json:"ticketPrice"`
	PricePerMinute    float64    `json:"pricePerMinute"`
	TotalCost         float64    `json:"totalCost"`
	BillingDuration   uint32     `json:"billingDuration"`
	ChargeTime        *time.Time `json:"chargeTime"`
	CreatedAt         *time.Time `json:"createdAt"`
}
