package syndb

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/shutdown"
	"xr-game-server/core/xrtimer"
)

// InitSynCache 初始化同步数据库缓存
func InitSynCache() {
	//先注册好通道,延迟 10 秒后开始单线程消费缓冲数据(quick 目标约 1 秒内入库)
	xrtimer.AddOnce(gctx.New(), 10*time.Second, func(ctx context.Context) {
		g.Log().Infof(ctx, "数据库同步开启成功")
		xrtimer.AddSingleton(ctx, 500*time.Millisecond, consume)
	})
	shutdown.RegCommonShutDownHandler(SysExit)
}

// RegWithSelf 指定同步频率
func RegWithSelf(tbName db.TbName, tbCol db.TbCol, synTime time.Duration, _ int) {
	colKey := string(tbName) + ":" + string(tbCol)
	lazyMap[colKey] = newColSynCache(string(tbName), string(tbCol), synTime)
}
