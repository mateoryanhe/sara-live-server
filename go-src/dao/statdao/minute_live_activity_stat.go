package statdao

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/gameevent"
	statentity "xr-game-server/entity/stat"
)

var minuteLiveActivityStatCacheMgr *cache.RowCache[*statentity.MinuteLiveActivityStat]

func initMinuteLiveActivityStatDao() {
	minuteLiveActivityStatCacheMgr = cache.NewRowCache[*statentity.MinuteLiveActivityStat]()
}

func normalizeMinuteLiveActivityPlatform(platform string) string {
	return gameevent.NormalizeStatClientPlatform(platform)
}

func minuteLiveActivityCacheKey(at time.Time, platform string) string {
	return fmt.Sprintf("%s|%s", strconv.FormatUint(statentity.MinuteLiveActivityStatID(at, platform), 10), platform)
}

func GetMinuteLiveActivityStatByTimePlatform(at time.Time, platform string) *statentity.MinuteLiveActivityStat {
	platform = normalizeMinuteLiveActivityPlatform(platform)
	key := minuteLiveActivityCacheKey(at, platform)
	return minuteLiveActivityStatCacheMgr.MustGetRow(gctx.New(), key, func(ctx context.Context) (*statentity.MinuteLiveActivityStat, error) {
		bucket := statentity.TruncateUTCMinute(at)
		var data *statentity.MinuteLiveActivityStat
		_ = g.Model(string(statentity.TbMinuteLiveActivityStat)).Unscoped().Where(g.Map{
			string(statentity.MinuteLiveActivityStatBucketAt): bucket,
			string(statentity.MinuteLiveActivityStatPlatform): platform,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return statentity.NewMinuteLiveActivityStat(at, platform), nil
	})
}

func PublishMinuteLiveActivityStat(at time.Time, platform string, data *statentity.MinuteLiveActivityStat) {
	if data == nil || minuteLiveActivityStatCacheMgr == nil {
		return
	}
	platform = normalizeMinuteLiveActivityPlatform(platform)
	minuteLiveActivityStatCacheMgr.PublishRow(gctx.New(), minuteLiveActivityCacheKey(at, platform), data)
}

func ListMinuteLiveActivityStats(start, end time.Time, limit int, platform string) []*statentity.MinuteLiveActivityStat {
	if limit <= 0 {
		limit = 4320
	}
	platform = normalizeMinuteLiveActivityPlatform(platform)
	m := g.Model(string(statentity.TbMinuteLiveActivityStat)).Ctx(gctx.New())
	if !start.IsZero() {
		m = m.Where("bucket_at >= ?", start.UTC())
	}
	if !end.IsZero() {
		m = m.Where("bucket_at <= ?", end.UTC())
	}
	if platform != "" {
		list := make([]*statentity.MinuteLiveActivityStat, 0)
		_ = m.Where(string(statentity.MinuteLiveActivityStatPlatform), platform).
			Order("bucket_at asc").Limit(limit).Scan(&list)
		return list
	}
	list := make([]*statentity.MinuteLiveActivityStat, 0)
	_ = m.Fields(`
bucket_at,
SUM(join_count) AS join_count,
SUM(leave_count) AS leave_count,
SUM(gift_count) AS gift_count,
SUM(chat_count) AS chat_count,
SUM(bet_count) AS bet_count,
SUM(call1v1_count) AS call1v1_count,
SUM(call1v1_connect_success_count) AS call1v1_connect_success_count,
SUM(call1v1_consume_count) AS call1v1_consume_count,
SUM(recharge_order_count) AS recharge_order_count
`).Group("bucket_at").Order("bucket_at asc").Limit(limit).Scan(&list)
	return list
}

func GetMinuteLiveActivityLiveBucketAggregate(at time.Time) *statentity.MinuteLiveActivityStat {
	bucket := statentity.TruncateUTCMinute(at)
	var row *statentity.MinuteLiveActivityStat
	_ = g.Model(string(statentity.TbMinuteLiveActivityStat)).Ctx(gctx.New()).
		Where(string(statentity.MinuteLiveActivityStatBucketAt), bucket).
		Fields(`
bucket_at,
SUM(join_count) AS join_count,
SUM(leave_count) AS leave_count,
SUM(gift_count) AS gift_count,
SUM(chat_count) AS chat_count,
SUM(bet_count) AS bet_count,
SUM(call1v1_count) AS call1v1_count,
SUM(call1v1_connect_success_count) AS call1v1_connect_success_count,
SUM(call1v1_consume_count) AS call1v1_consume_count,
SUM(recharge_order_count) AS recharge_order_count
`).Scan(&row)
	return row
}

func DeleteMinuteLiveActivityBefore(before time.Time) {
	if before.IsZero() {
		return
	}
	_, _ = g.Model(string(statentity.TbMinuteLiveActivityStat)).Ctx(gctx.New()).
		Where("bucket_at < ?", before.UTC()).
		Delete()
}

// GetMinuteLiveActivityStatByTime 兼容旧调用：未知平台桶(platform 为空).
func GetMinuteLiveActivityStatByTime(at time.Time) *statentity.MinuteLiveActivityStat {
	return GetMinuteLiveActivityStatByTimePlatform(at, "")
}

func PublishMinuteLiveActivityStatLegacy(at time.Time, data *statentity.MinuteLiveActivityStat) {
	PublishMinuteLiveActivityStat(at, "", data)
}
