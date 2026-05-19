package httpserver

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
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
	//添加静态路径
	shutdown.RegCommonShutDownHandler(closeServer)
	httpServer.BindHookHandler("/*", ghttp.HookBeforeServe, beforeServeHook)
	httpServer.Run()
}

func GetAuthId(ctx context.Context) uint64 {
	authId := g.RequestFromCtx(ctx).GetHeader(AuthId)
	return gconv.Uint64(authId)
}

func beforeServeHook(r *ghttp.Request) {
	g.Log().Debugf(r.GetCtx(), "beforeServeHook [is file:%v] URI:%s", r.IsFileRequest(), r.RequestURI)
	r.Response.CORSDefault()
}
