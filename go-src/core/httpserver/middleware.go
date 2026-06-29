package httpserver

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"os"
	"strings"
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
	userId := ""
	token := r.GetHeader("Authorization", "")
	if token != "" {
		userId = strings.Split(token, ".")[0]
	}
	//收到关机指令,服务器不再接受请求
	if !canDo {
		g.Log().Infof(r.Context(), "关机了,收到前端请求,reqId=%v,url=%v,ip=%v,authId=%v,请求数据=%v", r.GetHeader(ReqId, ""), r.URL.RequestURI(), r.GetClientIp(), userId, requestBodyForLog(r))
		WriteFailJson(r, int(errercode.ServerClose))
		return
	}

	g.Log().Infof(r.Context(), "收到前端请求,reqId=%v,url=%v,ip=%v,authId=%v,请求数据=%v", r.GetHeader(ReqId, ""), r.URL.RequestURI(), r.GetClientIp(), userId, requestBodyForLog(r))
	r.Middleware.Next()
}
