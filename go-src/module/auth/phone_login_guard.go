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

var phoneLoginCacheMgr *cache.RowCache[int64]

func initPhoneLoginGuard() {
	phoneLoginCacheMgr = cache.NewRowCache[int64]()
}

func checkPhoneLoginLimit(phoneKey string) error {
	ctx := gctx.New()
	if phoneLoginCacheMgr.ContainsRow(ctx, phoneLoginBlockKey(phoneKey)) {
		return errercode.CreateCode(errercode.RequestTooFrequent)
	}
	return nil
}

func markPhoneLoginFailure(phoneKey string) error {
	ctx := gctx.New()
	key := phoneLoginFailKey(phoneKey)
	count := 1

	current, ok := phoneLoginCacheMgr.GetRowCached(ctx, key)
	if ok && current > 0 {
		count = int(current + 1)
	}

	if count >= phoneLoginFailLimit {
		setPhoneLoginCache(phoneLoginBlockKey(phoneKey), time.Now().Unix(), phoneLoginFailExpire)
		phoneLoginCacheMgr.RemoveRow(ctx, key)
		return errercode.CreateCode(errercode.RequestTooFrequent)
	}

	setPhoneLoginCache(key, int64(count), phoneLoginFailExpire)
	return nil
}

func clearPhoneLoginFailure(phoneKey string) {
	ctx := gctx.New()
	phoneLoginCacheMgr.RemoveRow(ctx, phoneLoginFailKey(phoneKey))
	phoneLoginCacheMgr.RemoveRow(ctx, phoneLoginBlockKey(phoneKey))
}

func phoneLoginFailKey(phoneKey string) string {
	return "phone_login_fail:" + phoneKey
}

func phoneLoginBlockKey(phoneKey string) string {
	return "phone_login_block:" + phoneKey
}

func setPhoneLoginCache(key string, data int64, ttl time.Duration) {
	_ = phoneLoginCacheMgr.SetRow(gctx.New(), key, data, ttl)
}
