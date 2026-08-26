package httpserver

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"xr-game-server/core/hotrestart"
	"xr-game-server/core/startup"
	"xr-game-server/core/xrlog"
)

const (
	Ws                = "/ws"
	Token             = "token"
	AuthId            = "authId"
	PackageNameHeader = "X-App-Package"
	contentTypeJson   = "application/json"
	ReqId             = "reqId"
)

var httpServer = g.Server()

// InitHttpServer 初始化http服务器
func InitHttpServer() {
	hotrestart.RegisterEnterRestartPhase(EnterRestartClosingPhase)
	setupDomainSites()
	setupStaticPaths()
	initHTTPServerLogger()
	httpServer.SetErrorStack(true)
	httpServer.Use(middlewareCORS, middlewareRestartGuard)
	if g.Cfg().MustGet(context.Background(), "server.gzipEnabled").Bool() {
		httpServer.Use(ghttp.MiddlewareGzip)
	}
	httpServer.BindHookHandler("/*", ghttp.HookAfterOutput, hookAccessLogAfterOutput)
	httpServer.BindHookHandler("/*", ghttp.HookAfterOutput, hookHTTPErrorLogAfterOutput)
	httpServer.BindHookHandler("/*", ghttp.HookAfterOutput, hookAPIRequestAfterOutput)
	setupAppOpenApiHook()
	enableHotRestartGraceful()
	go waitHTTPServerReadyAndLogStartupEnd()
	httpServer.Run()
	hotrestart.NotifyOldProcessExit()
}

func waitHTTPServerReadyAndLogStartupEnd() {
	ctx := gctx.New()
	for httpServer.Status() != ghttp.ServerStatusRunning {
		time.Sleep(5 * time.Millisecond)
	}
	startup.LogEnd(ctx)
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
	xrlog.DetailLog.Infof(r.Context(), "beforeServeHook [is file:%v] URI:%s ip:%s", r.IsFileRequest(), r.RequestURI, r.GetClientIp())
	applyRequestCORS(r)
}

func GetReqId(ctx context.Context) uint64 {
	tokenStr := g.RequestFromCtx(ctx).GetHeader(ReqId, "")
	if tokenStr == "" {
		return 0
	}
	return gconv.Uint64(tokenStr)
}

// GetPackageNameFromContext 从请求头读取 App 包名
func GetPackageNameFromContext(ctx context.Context) string {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return ""
	}
	return strings.TrimSpace(r.GetHeader(PackageNameHeader, ""))
}
