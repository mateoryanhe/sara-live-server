package luckygiftpoolcfg

import "sync/atomic"

// settleNoWinOn: 冷启/热启未合并完成、或 PrepareRestart 排空期间，
// 返奖按不中奖处理；应入池金额写入 pending 表，不直加 pool_balance。
var settleNoWinOn atomic.Bool

func isSettleNoWin() bool {
	return settleNoWinOn.Load()
}

func setSettleNoWin(on bool) {
	settleNoWinOn.Store(on)
}
