package hotrestart

import (
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gutil"
	"xr-game-server/core/syndb"
	"xr-game-server/core/xrlog"
)

var (
	statFlushMu    sync.Mutex
	statDrainFn    func()
	statIdleFn     func() bool
	phase1IdlePoll = 50 * time.Millisecond
)

// RegisterStatQueueFlush 由 module/stat 注册:热重启刷盘时先排空统计队列再落 syndb.
func RegisterStatQueueFlush(drain func(), idle func() bool) {
	statFlushMu.Lock()
	defer statFlushMu.Unlock()
	statDrainFn = drain
	statIdleFn = idle
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

// startPhase1FlushBackground 第一阶段后台刷盘(recover),不阻塞热重启主流程.
func startPhase1FlushBackground(_ context.Context) {
	go func() {
		gutil.TryCatch(gctx.New(), func(ctx context.Context) {
			flushPhase1UntilIdle(ctx)
		}, func(catch context.Context, err error) {
			xrlog.ErrorWithErr(catch, "HotRestart", "syndb第一阶段刷盘panic", err)
		})
	}()
}

// flushPhase1UntilIdle 先排空 stat 队列再刷 syndb,循环直到两者皆空.
func flushPhase1UntilIdle(ctx context.Context) {
	for {
		drainStatQueue()
		rows := syndb.FlushShutdownOnce()
		if rows == 0 && syndb.AllCachesIdle() && statQueueIdle() {
			// 短暂等待,吸收竞态入队后再确认一轮
			time.Sleep(phase1IdlePoll)
			drainStatQueue()
			if syndb.FlushShutdownOnce() == 0 && syndb.AllCachesIdle() && statQueueIdle() {
				xrlog.DetailLog.Info(ctx, "热重启第一阶段刷盘完成(含stat队列)")
				return
			}
		}
	}
}

// flushPhase2Once 第二阶段退出前尽力刷一轮,不阻塞等待队列空.
func flushPhase2Once(ctx context.Context) {
	gutil.TryCatch(ctx, func(_ context.Context) {
		drainStatQueue()
		syndb.FlushShutdownOnce()
	}, func(catch context.Context, err error) {
		xrlog.ErrorWithErr(catch, "HotRestart", "syndb第二阶段刷盘panic", err)
	})
}
