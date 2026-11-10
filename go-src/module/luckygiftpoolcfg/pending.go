package luckygiftpoolcfg

import (
	liveentity "xr-game-server/entity/live"
)

// enqueuePendingPoolCredit 应入奖池金额写入临时表(不直加 pool_balance)
func enqueuePendingPoolCredit(amount float64) {
	if amount <= 0 {
		return
	}
	_ = liveentity.NewLuckyGiftPoolPendingBalance(amount)
}
