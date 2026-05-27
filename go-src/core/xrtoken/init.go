package xrtoken

import (
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/guid"
	"time"
	"xr-game-server/core/cfg"
	"xr-game-server/core/event"
)

const (
	Time = 7 * 24 * time.Hour
)

var appCache = gcache.New()
var cmsCache = gcache.New()

func AddAppToken(authId uint64) string {
	token := guid.S()
	appCache.Set(gctx.New(), authId, token, Time)
	event.Pub(event.AppToken, &event.AppTokenData{Token: token, Id: authId})
	return token
}

func InitAppToken(authId uint64, token string, val time.Time) {
	expire := val.Sub(time.Now())
	appCache.Set(gctx.New(), authId, token, expire)
}

func AddCmsToken(authId uint64) string {
	token := guid.S()
	cmsCache.Set(gctx.New(), authId, token, Time)
	return token
}

func HasAppToken(authId uint64, token string) bool {
	ctx := gctx.New()
	if cfg.GetAuthCfg().LoginOff {
		return true
	}
	cacheToken, e := appCache.Get(gctx.New(), authId)
	if e == nil && cacheToken.String() == token {
		appCache.UpdateExpire(ctx, authId, Time)
		event.Pub(event.AppToken, &event.AppTokenData{Token: token, Id: authId})
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
