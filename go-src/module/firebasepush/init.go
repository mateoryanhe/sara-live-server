package firebasepush

import "xr-game-server/dao/cfgdao"

func Init() {
	cfgdao.InitFirebasePushCfgDao()
	cfgdao.ReloadFirebasePushCfgCache()
}
