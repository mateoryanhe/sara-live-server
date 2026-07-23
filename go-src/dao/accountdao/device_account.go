package accountdao

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var deviceAccountCacheMgr *cache.CacheMgr

const deviceAccountMissCacheTime = 10 * time.Minute

func initDeviceAccountDao() {
	deviceAccountCacheMgr = cache.NewCacheMgr()
}

func deviceAccountCacheKey(deviceId string, channel uint) string {
	return fmt.Sprintf("%s:%d", deviceId, channel)
}

func deviceAccountMissCacheKey(key string) string {
	return "miss:" + key
}

func clearDeviceAccountMissCache(key string) {
	ctx := gctx.New()
	_, _ = deviceAccountCacheMgr.Cache.Remove(ctx, deviceAccountMissCacheKey(key))
}

// FindDeviceAccount 只读查询设备账号,不存在时返回 nil,不会创建新记录;未命中结果负缓存 10 分钟
func FindDeviceAccount(deviceId string, channel uint) *entity.Account {
	deviceId = strings.TrimSpace(deviceId)
	if deviceId == "" || channel == 0 {
		return nil
	}

	key := deviceAccountCacheKey(deviceId, channel)
	ctx := gctx.New()

	if ok, _ := deviceAccountCacheMgr.Cache.Contains(ctx, deviceAccountMissCacheKey(key)); ok {
		return nil
	}

	if cached := deviceAccountCacheMgr.GetFromCache(key); cached != nil {
		return cached.(*entity.Account)
	}

	var account *entity.Account
	_ = g.Model(string(entity.TbAccount)).Unscoped().Where(g.Map{
		string(entity.AccountOpenId):  deviceId,
		string(entity.AccountChannel): channel,
	}).Scan(&account)
	if account == nil {
		_ = deviceAccountCacheMgr.Cache.Set(ctx, deviceAccountMissCacheKey(key), 1, deviceAccountMissCacheTime)
		return nil
	}

	deviceAccountCacheMgr.FlushCache(key, account)
	return account
}

// GetDeviceAccount 根据设备码+渠道获取账号,不存在则创建并异步入库
func GetDeviceAccount(deviceId string, channel uint) (*entity.Account, bool) {
	deviceId = strings.TrimSpace(deviceId)
	if deviceId == "" || channel == 0 {
		return nil, false
	}

	isNew := false
	key := deviceAccountCacheKey(deviceId, channel)
	cacheData := deviceAccountCacheMgr.GetData(key, func(ctx context.Context) (value interface{}, err error) {
		var account *entity.Account
		err = g.Model(string(entity.TbAccount)).Unscoped().Where(g.Map{
			string(entity.AccountOpenId):  deviceId,
			string(entity.AccountChannel): channel,
		}).Scan(&account)
		if account != nil {
			return account, nil
		}
		isNew = true
		acc := entity.NewAccount(deviceId, channel)
		acc.SetPhoneAreaCode("")
		return acc, nil
	})
	if cacheData == nil {
		return nil, false
	}
	clearDeviceAccountMissCache(key)
	return cacheData.(*entity.Account), isNew
}
