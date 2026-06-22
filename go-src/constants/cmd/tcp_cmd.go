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
	//主播封禁推送(推送给主播及直播间在线观众)
	LiveRoomAnchorBan = 15
	//观众禁言状态推送(推送给被禁言/解禁的用户)
	LiveRoomAudienceMute = 17
	//观众被踢出推送(推送给被踢出的用户)
	LiveRoomAudienceKick = 18
	//取消观众进入限制推送(推送给被取消限制的用户)
	LiveRoomAudienceKickCancel = 19
	//观众进入直播间推送(房间内全体在线用户)
	LiveRoomAudienceJoin = 20
	//观众离开直播间推送(房间内全体在线用户)
	LiveRoomAudienceLeave = 21
	//直播间付费弹幕推送(房间内全体在线用户)
	LiveRoomPaidDanmaku = 22
	//主播下播推送(推送给直播间在线观众,不含主播)
	LiveRoomStopLive = 23
	//观众列表刷新推送(房间内全体在线用户,含主播)
	LiveRoomAudienceListRefresh = 24
)
