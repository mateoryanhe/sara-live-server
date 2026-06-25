package xrpool

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/util/gutil"
	"xr-game-server/constants/common"
	"xr-game-server/core/cfg"
)

var pool *grpool.Pool

func Init() {
	if cfg.GoPoolCfgModel.Size == common.Zero {
		pool = grpool.New(1024)
	} else {
		pool = grpool.New(cfg.GoPoolCfgModel.Size)
	}
}

// GetGoPool 防止go 过多,消耗内存

func AddWithRecover(ctx context.Context, userFunc grpool.Func) {
	pool.AddWithRecover(ctx, userFunc, func(ctx context.Context, exception error) {

		gutil.TryCatch(ctx, func(try context.Context) {
			userFunc(ctx)
		}, func(catch context.Context, exception error) {
			g.Log().Errorf(catch, "协程池发生错误 AddWithRecover error %v", gerror.Stack(exception))
		})
	})
}
