package dao

import (
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/cmsuserdao"
	"xr-game-server/dao/globalcfgdao"
	"xr-game-server/dao/namedao"
	"xr-game-server/dao/userinfodao"
)

func Init() {
	accountdao.InitAccountDao()
	namedao.InitName()
	globalcfgdao.InitGlobalCfg()
	cmsuserdao.InitCMSUser()
	userinfodao.InitUserInfoDao()
}
