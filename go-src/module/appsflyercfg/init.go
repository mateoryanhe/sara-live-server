package appsflyercfg

import "xr-game-server/dao/cfgdao"

func Init() {
	cfgdao.InitAppsFlyerCfgDao()
	cfgdao.ReloadAppsFlyerCfgCache()
}