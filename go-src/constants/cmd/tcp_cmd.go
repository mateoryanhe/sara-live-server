package cmd

const (
	// Enter websocket 连接
	Enter       = 1
	Heart       = 2
	Kick        = 3
	CloseServer = 4
	RepeatLogin = 5
	//无参数错误
	Error = 6
	//带参数错误
	ErrorParam = 7
	//直播间送礼推送(房间内全体在线用户)
	LiveRoomGift = 8
	//直播间文字消息推送(房间内全体在线用户)
	LiveRoomChat = 9
	//钻石余额推送(推送给指定用户)
	DiamondPush = 10
	//金币余额推送(推送给指定用户)
	GoldPush = 11
	//VIP等级推送(推送给指定用户)
	VipLevelPush = 12
	//私信推送(推送给接收者)
	PrivateMessagePush = 13
	//系统消息推送(推送给接收者)
	SystemMessagePush = 14
)
