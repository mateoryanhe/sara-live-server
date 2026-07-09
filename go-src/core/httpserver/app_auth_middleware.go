package httpserver

import (
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"xr-game-server/core/xrtoken"
	"xr-game-server/errercode"
)

func MiddlewareAppAuth(r *ghttp.Request) {
	authStart := gtime.Now()
	token := r.GetHeader("Authorization", "")
	if token == "" || len(strings.Split(token, ".")) != 2 {
		logAPIRequestAuth(r, elapsedMs(authStart))
		WriteFailJson(r, errercode.Token)
		return
	}

	tokenStr := strings.Split(token, ".")[1]
	userId := strings.Split(token, ".")[0]
	if flag := xrtoken.HasAppToken(gconv.Uint64(userId), tokenStr); !flag {
		logAPIRequestAuth(r, elapsedMs(authStart))
		WriteFailJson(r, errercode.Token)
		return
	}
	logAPIRequestAuth(r, elapsedMs(authStart))
	r.Middleware.Next()
}
