package xrtoken

import (
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/guid"
	"time"
	"xr-game-server/core/event"
)

const (
	Time = 30 * 24 * time.Hour
)

var appCache = gcache.New()
var cmsCache = gcache.New()

func AddAppToken(authId uint64) string {
	token := guid.S()
	appCache.Set(gctx.New(), authId, token, Time)
	event.Pub(event.AppToken, &event.AppTokenData{Token: token, Id: authId})
	return token
}

func DelToken(authId uint64) {
	appCache.Remove(gctx.New(), authId)
}

func InitAppToken(authId uint64, token string, val time.Time) {
	expire := val.Sub(time.Now())
	if expire <= 0 {
		return
	}
	appCache.Set(gctx.New(), authId, token, expire)
}

// ReloadAppTokenCache 清空并重新加载App Token内存缓存
func ReloadAppTokenCache(load func() []AppTokenCacheItem) {
	ctx := gctx.New()
	keys, _ := appCache.Keys(ctx)
	if len(keys) > 0 {
		_, _ = appCache.Remove(ctx, keys...)
	}
	if load == nil {
		return
	}
	for _, item := range load() {
		InitAppToken(item.UserId, item.Token, item.ExpireAt)
	}
}

type AppTokenCacheItem struct {
	UserId   uint64
	Token    string
	ExpireAt time.Time
}

func AddCmsToken(authId uint64) string {
	token := guid.S()
	cmsCache.Set(gctx.New(), authId, token, Time)
	event.Pub(event.CmsToken, &event.CmsTokenData{Token: token, Id: authId})
	return token
}

func InitCmsToken(authId uint64, token string, val time.Time) {
	expire := val.Sub(time.Now())
	if expire <= 0 {
		return
	}
	cmsCache.Set(gctx.New(), authId, token, expire)
}

func HasAppToken(authId uint64, token string) bool {
	ctx := gctx.New()

	cacheToken, e := appCache.Get(gctx.New(), authId)
	duration, _ := appCache.GetExpire(gctx.New(), authId)
	if e == nil && cacheToken.String() == token {
		if (7 * 24 * time.Hour) > duration {
			//2天快失效了,才刷新
			appCache.UpdateExpire(ctx, authId, Time)
			event.Pub(event.AppToken, &event.AppTokenData{Token: token, Id: authId})
		}
		return true
	}
	return false
}

func HasCmsToken(authId uint64, token string) bool {
	ctx := gctx.New()

	cacheToken, e := cmsCache.Get(gctx.New(), authId)
	duration, _ := cmsCache.GetExpire(gctx.New(), authId)
	if e == nil && cacheToken.String() == token {

		if (7 * 24 * time.Hour) > duration {
			//2天快失效了,才刷新
			cmsCache.UpdateExpire(ctx, authId, Time)
			event.Pub(event.CmsToken, &event.CmsTokenData{Token: token, Id: authId})
		}
		return true
	}
	return false
}
