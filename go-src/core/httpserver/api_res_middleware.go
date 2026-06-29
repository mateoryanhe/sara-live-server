package httpserver

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

func apiResponseMiddleware(r *ghttp.Request) {
	r.Middleware.Next()
	writeResponseAndLog(r, authIdFromToken(r), func(res any) any {
		return CreateSuccess(res)
	})
}
