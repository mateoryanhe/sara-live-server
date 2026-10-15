package userwithdraw

import "xr-game-server/module/recharge"

func Init() {
	recharge.RegisterUserWithdrawPayoutNotifyHandler(handleUserWithdrawPayoutNotify)
	ensureUserWithdrawTierCfgGameConsumeGoldColumn()
	loadUserWithdrawTierCfgCache()
}
