package httpserver

import (
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
)

// cms应答结果
func customResponseMiddleware(r *ghttp.Request) {
	preHandlerMs := preHandlerDurationMs(r)
	handlerStart := time.Now()
	r.Middleware.Next()
	handlerMs := time.Since(handlerStart).Milliseconds()
	writeResponseAndLog(r, authIdFromRequest(r), preHandlerMs, handlerMs, nil)
}
