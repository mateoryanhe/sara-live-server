package httpserver

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"xr-game-server/constants/common"
	"xr-game-server/core/xrtoken"
	"xr-game-server/errercode"
)

func MiddlewareCmsAuth(r *ghttp.Request) {
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
	if flag := xrtoken.HasCmsToken(gconv.Uint64(userId), tokenStr); !flag {
		WriteFailJson(r, errercode.Token)
		return
	}
	r.Middleware.Next()
}
