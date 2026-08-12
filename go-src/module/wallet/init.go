package wallet

import "xr-game-server/dao/cfgdao"

func Init() {
	cfgdao.InitWalletExchangeCfgDao()
	cfgdao.ReloadWalletExchangeCfgCache()
}
