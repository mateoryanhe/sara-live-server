package dao

import (
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/cmsuserdao"
	"xr-game-server/dao/dailyloginstatdao"
	"xr-game-server/dao/dailyuserlogindao"
	"xr-game-server/dao/globalcfgdao"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/livefollowdao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/messagedao"
	"xr-game-server/dao/namedao"
	"xr-game-server/dao/rechargeorderdao"
	"xr-game-server/dao/shortvideocfgdao"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dao/shortvideolikedao"
	"xr-game-server/dao/shortvideowatchdao"
	"xr-game-server/dao/statdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dao/userlogindevicedao"
)

func Init() {
	accountdao.InitAccountDao()
	namedao.InitName()
	globalcfgdao.InitGlobalCfg()
	cmsuserdao.InitCMSUser()
	userinfodao.InitUserInfoDao()
	dailyloginstatdao.InitDailyLoginStatDao()
	dailyuserlogindao.InitDailyUserLoginDao()
	userlogindevicedao.InitUserLoginDeviceDao()
	guilddao.InitGuildDao()
	liveroomdao.InitLiveRoomDao()
	liveroomdao.InitLiveRoomOnlineDao()
	liveroomdao.InitLiveRecordDao()
	livefollowdao.InitLiveFollowDao()
	shortvideodao.InitShortVideoDao()
	shortvideocfgdao.InitShortVideoCfgDao()
	shortvideolikedao.InitShortVideoLikeDao()
	shortvideowatchdao.InitShortVideoWatchDao()
	rechargeorderdao.InitRechargeOrderDao()
	messagedao.Init()
	statdao.Init()
}
