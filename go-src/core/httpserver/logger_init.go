package httpserver

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"xr-game-server/core/shutdown"
)

// initHTTPServerLogger 关闭 GF Server 内置 access/error 文件日志,改由 logger.access / logger.error 独立配置处理;
// 同时将 logger.detail 同步到 httpServer.Logger() 与默认 g.Log(),避免 detail 日志丢失。
func initHTTPServerLogger() {
	ctx := context.Background()
	httpServer.SetAccessLogEnabled(false)
	httpServer.SetErrorLogEnabled(false)
	initDetailLogger(ctx)
	initAccessLogger()
	initErrorLogger()
	shutdown.RunAfterDetailLoggerHooks()
	shutdown.RegisterLoggerSyncHook(func() {
		shutdown.ApplySyncLogger(httpServer.Logger())
	})
}

func initDetailLogger(ctx context.Context) {
	if g.Cfg().MustGet(ctx, "logger.detail").IsEmpty() {
		return
	}
	detailLogger := g.Log("detail")
	detailCfg := detailLogger.GetConfig()
	_ = httpServer.Logger().SetConfig(detailCfg)
	defaultLogger := glog.New()
	if err := defaultLogger.SetConfig(detailCfg); err != nil {
		return
	}
	glog.SetDefaultLogger(defaultLogger)
}
