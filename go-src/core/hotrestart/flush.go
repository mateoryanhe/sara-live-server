package hotrestart

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gutil"
	"xr-game-server/core/syndb"
	"xr-game-server/core/xrlog"
	"xr-game-server/core/xrpool"
)

var (
	statFlushMu           sync.Mutex
	statDrainFn           func()
	statIdleFn            func() bool
	trackingFlushMu       sync.Mutex
	trackingDrainFn       func()
	trackingIdleFn        func() bool
	phase1IdlePoll        = 50 * time.Millisecond
	phase1FlushDone       atomic.Bool
	oldProcessFlushStop   atomic.Bool
	oldProcessFlushWorker sync.Once
)

// RegisterStatQueueFlush 由 module/stat 注册:热重启刷盘时先排空统计队列再落 syndb.
func RegisterStatQueueFlush(drain func(), idle func() bool) {
	statFlushMu.Lock()
	defer statFlushMu.Unlock()
	statDrainFn = drain
	statIdleFn = idle
}

// RegisterTrackingQueueFlush 由 module/tracking 注册:热重启刷盘时排空埋点队列再落 syndb.
func RegisterTrackingQueueFlush(drain func(), idle func() bool) {
	trackingFlushMu.Lock()
	defer trackingFlushMu.Unlock()
	trackingDrainFn = drain
	trackingIdleFn = idle
}

func drainStatQueue() {
	statFlushMu.Lock()
	fn := statDrainFn
	statFlushMu.Unlock()
	if fn != nil {
		fn()
	}
}

func statQueueIdle() bool {
	statFlushMu.Lock()
	fn := statIdleFn
	statFlushMu.Unlock()
	if fn == nil {
		return true
	}
	return fn()
}

func drainTrackingQueue() {
	trackingFlushMu.Lock()
	fn := trackingDrainFn
	trackingFlushMu.Unlock()
	if fn != nil {
		fn()
	}
}

func trackingQueueIdle() bool {
	trackingFlushMu.Lock()
	fn := trackingIdleFn
	trackingFlushMu.Unlock()
	if fn == nil {
		return true
	}
	return fn()
}

func oldProcessQueuesIdle() bool {
	return syndb.AllCachesIdle() && statQueueIdle() && trackingQueueIdle()
}

// flushOldProcessRound 旧进程一轮:业务队列排空 + syndb 落库.
func flushOldProcessRound() int {
	drainStatQueue()
	drainTrackingQueue()
	return syndb.FlushShutdownOnce()
}

// startPhase1FlushBackground 旧进程自 PrepareRestart 起持续刷盘,直至退出前 stop,防在途入队不落库.
func startPhase1FlushBackground(_ context.Context) {
	oldProcessFlushStop.Store(false)
	phase1FlushDone.Store(false)
	oldProcessFlushWorker.Do(func() {
		xrpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
			runOldProcessContinuousFlush(ctx)
		})
	})
}

func markPhase1FlushDone() {
	phase1FlushDone.Store(true)
}

// waitHotRestartPhase1 第一阶段固定等待 maxWait(配置 hotRestartFlushTimeout),后台持续刷盘.
// 不因队列提前 idle 而缩短,给在途请求与异步入队留足落库窗口,再 RestartAllServer.
func waitHotRestartPhase1(ctx context.Context, maxWait time.Duration) {
	const poll = 100 * time.Millisecond
	start := time.Now()
	deadline := start.Add(maxWait)
	for time.Now().Before(deadline) {
		time.Sleep(poll)
	}
	if phase1FlushDone.Load() {
		xrlog.DetailLog.Info(ctx, "热重启第一阶段刷盘已结束,进入 RestartAllServer")
	} else {
		xrlog.DetailLog.Warningf(ctx, "热重启第一阶段等待刷盘超时(%ds),仍继续 RestartAllServer", int(maxWait.Seconds()))
	}
	xrlog.DetailLog.Info(ctx, "hot_restart phase1 ready")
}

// runOldProcessContinuousFlush 循环刷盘,覆盖阶段一/fork 后等子进程/阶段二全程,直到 stop.
func runOldProcessContinuousFlush(ctx context.Context) {
	for !oldProcessFlushStop.Load() {
		rows := flushOldProcessRound()
		if rows == 0 && oldProcessQueuesIdle() {
			if !phase1FlushDone.Load() {
				time.Sleep(phase1IdlePoll)
				if flushOldProcessRound() == 0 && oldProcessQueuesIdle() {
					xrlog.DetailLog.Info(ctx, "热重启旧进程刷盘曾全部空闲(后台仍持续刷盘)")
					xrlog.DetailLog.Info(ctx, "hot_restart old_process flush idle once")
					markPhase1FlushDone()
				}
			}
		}
		time.Sleep(phase1IdlePoll)
	}
	flushOldProcessRound()
}

// stopOldProcessContinuousFlush 退出前停止后台循环并做最后一轮刷盘.
func stopOldProcessContinuousFlush(ctx context.Context) {
	oldProcessFlushStop.Store(true)
	time.Sleep(phase1IdlePoll * 2)
	flushPhase2Once(ctx)
}

// flushPhase2Once 旧进程退出前尽力刷盘(与后台循环互补).
func flushPhase2Once(ctx context.Context) {
	gutil.TryCatch(ctx, func(_ context.Context) {
		flushOldProcessRound()
	}, func(catch context.Context, err error) {
		xrlog.ErrorWithErr(catch, "HotRestart", "syndb旧进程退出刷盘panic", err)
	})
}
