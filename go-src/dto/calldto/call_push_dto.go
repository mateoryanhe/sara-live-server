package calldto

// CallRequestPushItem 直播间通话请求推送载荷(推送给主播)
type CallRequestPushItem struct {
	OrderId        string `json:"orderId" dc:"通话订单ID"`
	CallerId       string `json:"callerId" dc:"呼叫者ID"`
	ChannelName    string `json:"channelName" dc:"声网频道名"`
	CallType       uint8  `json:"callType" dc:"通话类型(1-语音,2-视频)"`
	CallerNickname string `json:"callerNickname" dc:"呼叫者昵称"`
	CallerAvatar   string `json:"callerAvatar" dc:"呼叫者头像"`
	ReceiverToken  string `json:"receiverToken" dc:"接收者声网Token"`
	Message        string `json:"message" dc:"提示文案"`
}

// CallRejectedPushItem 通话被拒接推送载荷(推送给呼叫者)
type CallRejectedPushItem struct {
	OrderId          string `json:"orderId" dc:"通话订单ID"`
	ReceiverId       string `json:"receiverId" dc:"接收者ID"`
	ReceiverNickname string `json:"receiverNickname" dc:"接收者昵称"`
	ReceiverAvatar   string `json:"receiverAvatar" dc:"接收者头像"`
	Message          string `json:"message" dc:"提示文案"`
}

// CallAcceptedPushItem 通话被接听推送载荷(推送给呼叫者)
type CallAcceptedPushItem struct {
	OrderId          string `json:"orderId" dc:"通话订单ID"`
	ReceiverId       string `json:"receiverId" dc:"接收者ID"`
	ReceiverNickname string `json:"receiverNickname" dc:"接收者昵称"`
	ReceiverAvatar   string `json:"receiverAvatar" dc:"接收者头像"`
	ChannelName      string `json:"channelName" dc:"声网频道名"`
	CallType         uint8  `json:"callType" dc:"通话类型(1-语音,2-视频)"`
	Message          string `json:"message" dc:"提示文案"`
}

// CallEndedPushItem 通话结束推送载荷(推送给对方)
type CallEndedPushItem struct {
	OrderId         string  `json:"orderId" dc:"通话订单ID"`
	EndUserId       string  `json:"endUserId" dc:"挂断方用户ID"`
	EndUserNickname string  `json:"endUserNickname" dc:"挂断方昵称"`
	EndUserAvatar   string  `json:"endUserAvatar" dc:"挂断方头像"`
	CallDuration    uint32  `json:"callDuration" dc:"通话时长(秒)"`
	BillingDuration uint32  `json:"billingDuration" dc:"计费时长(分钟)"`
	TotalCost       float64 `json:"totalCost" dc:"总费用(钻石)"`
	Message         string  `json:"message" dc:"提示文案"`
}

// CallStartedPushItem 通话开始推送载荷(推送给呼叫者与接听者)
type CallStartedPushItem struct {
	OrderId     string `json:"orderId" dc:"通话订单ID"`
	CallerId    string `json:"callerId" dc:"呼叫者ID"`
	ReceiverId  string `json:"receiverId" dc:"接听者ID"`
	ChannelName string `json:"channelName" dc:"声网频道名"`
	CallType    uint8  `json:"callType" dc:"通话类型(1-语音,2-视频)"`
	StartedAt   int64  `json:"startedAt" dc:"通话开始时间(秒)"`
	Message     string `json:"message" dc:"提示文案"`
}

// CallTimeoutPushItem 呼叫超时推送载荷(推送给呼叫者与接听者)
type CallTimeoutPushItem struct {
	OrderId    string `json:"orderId" dc:"通话订单ID"`
	CallerId   string `json:"callerId" dc:"呼叫者ID"`
	ReceiverId string `json:"receiverId" dc:"接听者ID"`
	Message    string `json:"message" dc:"提示文案"`
}
