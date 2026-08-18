package cfgdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
)

const liveRevenueShareCfgCacheKey = "live_revenue_share_cfg"

var liveRevenueShareCfgCacheMgr *cache.CacheMgr

func InitLiveRevenueShareCfgDao() {
	liveRevenueShareCfgCacheMgr = cache.NewCacheMgr()
}

func loadLiveRevenueShareCfgFromDB() *entity.LiveRevenueShareCfg {
	var row entity.LiveRevenueShareCfg
	if err := g.DB().Model(string(entity.TbLiveRevenueShareCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveLiveRevenueShareCfg(row *entity.LiveRevenueShareCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbLiveRevenueShareCfg)).Save(row)
	return err
}

func ReloadLiveRevenueShareCfgCache() {
	if liveRevenueShareCfgCacheMgr == nil {
		return
	}
	liveRevenueShareCfgCacheMgr.FlushCache(liveRevenueShareCfgCacheKey, loadLiveRevenueShareCfgFromDB())
	liveRevenueShareCfgCacheMgr.Cache.UpdateExpire(gctx.New(), liveRevenueShareCfgCacheKey, time.Hour*24*365*100)
}

func GetLiveRevenueShareCfgCached() *entity.LiveRevenueShareCfg {
	if liveRevenueShareCfgCacheMgr == nil {
		return loadLiveRevenueShareCfgFromDB()
	}
	v := liveRevenueShareCfgCacheMgr.GetData(liveRevenueShareCfgCacheKey, func(ctx context.Context) (value interface{}, err error) {
		return loadLiveRevenueShareCfgFromDB(), nil
	})
	liveRevenueShareCfgCacheMgr.Cache.UpdateExpire(gctx.New(), liveRevenueShareCfgCacheKey, time.Hour*24*365*100)
	if v == nil {
		return nil
	}
	row, _ := v.(*entity.LiveRevenueShareCfg)
	return row
}
