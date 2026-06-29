package httpserver

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// cms应答结果
func customResponseMiddleware(r *ghttp.Request) {
	r.Middleware.Next()
	writeResponseAndLog(r, r.GetHeader(AuthId), nil)
}
