package xrlog

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
)

const DetailLoggerName = "detail"

// DetailLog 业务 detail 日志(对应 logger.detail 配置).
var DetailLog detailLogger

type detailLogger struct{}

func (detailLogger) Logger() *glog.Logger {
	return g.Log(DetailLoggerName)
}

func (l detailLogger) Info(ctx context.Context, v ...any) {
	l.Logger().Info(ctx, v...)
}

func (l detailLogger) Infof(ctx context.Context, format string, args ...any) {
	l.Logger().Infof(ctx, format, args...)
}

func (l detailLogger) Warning(ctx context.Context, v ...any) {
	l.Logger().Warning(ctx, v...)
}

func (l detailLogger) Warningf(ctx context.Context, format string, args ...any) {
	l.Logger().Warningf(ctx, format, args...)
}

func (l detailLogger) Error(ctx context.Context, v ...any) {
	l.Logger().Error(ctx, v...)
}

func (l detailLogger) Errorf(ctx context.Context, format string, args ...any) {
	l.Logger().Errorf(ctx, format, args...)
}
