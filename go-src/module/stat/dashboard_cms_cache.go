package stat

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/dto/statdto"
)

const (
	dashboardSysStatCacheSuffix       = "sysStat"
	dashboardUserStatTrendCacheSuffix = "userTrend"
)

var (
	dashboardSysStatCacheOnce sync.Once
	dashboardSysStatCacheMgr  *cache.RowCache[*statdto.CMSSysStatRes]

	dashboardUserStatTrendCacheOnce sync.Once
	dashboardUserStatTrendCacheMgr  *cache.RowCache[*statdto.CMSUserStatTrendRes]
)

func dashboardSysStatCache() *cache.RowCache[*statdto.CMSSysStatRes] {
	dashboardSysStatCacheOnce.Do(func() {
		dashboardSysStatCacheMgr = cache.NewRowCacheWithTTL[*statdto.CMSSysStatRes](cache.DashboardCacheTTL)
	})
	return dashboardSysStatCacheMgr
}

func dashboardUserStatTrendCache() *cache.RowCache[*statdto.CMSUserStatTrendRes] {
	dashboardUserStatTrendCacheOnce.Do(func() {
		dashboardUserStatTrendCacheMgr = cache.NewRowCacheWithTTL[*statdto.CMSUserStatTrendRes](cache.DashboardCacheTTL)
	})
	return dashboardUserStatTrendCacheMgr
}

func getDashboardCachedRow[T any](mgr *cache.RowCache[T], cacheSuffix string, bypass bool, build func() (T, error)) (T, error) {
	key := cache.DashboardCMSCacheBucketKey(cacheSuffix)
	if bypass {
		v, err := build()
		if err != nil {
			var zero T
			return zero, err
		}
		if mgr != nil {
			_ = mgr.SetRow(gctx.New(), key, v)
		}
		return v, nil
	}
	if mgr == nil {
		return build()
	}
	res := mgr.MustGetRow(gctx.New(), key, func(ctx context.Context) (T, error) {
		return build()
	})
	return res, nil
}

func getCachedDashboardSysStat(bypass bool, build func() (*statdto.CMSSysStatRes, error)) (*statdto.CMSSysStatRes, error) {
	return getDashboardCachedRow(dashboardSysStatCache(), dashboardSysStatCacheSuffix, bypass, build)
}

func getCachedDashboardUserStatTrend(bypass bool, build func() (*statdto.CMSUserStatTrendRes, error)) (*statdto.CMSUserStatTrendRes, error) {
	return getDashboardCachedRow(dashboardUserStatTrendCache(), dashboardUserStatTrendCacheSuffix, bypass, build)
}
