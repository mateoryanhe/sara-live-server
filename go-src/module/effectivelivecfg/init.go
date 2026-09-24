package effectivelivecfg

import "xr-game-server/dao/cfgdao"

func Init() {
	cfgdao.InitEffectiveLiveCfgDao()
	cfgdao.ReloadEffectiveLiveCfgCache()
}
