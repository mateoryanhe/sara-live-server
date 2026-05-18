package httpserver

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"xr-game-server/constants/common"
	"xr-game-server/core/cfg"
	"xr-game-server/core/xrtoken"
	"xr-game-server/errercode"
)

func MiddlewareAppAuth(r *ghttp.Request) {
	if cfg.GetAuthCfg().LoginOff {
		r.Middleware.Next()
		return
	}
	tokenStr := r.GetHeader(Token)
	userId := r.GetHeader(AuthId)
	if len(tokenStr) == common.Zero {
		WriteFailJson(r, errercode.EmptyToken)
		return
	}
	if len(userId) == common.Zero {
		WriteFailJson(r, errercode.EmptyUserId)
		return
	}
	if flag := xrtoken.HasAppToken(gconv.Uint64(userId)); !flag {
		WriteFailJson(r, errercode.Token)
		return
	}
	r.Middleware.Next()
}
