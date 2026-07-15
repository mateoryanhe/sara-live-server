package dao

import (
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/calldao"
	"xr-game-server/dao/cmsuserdao"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/livefollowdao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/messagedao"
	"xr-game-server/dao/namedao"
	"xr-game-server/dao/rechargeorderdao"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dao/statdao"
	"xr-game-server/dao/userchanneltokendao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dao/userlogindevicedao"
)

func Init() {
	accountdao.InitAccountDao()
	namedao.InitName()
	cmsuserdao.InitCMSUser()
	userinfodao.InitUserInfoDao()
	userlogindevicedao.InitUserLoginDeviceDao()
	guilddao.InitGuildDao()
	liveroomdao.Init()
	livefollowdao.InitLiveFollowDao()
	shortvideodao.Init()
	rechargeorderdao.InitRechargeOrderDao()
	calldao.Init()
	userchanneltokendao.Init()
	messagedao.Init()
	statdao.Init()
}
