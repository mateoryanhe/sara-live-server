package httpserver

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

// initHTTPServerLogger 将 HTTP Server 内置 logger 对齐 detail 配置，
// 避免 GF 默认写入 {Y-m-d}.log（无前缀文件）。
func initHTTPServerLogger() {
	ctx := context.Background()
	if g.Cfg().MustGet(ctx, "logger.detail").IsEmpty() {
		return
	}
	detailCfg := g.Log("detail").GetConfig()
	detailCfg.StdoutPrint = g.Cfg().MustGet(ctx, "server.logStdout").Bool()
	_ = httpServer.Logger().SetConfig(detailCfg)
}
