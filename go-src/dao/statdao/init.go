package statdao

func Init() {
	initSysSata()
	initDailyLoginStatDao()
	initDailyUserLoginDao()
	initDailyUserRechargeDao()
	initDailyUserGoldConsumeDao()
	initDailyUserDiamondConsumeDao()
	initDailyUserAudienceDao()
	initWeeklyLoginStatDao()
	initWeeklyUserLoginDao()
	initWeeklyUserRechargeDao()
	initWeeklyUserGoldConsumeDao()
	initWeeklyUserDiamondConsumeDao()
	initWeeklyUserAudienceDao()
	initMonthlyLoginStatDao()
	initMonthlyUserLoginDao()
	initMonthlyUserRechargeDao()
	initMonthlyUserGoldConsumeDao()
	initMonthlyUserDiamondConsumeDao()
	initMonthlyUserAudienceDao()
}
