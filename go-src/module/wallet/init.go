package wallet

import (
	"xr-game-server/core/event"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/gameevent"
)

func Init() {
	cfgdao.InitWalletExchangeCfgDao()
	cfgdao.ReloadWalletExchangeCfgCache()
	event.Sub(gameevent.UsdIncomeArrivedEvent, onUsdIncomeArrived)
}
