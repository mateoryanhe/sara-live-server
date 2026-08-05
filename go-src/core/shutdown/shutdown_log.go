package shutdown

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
)

var (
	loggerSyncHooks        []func()
	afterDetailLoggerHooks []func()
	loggerNames            = []string{"detail", "access", "error"}
)

// RegisterAfterDetailLoggerHook 注册 detail logger 初始化完成后执行的回调.
func RegisterAfterDetailLoggerHook(fn func()) {
	if fn == nil {
		return
	}
	afterDetailLoggerHooks = append(afterDetailLoggerHooks, fn)
}

// RunAfterDetailLoggerHooks 执行 detail logger 初始化后的回调.
func RunAfterDetailLoggerHooks() {
	for _, hook := range afterDetailLoggerHooks {
		hook()
	}
}

// RegisterLoggerSyncHook 注册额外 logger 的同步切换(如 httpServer.Logger)。
func RegisterLoggerSyncHook(fn func()) {
	if fn == nil {
		return
	}
	loggerSyncHooks = append(loggerSyncHooks, fn)
}

// ApplySyncLogger 关闭异步写盘,改为同步写入。
func ApplySyncLogger(logger *glog.Logger) {
	if logger == nil {
		return
	}
	cfg := logger.GetConfig()
	cfg.Flags &^= glog.F_ASYNC
	_ = logger.SetConfig(cfg)
}

func enableSyncLoggers() {
	ctx := gctx.New()
	for _, name := range loggerNames {
		v, err := g.Cfg().Get(ctx, "logger."+name)
		if err != nil || v.IsEmpty() {
			continue
		}
		ApplySyncLogger(g.Log(name))
	}
	ApplySyncLogger(g.Log())

	if db := g.DB(); db != nil {
		if logger, ok := db.GetLogger().(*glog.Logger); ok {
			ApplySyncLogger(logger)
		}
	}
	for _, hook := range loggerSyncHooks {
		hook()
	}
}
