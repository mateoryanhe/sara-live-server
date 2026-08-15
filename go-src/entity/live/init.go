package entity

// Init 直播相关表迁移与 syndb 注册
func Init() {
	initLiveRoom()
	initLiveRoomTag()
	initLiveRoomGameRecommend()
	initLiveRoomOnline()
	initLiveRoomBillingPay()
	initLiveRecord()
	initLiveRecordUser()
	initLiveRevenueLog()
	initLiveFollow()
	initLiveTicket()
	initLivePrivateRoomBilling()
	initLiveCfg()
	initLiveGuild()
	InitLiveGift()
	initHomeBanner()
	initVipCfg()
	initAgoraCfg()
	initAliyunTextModerationCfg()
}
