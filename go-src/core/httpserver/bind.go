package httpserver

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// RegAPI 绑定app控制器 检查权限 返回自定义格式数据
func RegAPI(prefix string, handlerOrObject ...interface{}) {
	httpServer.Group(prefix, func(group *ghttp.RouterGroup) {
		group.Middleware(middlewareLogReq, MiddlewareAppAuth, apiResponseMiddleware)
		group.Bind(handlerOrObject...)
	})
}

// RegNonAuthAPI 绑定无需权限检查控制器 返回自定义格式数据
func RegNonAuthAPI(prefix string, handlerOrObject ...interface{}) {
	httpServer.Group(prefix, func(group *ghttp.RouterGroup) {
		group.Middleware(middlewareLogReq, apiResponseMiddleware)
		group.Bind(handlerOrObject...)
	})
}

// RegAppCustomizeRes 自定义查询结果 需要鉴权 app专用
func RegAppCustomizeRes(prefix string, handlerOrObject ...interface{}) {
	httpServer.Group(prefix, func(group *ghttp.RouterGroup) {
		group.Middleware(middlewareLogReq, MiddlewareAppAuth, customResponseMiddleware)
		group.Bind(handlerOrObject...)
	})
}

// RegAPIHandler 绑定App原始Handler,不触发 GoFrame ParseMultipartForm(适合大文件流式上传)
func RegAPIHandler(prefix, pattern string, handler ghttp.HandlerFunc) {
	httpServer.Group(prefix, func(group *ghttp.RouterGroup) {
		group.Middleware(middlewareLogReq, MiddlewareAppAuth, apiResponseMiddleware)
		group.POST(pattern, handler)
	})
}

// RegCMS 绑定CMS控制器 需要鉴权
func RegCMS(prefix string, handlerOrObject ...interface{}) {
	httpServer.Group(prefix, func(group *ghttp.RouterGroup) {
		group.Middleware(middlewareLogReq, MiddlewareCmsAuth, apiResponseMiddleware)
		group.Bind(handlerOrObject...)
	})
}

// RegCMSHandler 绑定CMS原始Handler,不触发 GoFrame ParseMultipartForm
func RegCMSHandler(prefix, pattern string, handler ghttp.HandlerFunc) {
	httpServer.Group(prefix, func(group *ghttp.RouterGroup) {
		group.Middleware(middlewareLogReq, MiddlewareCmsAuth, apiResponseMiddleware)
		group.POST(pattern, handler)
	})
}

// RegCMSCustomizeRes 绑定cms自定义结果 需要鉴权
func RegCMSCustomizeRes(prefix string, handlerOrObject ...interface{}) {
	httpServer.Group(prefix, func(group *ghttp.RouterGroup) {
		group.Middleware(middlewareLogReq, MiddlewareCmsAuth, customResponseMiddleware)
		group.Bind(handlerOrObject...)
	})
}

// RegInternalHandler 绑定内部运维接口,不走 canDo 拦截.
func RegInternalHandler(pattern string, handler ghttp.HandlerFunc) {
	httpServer.BindHandler("GET:"+pattern, handler)
	httpServer.BindHandler("POST:"+pattern, handler)
}

// RegNonAuthHandler 绑定无需鉴权的原始 Handler(Pub/Sub 等 webhook)
func RegNonAuthHandler(prefix, pattern string, handler ghttp.HandlerFunc) {
	httpServer.Group(prefix, func(group *ghttp.RouterGroup) {
		group.Middleware(middlewareLogReq)
		group.POST(pattern, handler)
	})
}

// RegRootNonAuthCustomizeRes 根路径绑定,无需鉴权,自定义响应(第三方 webhook 等).
func RegRootNonAuthCustomizeRes(handlerOrObject ...interface{}) {
	httpServer.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(middlewareLogReq, customResponseMiddleware)
		group.Bind(handlerOrObject...)
	})
}

// RegCMSNonAuthCustomizeRes 自定义结果 不需鉴权
func RegCMSNonAuthCustomizeRes(prefix string, handlerOrObject ...interface{}) {
	httpServer.Group(prefix, func(group *ghttp.RouterGroup) {
		group.Middleware(middlewareLogReq, customResponseMiddleware)
		group.Bind(handlerOrObject...)
	})
}

// RegDataSyncReceive 数据同步接收端点(校验同步 Token,不走 CMS 登录)
func RegDataSyncReceive(prefix string, authMiddleware ghttp.HandlerFunc, handlerOrObject ...interface{}) {
	httpServer.Group(prefix, func(group *ghttp.RouterGroup) {
		group.Middleware(middlewareLogReq, authMiddleware, apiResponseMiddleware)
		group.Bind(handlerOrObject...)
	})
}
