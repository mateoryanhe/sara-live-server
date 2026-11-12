package gameactivity

import "xr-game-server/dao/cfgdao"

func Init() {
	cfgdao.InitGameActivityCfgDao()
	cfgdao.ReloadGameActivityCfgCache()
}
