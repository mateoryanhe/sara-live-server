package xrcorn

import (
	"context"
	"xr-game-server/core/xrlog"

	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/util/gutil"
)

func AddSingleton(ctx context.Context, pattern string, job gcron.JobFunc) (*gcron.Entry, error) {
	return gcron.AddSingleton(ctx, pattern, func(ctx context.Context) {
		gutil.TryCatch(ctx, func(try context.Context) {
			job(ctx)
		}, func(catch context.Context, exception error) {
			xrlog.ErrorWithErr(catch, "Cron", "AddSingleton", exception)
		})
	})
}
