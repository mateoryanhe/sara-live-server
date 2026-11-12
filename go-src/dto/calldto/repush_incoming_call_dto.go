package calldto

import "github.com/gogf/gf/v2/frame/g"

// RepushIncomingCallReq 若存在呼叫中订单,向当前用户重推呼叫消息(WS cmd=LiveRoomCallRequest).
type RepushIncomingCallReq struct {
	g.Meta `path:"/repushIncomingCall" method:"post" summary:"重推进行中的来电(恢复1v1/通话流程)" tags:"通话"`
}

// RepushIncomingCallRes 重推结果; calling=false 表示当前无人呼叫.
type RepushIncomingCallRes struct {
	Calling        bool   `json:"calling" dc:"是否存在呼叫中订单并已重推"`
	OrderId        string `json:"orderId,omitempty" dc:"通话订单ID"`
	Source         uint8  `json:"source,omitempty" dc:"业务来源(1-直播间,2-私信,3-1v1房间)"`
	CallerId       string `json:"callerId,omitempty" dc:"呼叫者ID"`
	ChannelName    string `json:"channelName,omitempty" dc:"声网频道名"`
	CallType       uint8  `json:"callType,omitempty" dc:"通话类型(1-语音,2-视频)"`
	CallerNickname string `json:"callerNickname,omitempty" dc:"呼叫者昵称"`
	CallerAvatar   string `json:"callerAvatar,omitempty" dc:"呼叫者头像"`
	ReceiverToken  string `json:"receiverToken,omitempty" dc:"接收者声网Token"`
	Message        string `json:"message,omitempty" dc:"提示文案"`
}
