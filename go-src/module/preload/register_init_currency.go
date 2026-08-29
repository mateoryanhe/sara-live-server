package preload

import (
	"xr-game-server/constants/currency"
	"xr-game-server/core/event"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/gameevent"
	"xr-game-server/module/wallet"
)

func initRegisterInitCurrency() {
	event.Sub(gameevent.RegisterEvent, onRegisterInitCurrency)
}

func onRegisterInitCurrency(data any) {
	val, ok := data.(*gameevent.RegisterEventData)
	if !ok || val == nil || val.UserId == 0 {
		return
	}
	userinfodao.GetUserInfoByUserId(val.UserId)
	if gold := cfgdao.GetInitGold(); gold > 0 {
		_, _ = wallet.GoldAdd(val.UserId, gold, currency.ReasonSystemGrant)
	}
	if diamond := cfgdao.GetInitDiamond(); diamond > 0 {
		_, _ = wallet.DiamondAdd(val.UserId, diamond, currency.ReasonSystemGrant)
	}
}
