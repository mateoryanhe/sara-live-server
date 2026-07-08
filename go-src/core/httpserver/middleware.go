package httpserver

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"os"
	"xr-game-server/errercode"
)

// 接受请求
var canDo = false

func closeServer(sig os.Signal) {
	canDo = false
}

func middlewareCORS(r *ghttp.Request) {
	options := r.Response.DefaultCORSOptions()
	// CMS/App 自定义鉴权头
	for _, header := range []string{"token", "authId", "reqId"} {
		if options.AllowHeaders == "" {
			options.AllowHeaders = header
			continue
		}
		options.AllowHeaders += "," + header
	}
	r.Response.CORS(options)
	r.Middleware.Next()
}

// 记录前端请求日志
func middlewareLogReq(r *ghttp.Request) {
	stashRequestTiming(r)
	authId := authIdFromRequest(r)
	if !canDo {
		logAPIRequestStart(r, authId, "关机了,收到前端请求")
		WriteFailJson(r, int(errercode.ServerClose))
		return
	}

	logAPIRequestStart(r, authId, "收到前端请求")
	r.Middleware.Next()
}

func Ready() {
	canDo = true
}
