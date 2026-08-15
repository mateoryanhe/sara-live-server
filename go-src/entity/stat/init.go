package entity

// Init 日/周/月统计相关表迁移与 syndb 注册
func Init() {
	initDailyLoginStat()
	initDailyUserLogin()
	initDailyUserRecharge()
	initDailyUserGoldConsume()
	initDailyUserDiamondConsume()
	initDailyUserAudience()
	initWeeklyLoginStat()
	initWeeklyUserLogin()
	initWeeklyUserRecharge()
	initWeeklyUserGoldConsume()
	initWeeklyUserDiamondConsume()
	initWeeklyUserAudience()
	initMonthlyLoginStat()
	initMonthlyUserLogin()
	initMonthlyUserRecharge()
	initMonthlyUserGoldConsume()
	initMonthlyUserDiamondConsume()
	initMonthlyUserAudience()
	initSystemTotalStat()
}
