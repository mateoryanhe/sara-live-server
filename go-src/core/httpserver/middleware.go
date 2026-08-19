package httpserver

import (
	"sync/atomic"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/xrlog"
	"xr-game-server/errercode"
)

// 接受请求
var canDo = false

// restartClosing 热重启执行阶段:旧进程收到请求直接关闭连接.
var restartClosing atomic.Bool

// RejectNewRequests 拒绝新 HTTP 请求(SIGTERM 等普通关机).
func RejectNewRequests() {
	canDo = false
}

// EnterRestartClosingPhase RestartAllServer 完成后,旧进程不再接收新 HTTP/WS.
func EnterRestartClosingPhase() {
	restartClosing.Store(true)
	canDo = false
	xrlog.DetailLog.Warning(gctx.New(), "热重启:旧进程 Restart 完成,拒绝新 HTTP/WS 连接")
}

func middlewareRestartGuard(r *ghttp.Request) {
	if restartClosing.Load() {
		xrlog.DetailLog.Warningf(r.Context(), "热重启:旧进程收到新 HTTP 请求并关闭连接,URI=%s ip=%s", r.RequestURI, r.GetClientIp())
		r.Response.Header().Set("Connection", "close")
		r.Exit()
		return
	}
	r.Middleware.Next()
}

func middlewareCORS(r *ghttp.Request) {
	options := r.Response.DefaultCORSOptions()
	// CMS/App 自定义鉴权头
	for _, header := range []string{"token", "authId", "reqId", PackageNameHeader} {
		if options.AllowHeaders == "" {
			options.AllowHeaders = header
			continue
		}
		options.AllowHeaders += "," + header
	}
	r.Response.CORS(options)
	r.Middleware.Next()
}

// 请求入口:收到 header 后记录日志,关机时拒绝新请求
func middlewareLogReq(r *ghttp.Request) {
	if !shouldSkipAPILogChain(r) {
		logAPIRequestStart(r)
	}
	if !canDo {
		WriteFailJson(r, int(errercode.ServerClose))
		return
	}
	r.Middleware.Next()
}

func Ready() {
	canDo = true
}
