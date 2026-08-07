package syndb

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/shutdown"
	"xr-game-server/core/xrlog"
	"xr-game-server/core/xrtimer"
)

var synCfg runtimeConfig

// InitSynCache 初始化同步数据库缓存
func InitSynCache() {
	synCfg = loadRuntimeConfig()
	xrtimer.AddOnce(gctx.New(), 10*time.Second, func(ctx context.Context) {
		xrlog.DetailLog.Infof(ctx, "syndb开启,tick=%vms,cpuIdle>=%v%%,maxWait=%vms,batch=%v",
			synCfg.tickInterval.Milliseconds(),
			int(synCfg.cpuIdlePercent),
			synCfg.maxPendingWait.Milliseconds(),
			synCfg.batchSize,
		)
		xrtimer.AddSingleton(ctx, synCfg.tickInterval, consume)
	})
	shutdown.RegCommonShutDownHandler(SysExit)
}
