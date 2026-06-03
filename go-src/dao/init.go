package dao

import (
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/cmsuserdao"
	"xr-game-server/dao/dailyloginstatdao"
	"xr-game-server/dao/dailyuseraudiencedao"
	"xr
	"xr-game-server/dao/dailyuserdiamondconsumdao"
	"xr-game-server/dao/dailyuseraudiencedao"
	"xr-game-server/dao/dailyuserlogindao"
	"xr-game-server/dao/monthlyloginstatdao"
	"xr-game-server/dao/userlogindevicedao"

	"xr-game-server/dao/dailyuserrechargdao"
	"xr-game-server/dao/globalcfgdao"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/livefollowdao"
	"xr-game-server/dao/liveroomdao"
r-game-server/dao/messagedao"
	"
	"xr-game-server/dao/messagedao"
	"xr-game-server/dao/monthlyusergoldconsumdao"
	"xr-game-server/dao/monthlyuseraudiencedao"
	"xr-game-server/dao/monthlyuserlogindao"

	"xr-game-server/dao/monthlyuserrechargdao"
	"xr-game-server/dao/namedao"
	"xr-game-server/dao/rechargeorderdao"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dao/statdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dao/weeklyloginstatdao"
	"xr-game-server/dao/weeklyuseraudiencedao"
	"xr-game-server/dao/weeklyuserlogindao"

	"xr-game-server/dao/weeklyuserdiamondconsumdao"
	"xr-game-server/dao/weeklyusergoldconsumdao"
	"xr-game-server/dao/weeklyuserrechargdao"
)

func Init() {
	accountdao.InitAccountDao()
	namedao.InitName()
	globalcfgdao.InitGlobalCfg()
	cmsuserdao.InitCMSUser()
	userinfodao.InitUserInfoDao()
	dailyloginstatdao.InitDailyLoginStatDao()
	dailyuserlogindao.InitDailyUserLoginDao()
	dailyuserrechargdao.InitDailyUserRechargeDao()
	dailyusergoldconsumdao.InitDailyUserGoldConsumeDao()
	dailyuserdiamondconsumdao.InitDailyUserDiamondConsumeDao()
	dailyuseraudiencedao.InitDailyUserAudienceDao()
	weeklyloginstatdao.InitWeeklyLoginStatDao()
	weeklyuserlogindao.InitWeeklyUserLoginDao()
	weeklyuserrechargdao.InitWeeklyUserRechargeDao()
	weeklyusergoldconsumdao.InitWeeklyUserGoldConsumeDao()
	weeklyuserdiamondconsumdao.InitWeeklyUserDiamondConsumeDao()
	weeklyuseraudiencedao.InitWeeklyUserAudienceDao()
	monthlyloginstatdao.InitMonthlyLoginStatDao()
	monthlyuserlogindao.InitMonthlyUserLoginDao()
	monthlyuserrechargdao.InitMonthlyUserRechargeDao()
	monthlyusergoldconsumdao.InitMonthlyUserGoldConsumeDao()
	monthlyuserdiamondconsumdao.InitMonthlyUserDiamondConsumeDao()
	monthlyuseraudiencedao.InitMonthlyUserAudienceDao()
	userlogindevicedao.InitUserLoginDeviceDao()
	guilddao.InitGuildDao()
	liveroomdao.InitLiveRoomDao()
	liveroomdao.InitLiveRoomOnlineDao()
	liveroomdao.InitLiveRecordDao()
	livefollowdao.InitLiveFollowDao()
	shortvideodao.InitShortVideoDao()
	shortvideodao.InitShortVideoStatDao()
	shortvideodao.InitShortVideoCfgDao()
	shortvideodao.InitShortVideoLikeDao()
	shortvideodao.InitShortVideoWatchDao()
	rechargeorderdao.InitRechargeOrderDao()
	messagedao.Init()
	statdao.Init()
}
