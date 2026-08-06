package cfg

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
)

// initDefaultLoggerFromDetail 将 logger.detail 配置应用到默认 g.Log()。
// GF 默认读取 logger 根节点；此处启动时把 detail 子配置同步给 default 实例。
func initDefaultLoggerFromDetail() {
	ctx := gctx.New()
	if g.Cfg().MustGet(ctx, "logger.detail").IsEmpty() {
		return
	}
	detailLogger := g.Log("detail")
	defaultLogger := glog.New()
	if err := defaultLogger.SetConfig(detailLogger.GetConfig()); err != nil {
		return
	}
	glog.SetDefaultLogger(defaultLogger)
}

// EnsureErrorLogger 加载 logger.error 配置并同步写盘。
// 数据库迁移等启动早期失败会走 xrlog -> g.Log("error"),需在连库前调用。
func EnsureErrorLogger() {
	ctx := gctx.New()
	if g.Cfg().MustGet(ctx, "logger.error").IsEmpty() {
		return
	}
	applySyncLogger(g.Log("error"))
}

// disableHTTPServerBuiltinLogs 关闭 GF HTTP Server 内置 access/error 文件日志,
// 避免写入 server.accessLogPattern / server.errorLogPattern 指定的按日文件。
func disableHTTPServerBuiltinLogs() {
	ctx := gctx.New()
	if !g.Cfg().Available(ctx) {
		return
	}
	server := g.Server()
	server.SetAccessLogEnabled(false)
	server.SetErrorLogEnabled(false)
}

func applySyncLogger(logger *glog.Logger) {
	if logger == nil {
		return
	}
	cfg := logger.GetConfig()
	cfg.Flags &^= glog.F_ASYNC
	_ = logger.SetConfig(cfg)
}
