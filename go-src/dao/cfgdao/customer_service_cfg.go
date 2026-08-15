package cfgdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/user"
)

const customerServiceCfgCacheKey = "cfg"

var customerServiceCfgCacheMgr *cache.CacheMgr

func InitCustomerServiceCfgDao() {
	customerServiceCfgCacheMgr = cache.NewCacheMgr()
}

func loadCustomerServiceCfgFromDB() *entity.CustomerServiceCfg {
	var row entity.CustomerServiceCfg
	if err := g.DB().Model(string(entity.TbCustomerServiceCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveCustomerServiceCfg(row *entity.CustomerServiceCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbCustomerServiceCfg)).Save(row)
	return err
}

// ReloadCustomerServiceCfgCache 配置变更后刷新 gcache 并从数据库重新加载
func ReloadCustomerServiceCfgCache() {
	if customerServiceCfgCacheMgr == nil {
		return
	}
	customerServiceCfgCacheMgr.FlushCache(customerServiceCfgCacheKey, loadCustomerServiceCfgFromDB())
	customerServiceCfgCacheMgr.Cache.UpdateExpire(gctx.New(), customerServiceCfgCacheKey, time.Hour*24*365*100)
}

// GetCustomerServiceCfgCached 获取客服联系配置(优先读 gcache,未命中再查库)
func GetCustomerServiceCfgCached() *entity.CustomerServiceCfg {
	if customerServiceCfgCacheMgr == nil {
		return loadCustomerServiceCfgFromDB()
	}
	v := customerServiceCfgCacheMgr.GetData(customerServiceCfgCacheKey, func(ctx context.Context) (value interface{}, err error) {
		return loadCustomerServiceCfgFromDB(), nil
	})
	customerServiceCfgCacheMgr.Cache.UpdateExpire(gctx.New(), customerServiceCfgCacheKey, time.Hour*24*365*100)
	if v == nil {
		return nil
	}
	row, _ := v.(*entity.CustomerServiceCfg)
	return row
}
