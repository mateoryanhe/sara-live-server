package httpserver

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
	"xr-game-server/core/xrtoken"
	"xr-game-server/errercode"
)

func MiddlewareAppAuth(r *ghttp.Request) {
	token := r.GetHeader("Authorization", "")
	if token == "" || len(strings.Split(token, ".")) != 2 {
		WriteFailJson(r, errercode.Token)
		return
	}

	tokenStr := strings.Split(token, ".")[1]
	userId := strings.Split(token, ".")[0]
	if flag := xrtoken.HasAppToken(gconv.Uint64(userId), tokenStr); !flag {
		WriteFailJson(r, errercode.Token)
		return
	}
	r.Middleware.Next()
}
