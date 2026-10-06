package liveroomdao

func Init() {
	initLiveRoomDao()
	initOneToOneRoomDao()
	initOneToOneFreeDailyDao()
	initLiveRoomIncomeDao()
	initAnchorIncomeSettlementLogDao()
	initGuildIncomeSettlementLogDao()
	initLiveRecordDao()
	initLiveRoomOnlineDao()
	initDailyAnchorEffectiveLiveDao()
	initDailyGuildEffectiveLiveDao()
	InitLiveRoomGameRecommendDao()
}
