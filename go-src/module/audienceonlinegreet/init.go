package audienceonlinegreet

import "xr-game-server/dao/cfgdao"

func Init() {
	cfgdao.InitAudienceOnlineAnchorGreetCfgDao()
	cfgdao.ReloadAudienceOnlineAnchorGreetCfgCache()
	initOnlineGreet()
}