package customerservice

import "xr-game-server/dao/cfgdao"

func Init() {
	cfgdao.InitCustomerServiceCfgDao()
	cfgdao.ReloadCustomerServiceCfgCache()
}
