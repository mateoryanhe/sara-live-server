package httpserver

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"os"
	"xr-game-server/errercode"
)

// 接受请求
var canDo = true

func closeServer(sig os.Signal) {
	canDo = false
}

func middlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

// 记录前端请求日志
func middlewareLogReq(r *ghttp.Request) {
	authId := authIdFromToken(r)
	if !canDo {
		logAPIRequestStart(r, authId, "关机了,收到前端请求")
		WriteFailJson(r, int(errercode.ServerClose))
		return
	}

	logAPIRequestStart(r, authId, "收到前端请求")
	r.Middleware.Next()
}
