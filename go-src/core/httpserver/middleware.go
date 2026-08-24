package httpserver

import (
	"sync/atomic"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/hotrestart"
)

// restartClosing 热重启第二阶段:旧进程仍可处理请求,但需记录日志.
var restartClosing atomic.Bool

// EnterRestartClosingPhase RestartAllServer 完成后进入热重启第二阶段.
func EnterRestartClosingPhase() {
	restartClosing.Store(true)
}

func middlewareRestartGuard(r *ghttp.Request) {
	if restartClosing.Load() {
		hotrestart.LogPhase2HTTPRequest(r.Context(), r.RequestURI, r.GetClientIp())
	}
	r.Middleware.Next()
}

func middlewareCORS(r *ghttp.Request) {
	options := r.Response.DefaultCORSOptions()
	// CMS/App 自定义鉴权头
	for _, header := range []string{"token", "authId", "reqId", PackageNameHeader, H5ClientHeader} {
		if options.AllowHeaders == "" {
			options.AllowHeaders = header
			continue
		}
		options.AllowHeaders += "," + header
	}
	r.Response.CORS(options)
	r.Middleware.Next()
}

// 请求入口:收到 header 后记录日志.
func middlewareLogReq(r *ghttp.Request) {
	if !shouldSkipAPILogChain(r) {
		logAPIRequestStart(r)
	}
	r.Middleware.Next()
}
