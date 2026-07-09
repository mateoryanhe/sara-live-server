package httpserver

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

// cms应答结果
func customResponseMiddleware(r *ghttp.Request) {
	logRequestBodyBeforeHandler(r)
	handlerStart := gtime.Now()
	r.Middleware.Next()
	logAPIRequestHandler(r, elapsedMs(handlerStart))
	writeResponse(r, nil)
}
