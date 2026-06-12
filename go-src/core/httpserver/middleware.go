package httpserver

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"os"
	"strconv"
	"strings"
	"time"
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
	//收到关机指令,服务器不再接受请求
	if !canDo {
		WriteFailJson(r, int(errercode.ServerClose))
		return
	}
	userId := ""
	token := r.GetHeader("Authorization", "")
	if token != "" {
		userId = strings.Split(token, ".")[0]
	}
	//设置请求时间
	r.Header.Set(DoTime, strconv.FormatInt(time.Now().UnixMilli(), 10))
	g.Log().Infof(gctx.New(), "收到前端请求,url=%v,ip=%v,authId=%v,请求数据=%v", r.URL.RequestURI(), r.GetClientIp(), userId, requestBodyForLog(r))
	r.Middleware.Next()
}
