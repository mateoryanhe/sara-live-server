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
	_ = g.Log()
	detailLogger := g.Log("detail")
	defaultLogger := g.Log()
	if err := defaultLogger.SetConfig(detailLogger.GetConfig()); err != nil {
		return
	}
	glog.SetDefaultLogger(defaultLogger)
}
