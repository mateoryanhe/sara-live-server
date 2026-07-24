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
