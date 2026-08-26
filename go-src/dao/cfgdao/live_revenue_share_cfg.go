package cfgdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
)

const liveRevenueShareCfgCacheKey = "live_revenue_share_cfg"

var liveRevenueShareCfgCacheMgr *cache.RowCache[*entity.LiveRevenueShareCfg]

func InitLiveRevenueShareCfgDao() {
	liveRevenueShareCfgCacheMgr = cache.NewRowCache[*entity.LiveRevenueShareCfg]()
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
	liveRevenueShareCfgCacheMgr.PublishRow(gctx.New(), liveRevenueShareCfgCacheKey, loadLiveRevenueShareCfgFromDB())
}

func GetLiveRevenueShareCfgCached() *entity.LiveRevenueShareCfg {
	if liveRevenueShareCfgCacheMgr == nil {
		return loadLiveRevenueShareCfgFromDB()
	}
	v := liveRevenueShareCfgCacheMgr.MustGetRow(gctx.New(), liveRevenueShareCfgCacheKey, func(ctx context.Context) (*entity.LiveRevenueShareCfg, error) {
		return loadLiveRevenueShareCfgFromDB(), nil
	})
	return v
}
