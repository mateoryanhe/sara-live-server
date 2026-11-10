package luckygiftpoolcfg

import (
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dao/cfgdao"
)

// 启动成功后延时合并临时入账，避开启动瞬时并发与 syndb 刷盘窗口
const pendingMergeDelay = 10 * time.Second

var mergeScheduleOnce sync.Once

func schedulePendingMergeAfterStartOnce() {
	mergeScheduleOnce.Do(func() {
		xrtimer.AddOnce(gctx.New(), pendingMergeDelay, func(ctx context.Context) {
			EnqueueMergePendingBalances()
		})
	})
}

// EnqueueMergePendingBalances 入队合并未处理临时余额(串行,防并发)
func EnqueueMergePendingBalances() {
	_ = enqueueAndWait(&poolJob{kind: jobMergePending})
}

func resetStuckPendingOnStart() {
	cfgdao.ResetStuckLuckyGiftPoolPendingProcessing()
}

func processMergePendingLocked() {
	cfgdao.ResetStuckLuckyGiftPoolPendingProcessing()
	rows := cfgdao.LoadPendingLuckyGiftPoolBalances()
	if len(rows) == 0 {
		setSettleNoWin(false)
		return
	}

	runtimeMu.Lock()
	defer runtimeMu.Unlock()

	balance := runtimePoolBalanceLocked()
	changed := false
	for _, row := range rows {
		if row == nil || row.ID == 0 || row.Amount <= 0 {
			continue
		}
		if !cfgdao.ClaimLuckyGiftPoolPendingBalance(row.ID) {
			continue
		}
		balance += row.Amount
		changed = true
		_ = cfgdao.MarkLuckyGiftPoolPendingBalanceProcessed(row.ID)
	}
	if changed {
		setRuntimePoolBalanceLocked(balance)
		storePublicSnapshotLocked()
	}
	setSettleNoWin(false)
}
