package minuteactivity

import (
	"context"
	"fmt"
	"sync"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/dto/trackingdto"
	"xr-game-server/gameevent"
)

const dashboardUserEventsCacheSuffix = "userEvents"

var (
	dashboardUserEventsTrendOnce sync.Once
	dashboardUserEventsTrendMgr  *cache.RowCache[*trackingdto.CMSMinuteLiveActivityTrendRes]
)

func dashboardUserEventsTrendCache() *cache.RowCache[*trackingdto.CMSMinuteLiveActivityTrendRes] {
	dashboardUserEventsTrendOnce.Do(func() {
		dashboardUserEventsTrendMgr = cache.NewRowCacheWithTTL[*trackingdto.CMSMinuteLiveActivityTrendRes](cache.DashboardCacheTTL)
	})
	return dashboardUserEventsTrendMgr
}

func isDefaultMinuteLiveActivityTrendReq(req *trackingdto.CMSMinuteLiveActivityTrendReq) bool {
	if req == nil {
		return true
	}
	return req.StartTime == 0 && req.EndTime == 0
}

func dashboardUserEventsTrendCacheKey(req *trackingdto.CMSMinuteLiveActivityTrendReq) string {
	platform := ""
	if req != nil {
		platform = gameevent.NormalizeStatClientPlatform(req.Platform)
	}
	if platform == "" {
		return cache.DashboardCMSCacheBucketKey(dashboardUserEventsCacheSuffix)
	}
	return cache.DashboardCMSCacheBucketKey(fmt.Sprintf("%s:%s", dashboardUserEventsCacheSuffix, platform))
}

func getCachedDashboardUserEventsTrend(req *trackingdto.CMSMinuteLiveActivityTrendReq, bypass bool, build func() (*trackingdto.CMSMinuteLiveActivityTrendRes, error)) (*trackingdto.CMSMinuteLiveActivityTrendRes, error) {
	mgr := dashboardUserEventsTrendCache()
	key := dashboardUserEventsTrendCacheKey(req)
	if bypass {
		v, err := build()
		if err != nil {
			return nil, err
		}
		if mgr != nil {
			_ = mgr.SetRow(gctx.New(), key, v)
		}
		return v, nil
	}
	if mgr == nil {
		return build()
	}
	res := mgr.MustGetRow(gctx.New(), key, func(ctx context.Context) (*trackingdto.CMSMinuteLiveActivityTrendRes, error) {
		return build()
	})
	return res, nil
}
