package httpserver

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"xr-game-server/core/shutdown"
)

const (
	Ws              = "/ws"
	Token           = "token"
	AuthId          = "authId"
	contentTypeJson = "application/json"
	ReqId           = "reqId"
)

var httpServer = g.Server()

// InitHttpServer 初始化http服务器
func InitHttpServer() {
	shutdown.RegCommonShutDownHandler(closeServer)
	setupDomainSites()
	bindCMSStaticFallback(context.Background())
	httpServer.SetErrorStack(true)
	httpServer.Use(middlewareCORS)
	if g.Cfg().MustGet(context.Background(), "server.gzipEnabled").Bool() {
		httpServer.Use(ghttp.MiddlewareGzip)
	}
	httpServer.BindHookHandler("/*", ghttp.HookAfterOutput, hookAPIRequestAfterOutput)
	httpServer.Run()
}

func GetAuthId(ctx context.Context) uint64 {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return 0
	}
	return gconv.Uint64(authIdFromRequest(r))
}

func authIdFromToken(r *ghttp.Request) string {
	token := r.GetHeader("Authorization", "")
	if token == "" {
		return ""
	}
	return strings.Split(token, ".")[0]
}

// authIdFromRequest App 从 Authorization 解析,CMS 从 authId header 读取
func authIdFromRequest(r *ghttp.Request) string {
	if id := authIdFromToken(r); id != "" {
		return id
	}
	return r.GetHeader(AuthId)
}

func beforeServeHook(r *ghttp.Request) {
	g.Log().Infof(r.Context(), "beforeServeHook [is file:%v] URI:%s ip:%s", r.IsFileRequest(), r.RequestURI, r.GetClientIp())
	r.Response.CORSDefault()
}

func GetReqId(ctx context.Context) uint64 {
	tokenStr := g.RequestFromCtx(ctx).GetHeader(ReqId, "")
	if tokenStr == "" {
		return 0
	}
	return gconv.Uint64(tokenStr)
}
