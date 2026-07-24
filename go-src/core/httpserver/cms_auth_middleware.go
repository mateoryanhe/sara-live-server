package httpserver

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"xr-game-server/constants/common"
	"xr-game-server/core/xrtoken"
	"xr-game-server/errercode"
)

func MiddlewareCmsAuth(r *ghttp.Request) {
	authStart := gtime.Now()
	skipLog := shouldSkipAPILogChain(r)
	tokenStr := r.GetHeader(Token)
	userId := r.GetHeader(AuthId)
	if len(tokenStr) == common.Zero {
		if !skipLog {
			logAPIRequestAuth(r, elapsedMs(authStart))
		}
		WriteFailJson(r, errercode.EmptyToken)
		return
	}
	if len(userId) == common.Zero {
		if !skipLog {
			logAPIRequestAuth(r, elapsedMs(authStart))
		}
		WriteFailJson(r, errercode.EmptyUserId)
		return
	}
	if flag := xrtoken.HasCmsToken(gconv.Uint64(userId), tokenStr); !flag {
		if !skipLog {
			logAPIRequestAuth(r, elapsedMs(authStart))
		}
		WriteFailJson(r, errercode.Token)
		return
	}
	if !skipLog {
		logAPIRequestAuth(r, elapsedMs(authStart))
	}
	r.Middleware.Next()
}
