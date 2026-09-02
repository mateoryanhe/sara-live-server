package fiatcurrency

import (
	"xr-game-server/dao/cfgdao"
)

func Init() {
	cfgdao.InitFiatCurrencyCfgDao()
	cfgdao.ReloadFiatCurrencyCfgCache()
}
