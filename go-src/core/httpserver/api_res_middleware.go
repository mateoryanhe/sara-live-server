package httpserver

import (
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
)

func apiResponseMiddleware(r *ghttp.Request) {
	preHandlerMs := preHandlerDurationMs(r)
	handlerStart := time.Now()
	r.Middleware.Next()
	handlerMs := time.Since(handlerStart).Milliseconds()
	writeResponseAndLog(r, authIdFromRequest(r), preHandlerMs, handlerMs, func(res any) any {
		return CreateSuccess(res)
	})
}
