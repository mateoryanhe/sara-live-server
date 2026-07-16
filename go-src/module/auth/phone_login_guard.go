package auth

import (
	"time"
	"xr-game-server/core/cache"
	"xr-game-server/errercode"

	"github.com/gogf/gf/v2/os/gctx"
)

const (
	phoneLoginFailLimit  = 5
	phoneLoginFailExpire = 2 * time.Hour
)

var phoneLoginCacheMgr *cache.CacheMgr

func initPhoneLoginGuard() {
	phoneLoginCacheMgr = cache.NewCacheMgr()
}

func checkPhoneLoginLimit(phone string) error {
	ctx := gctx.New()
	if blocked, _ := phoneLoginCacheMgr.Cache.Contains(ctx, phoneLoginBlockKey(phone)); blocked {
		return errercode.CreateCode(errercode.RequestTooFrequent)
	}
	return nil
}

func markPhoneLoginFailure(phone string) error {
	ctx := gctx.New()
	key := phoneLoginFailKey(phone)
	count := 1

	val, _ := phoneLoginCacheMgr.Cache.Get(ctx, key)
	if val != nil {
		if current, ok := val.Val().(int); ok && current > 0 {
			count = current + 1
		}
	}

	if count >= phoneLoginFailLimit {
		setPhoneLoginCache(phoneLoginBlockKey(phone), time.Now().Unix(), phoneLoginFailExpire)
		_, _ = phoneLoginCacheMgr.Cache.Remove(ctx, key)
		return errercode.CreateCode(errercode.RequestTooFrequent)
	}

	setPhoneLoginCache(key, count, phoneLoginFailExpire)
	return nil
}

func clearPhoneLoginFailure(phone string) {
	ctx := gctx.New()
	_, _ = phoneLoginCacheMgr.Cache.Remove(ctx, phoneLoginFailKey(phone))
	_, _ = phoneLoginCacheMgr.Cache.Remove(ctx, phoneLoginBlockKey(phone))
}

func phoneLoginFailKey(phone string) string {
	return "phone_login_fail:" + phone
}

func phoneLoginBlockKey(phone string) string {
	return "phone_login_block:" + phone
}

func setPhoneLoginCache(key string, data any, ttl time.Duration) {
	_ = phoneLoginCacheMgr.Cache.Set(gctx.New(), key, data, ttl)
}
