package statdao

func Init() {
	initSysSata()
	initDailyLoginStatDao()
	initDailyUserLoginDao()
	initDailyUserRechargeDao()
	initDailyUserGoldConsumeDao()
	initDailyUserDiamondConsumeDao()
	initDailyUserAudienceDao()
}
