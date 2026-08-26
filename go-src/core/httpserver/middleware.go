package httpserver

import (
	"net/http"
	"strings"
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
	applyRequestCORS(r)
	if r.Method == http.MethodOptions {
		if r.Response.Status == 0 {
			r.Response.WriteHeader(http.StatusNoContent)
		}
		return
	}
	r.Middleware.Next()
}

// applyRequestCORS 默认允许任意跨域,不做来源/头/方法限制.
func applyRequestCORS(r *ghttp.Request) {
	if r == nil {
		return
	}
	h := r.Response.Header()
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" {
		origin = "*"
	}
	h.Set("Access-Control-Allow-Origin", origin)
	h.Set("Access-Control-Allow-Credentials", "true")
	h.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,HEAD,OPTIONS")
	allowHeaders := strings.TrimSpace(r.Header.Get("Access-Control-Request-Headers"))
	if allowHeaders == "" {
		allowHeaders = "*"
	}
	h.Set("Access-Control-Allow-Headers", allowHeaders)
	h.Set("Access-Control-Expose-Headers", "*")
	h.Set("Access-Control-Max-Age", "86400")
}

func handleStaticCORS(r *ghttp.Request) bool {
	if r == nil {
		return false
	}
	applyRequestCORS(r)
	if r.Method == http.MethodOptions {
		r.Response.WriteHeader(http.StatusNoContent)
		r.ExitAll()
		return true
	}
	return false
}

// 请求入口:收到 header 后记录日志.
func middlewareLogReq(r *ghttp.Request) {
	if !shouldSkipAPILogChain(r) {
		logAPIRequestStart(r)
	}
	r.Middleware.Next()
}
