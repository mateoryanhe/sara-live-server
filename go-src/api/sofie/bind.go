package sofie

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
)

func jsonHandler[Req any, Res any](fn func(context.Context, *Req) (*Res, error)) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		var req Req
		if err := r.Parse(&req); err != nil {
			r.SetError(err)
			return
		}
		res, err := fn(r.Context(), &req)
		if err != nil {
			r.SetError(err)
			return
		}
		httpserver.SetHandlerResponseData(r, res)
	}
}
