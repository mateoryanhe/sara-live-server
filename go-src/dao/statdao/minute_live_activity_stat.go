package statdao

import (
	"context"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/stat"
)

var minuteLiveActivityStatCacheMgr *cache.RowCache[*entity.MinuteLiveActivityStat]

func initMinuteLiveActivityStatDao() {
	minuteLiveActivityStatCacheMgr = cache.NewRowCache[*entity.MinuteLiveActivityStat]()
}

func minuteLiveActivityCacheKey(at time.Time) string {
	return strconv.FormatUint(entity.MinuteLiveActivityStatID(at), 10)
}

func GetMinuteLiveActivityStatByTime(at time.Time) *entity.MinuteLiveActivityStat {
	key := minuteLiveActivityCacheKey(at)
	return minuteLiveActivityStatCacheMgr.MustGetRow(gctx.New(), key, func(ctx context.Context) (*entity.MinuteLiveActivityStat, error) {
		id := entity.MinuteLiveActivityStatID(at)
		var data *entity.MinuteLiveActivityStat
		_ = g.Model(string(entity.TbMinuteLiveActivityStat)).Unscoped().Where(g.Map{
			string(db.IdName): id,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewMinuteLiveActivityStat(at), nil
	})
}

func PublishMinuteLiveActivityStat(at time.Time, data *entity.MinuteLiveActivityStat) {
	if data == nil || minuteLiveActivityStatCacheMgr == nil {
		return
	}
	minuteLiveActivityStatCacheMgr.PublishRow(gctx.New(), minuteLiveActivityCacheKey(at), data)
}

func ListMinuteLiveActivityStats(start, end time.Time, limit int) []*entity.MinuteLiveActivityStat {
	if limit <= 0 {
		limit = 4320
	}
	list := make([]*entity.MinuteLiveActivityStat, 0)
	m := g.Model(string(entity.TbMinuteLiveActivityStat)).Ctx(gctx.New())
	if !start.IsZero() {
		m = m.Where("bucket_at >= ?", start.UTC())
	}
	if !end.IsZero() {
		m = m.Where("bucket_at <= ?", end.UTC())
	}
	_ = m.Order("bucket_at asc").Limit(limit).Scan(&list)
	return list
}

func DeleteMinuteLiveActivityBefore(before time.Time) {
	if before.IsZero() {
		return
	}
	_, _ = g.Model(string(entity.TbMinuteLiveActivityStat)).Ctx(gctx.New()).
		Where("bucket_at < ?", before.UTC()).
		Delete()
}
