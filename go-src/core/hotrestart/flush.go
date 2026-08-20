package hotrestart

import (
	"context"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gutil"
	"xr-game-server/core/syndb"
	"xr-game-server/core/xrlog"
)

// startPhase1FlushBackground 第一阶段后台刷盘(recover),不阻塞热重启主流程.
func startPhase1FlushBackground(_ context.Context) {
	go func() {
		gutil.TryCatch(gctx.New(), func(_ context.Context) {
			syndb.FlushUntilIdle(0)
		}, func(catch context.Context, err error) {
			xrlog.ErrorWithErr(catch, "HotRestart", "syndb第一阶段刷盘panic", err)
		})
	}()
}

// flushPhase2Once 第二阶段退出前尽力刷一轮,不阻塞等待队列空.
func flushPhase2Once(ctx context.Context) {
	gutil.TryCatch(ctx, func(_ context.Context) {
		syndb.FlushShutdownOnce()
	}, func(catch context.Context, err error) {
		xrlog.ErrorWithErr(catch, "HotRestart", "syndb第二阶段刷盘panic", err)
	})
}
