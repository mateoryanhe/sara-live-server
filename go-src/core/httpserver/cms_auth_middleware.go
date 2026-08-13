package httpserver

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"xr-game-server/constants/common"
	"xr-game-server/core/xrtoken"
	"xr-game-server/errercode"
)

// 公共上传接口,各业务页共用,登录即可访问
var cmsAuthSkipApiPaths = map[string]struct{}{
	"/upload/uploadFile": {},
}

var cmsApiPermissionChecker func(userId uint64, apiPath string) bool

// SetCmsApiPermissionChecker 注册 CMS 接口权限检查(由 cmsuserdao 初始化时注入,避免 import cycle)
func SetCmsApiPermissionChecker(fn func(userId uint64, apiPath string) bool) {
	cmsApiPermissionChecker = fn
}

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
	if !cmsAuthHasApiPermission(gconv.Uint64(userId), r.URL.Path) {
		if !skipLog {
			logAPIRequestAuth(r, elapsedMs(authStart))
		}
		WriteFailJson(r, int(errercode.NoPermission))
		return
	}
	if !skipLog {
		logAPIRequestAuth(r, elapsedMs(authStart))
	}
	r.Middleware.Next()
}

func cmsAuthHasApiPermission(userId uint64, apiPath string) bool {
	if apiPath == "" {
		return false
	}
	if _, ok := cmsAuthSkipApiPaths[apiPath]; ok {
		return true
	}
	if cmsApiPermissionChecker == nil {
		return true
	}
	return cmsApiPermissionChecker(userId, apiPath)
}
