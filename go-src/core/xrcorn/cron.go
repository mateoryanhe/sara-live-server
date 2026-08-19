package xrcorn

import (
	"context"
	"sync/atomic"
	"xr-game-server/core/xrlog"

	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/util/gutil"
)

var paused atomic.Bool

// Pause 停止执行新的 cron 任务.
func Pause() {
	paused.Store(true)
}

func AddSingleton(ctx context.Context, pattern string, job gcron.JobFunc) (*gcron.Entry, error) {
	return gcron.AddSingleton(ctx, pattern, func(ctx context.Context) {
		if paused.Load() {
			return
		}
		gutil.TryCatch(ctx, func(try context.Context) {
			job(ctx)
		}, func(catch context.Context, exception error) {
			xrlog.ErrorWithErr(catch, "Cron", "AddSingleton", exception)
		})
	})
}
