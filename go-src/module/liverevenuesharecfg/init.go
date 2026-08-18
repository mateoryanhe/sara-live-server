package liverevenuesharecfg

import "xr-game-server/dao/cfgdao"

func Init() {
	cfgdao.InitLiveRevenueShareCfgDao()
	cfgdao.ReloadLiveRevenueShareCfgCache()
}
