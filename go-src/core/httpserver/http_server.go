package httpserver

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
	"xr-game-server/core/shutdown"
)

const (
	Ws              = "/ws"
	Token           = "token"
	AuthId          = "authId"
	DoTime          = "1ea75a33e32449c67dfd40fe8e23232d-doTime"
	LongDoTime      = 400
	contentTypeJson = "application/json"
)

var httpServer = g.Server()

// InitHttpServer 初始化http服务器
func InitHttpServer() {
	shutdown.RegCommonShutDownHandler(closeServer)
	setupDomainSites()
	if g.Cfg().MustGet(context.Background(), "server.gzipEnabled").Bool() {
		httpServer.Use(ghttp.MiddlewareGzip)
	}
	httpServer.BindHookHandler("/*", ghttp.HookBeforeServe, beforeServeHook)
	httpServer.Run()
}

func GetAuthId(ctx context.Context) uint64 {
	tokenStr := g.RequestFromCtx(ctx).GetHeader("Authorization", "")
	if tokenStr == "" {
		return 0
	}
	auth := strings.Split(tokenStr, ".")
	if 2 > len(auth) {
		return 0
	}
	return gconv.Uint64(auth[0])
}

func beforeServeHook(r *ghttp.Request) {
	g.Log().Debugf(r.GetCtx(), "beforeServeHook [is file:%v] URI:%s", r.IsFileRequest(), r.RequestURI)
	r.Response.CORSDefault()
}
