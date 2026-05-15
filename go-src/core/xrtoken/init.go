package xrtoken

import (
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/guid"
	"time"
	"xr-game-server/core/cfg"
)

const (
	Time = 10 * time.Hour
)

var appCache = gcache.New()
var cmsCache = gcache.New()

func AddAppToken(authId uint64) string {
	token := guid.S()
	appCache.Set(gctx.New(), authId, token, Time)
	return token
}

func AddCmsToken(authId uint64) string {
	token := guid.S()
	cmsCache.Set(gctx.New(), authId, token, Time)
	return token
}

func HasAppToken(authId uint64) bool {
	ctx := gctx.New()
	if cfg.GetAuthCfg().LoginOff {
		return true
	}
	if flag, _ := appCache.Contains(ctx, authId); flag {
		appCache.UpdateExpire(ctx, authId, Time)
		return true
	}
	return false
}

func HasCmsToken(authId uint64) bool {
	ctx := gctx.New()
	if cfg.GetAuthCfg().LoginOff {
		return true
	}
	if flag, _ := cmsCache.Contains(ctx, authId); flag {
		cmsCache.UpdateExpire(ctx, authId, Time)
		return true
	}
	return false
}
