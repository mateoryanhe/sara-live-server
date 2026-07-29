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

// RegCMSNonAuthCustomizeRes 自定义结果 不需鉴权
func RegCMSNonAuthCustomizeRes(prefix string, handlerOrObject ...interface{}) {
	httpServer.Group(prefix, func(group *ghttp.RouterGroup) {
		group.Middleware(middlewareLogReq, customResponseMiddleware)
		group.Bind(handlerOrObject...)
	})
}
