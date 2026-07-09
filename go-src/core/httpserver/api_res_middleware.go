package httpserver

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

func apiResponseMiddleware(r *ghttp.Request) {
	logRequestBodyBeforeHandler(r)
	handlerStart := gtime.Now()
	r.Middleware.Next()
	logAPIRequestHandler(r, elapsedMs(handlerStart))
	writeResponse(r, func(res any) any {
		return CreateSuccess(res)
	})
}
