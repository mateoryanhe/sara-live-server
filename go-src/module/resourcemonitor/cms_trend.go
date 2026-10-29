package resourcemonitor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/dao/resourcemetricdao"
	"xr-game-server/dto/resourcemetricdto"
	"xr-game-server/entity/sys"
)

// cmsTrendCacheTTL 与仪表盘 CMS 固定 5 分钟刷新对齐; 默认查询先读内存, 过期再查库.
const cmsTrendCacheTTL = cache.DashboardCacheTTL

var cmsTrendPointsCacheMgr *cache.RowCache[[]*trendPoint]

func cmsTrendPointsCache() *cache.RowCache[[]*trendPoint] {
	if cmsTrendPointsCacheMgr == nil {
		cmsTrendPointsCacheMgr = cache.NewRowCacheWithTTL[[]*trendPoint](cmsTrendCacheTTL)
	}
	return cmsTrendPointsCacheMgr
}

const trendTimeLayout = "2006-01-02 15:04:05"

type trendQuery struct {
	StartTime            string
	EndTime              string
	Limit                int
	BypassDashboardCache bool
}

// trendPoint 细采样 CMS 映射
type trendPoint struct {
	RecordedAt          time.Time
	ProcMemMb           float64
	ProcHeapAllocMb     float64
	ProcHeapInuseMb     float64
	ProcHeapSysMb       float64
	ProcHeapUsedPercent float64
	ProcHeapIdlePercent float64
	ProcCpuPercent      float64
	SysMemUsedMb        float64
	SysMemTotalMb       float64
	SysMemUsedPercent   float64
	SysCpuPercent       float64
	OnlineCount         uint64
	LiveAnchorCount     uint64
	GameOnlineCount     uint64
	GameBetUserCount    uint64
}

func GetCMSResourceMetricMemoryTrend(_ context.Context, req *resourcemetricdto.CMSResourceMetricMemoryTrendReq) (*resourcemetricdto.CMSResourceMetricMemoryTrendRes, error) {
	rows := listTrendPoints(toTrendQuery(req.StartTime, req.EndTime, req.Limit, req.BypassDashboardCache))
	points := make([]*resourcemetricdto.CMSResourceMetricMemoryPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		points = append(points, &resourcemetricdto.CMSResourceMetricMemoryPoint{
			Time:              row.RecordedAt.Format(trendTimeLayout),
			ProcMemMb:         row.ProcMemMb,
			SysMemUsedMb:      row.SysMemUsedMb,
			SysMemTotalMb:     row.SysMemTotalMb,
			SysMemUsedPercent: row.SysMemUsedPercent,
		})
	}
	return &resourcemetricdto.CMSResourceMetricMemoryTrendRes{Points: points}, nil
}

func GetCMSResourceMetricHeapTrend(_ context.Context, req *resourcemetricdto.CMSResourceMetricHeapTrendReq) (*resourcemetricdto.CMSResourceMetricHeapTrendRes, error) {
	rows := listTrendPoints(toTrendQuery(req.StartTime, req.EndTime, req.Limit, false))
	points := make([]*resourcemetricdto.CMSResourceMetricHeapPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		points = append(points, &resourcemetricdto.CMSResourceMetricHeapPoint{
			Time:            row.RecordedAt.Format(trendTimeLayout),
			ProcHeapAllocMb: row.ProcHeapAllocMb,
			ProcHeapInuseMb: row.ProcHeapInuseMb,
			ProcHeapSysMb:   row.ProcHeapSysMb,
		})
	}
	return &resourcemetricdto.CMSResourceMetricHeapTrendRes{Points: points}, nil
}

func GetCMSResourceMetricRatioTrend(_ context.Context, req *resourcemetricdto.CMSResourceMetricRatioTrendReq) (*resourcemetricdto.CMSResourceMetricRatioTrendRes, error) {
	rows := listTrendPoints(toTrendQuery(req.StartTime, req.EndTime, req.Limit, false))
	points := make([]*resourcemetricdto.CMSResourceMetricRatioPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		points = append(points, &resourcemetricdto.CMSResourceMetricRatioPoint{
			Time:                row.RecordedAt.Format(trendTimeLayout),
			ProcHeapUsedPercent: row.ProcHeapUsedPercent,
			ProcHeapIdlePercent: row.ProcHeapIdlePercent,
		})
	}
	return &resourcemetricdto.CMSResourceMetricRatioTrendRes{Points: points}, nil
}

func GetCMSResourceMetricCpuTrend(_ context.Context, req *resourcemetricdto.CMSResourceMetricCpuTrendReq) (*resourcemetricdto.CMSResourceMetricCpuTrendRes, error) {
	rows := listTrendPoints(toTrendQuery(req.StartTime, req.EndTime, req.Limit, false))
	points := make([]*resourcemetricdto.CMSResourceMetricCpuPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		points = append(points, &resourcemetricdto.CMSResourceMetricCpuPoint{
			Time:           row.RecordedAt.Format(trendTimeLayout),
			ProcCpuPercent: row.ProcCpuPercent,
			SysCpuPercent:  row.SysCpuPercent,
		})
	}
	return &resourcemetricdto.CMSResourceMetricCpuTrendRes{Points: points}, nil
}

func GetCMSResourceMetricGameStatTrend(_ context.Context, req *resourcemetricdto.CMSResourceMetricGameStatTrendReq) (*resourcemetricdto.CMSResourceMetricGameStatTrendRes, error) {
	limit := resolveTrendLimit(req.Limit)
	rows := listTrendPoints(toTrendQuery(req.StartTime, req.EndTime, req.Limit, req.BypassDashboardCache))
	points := make([]*resourcemetricdto.CMSResourceMetricGameStatPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		points = append(points, &resourcemetricdto.CMSResourceMetricGameStatPoint{
			Time:        row.RecordedAt.Format(trendTimeLayout),
			Count:       row.GameBetUserCount,
			OnlineCount: row.GameOnlineCount,
		})
	}
	return &resourcemetricdto.CMSResourceMetricGameStatTrendRes{
		Points:         points,
		RangeTruncated: len(rows) >= limit,
		TrendMaxPoints: limit,
	}, nil
}

func GetCMSResourceMetricOnlineTrend(_ context.Context, req *resourcemetricdto.CMSResourceMetricOnlineTrendReq) (*resourcemetricdto.CMSResourceMetricOnlineTrendRes, error) {
	limit := resolveTrendLimit(req.Limit)
	rows := listTrendPoints(toTrendQuery(req.StartTime, req.EndTime, limit, req.BypassDashboardCache))
	points := make([]*resourcemetricdto.CMSResourceMetricOnlinePoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		points = append(points, &resourcemetricdto.CMSResourceMetricOnlinePoint{
			Time:            row.RecordedAt.Format(trendTimeLayout),
			OnlineCount:     row.OnlineCount,
			LiveAnchorCount: row.LiveAnchorCount,
			AudienceCount:   audienceOnlineCount(row.OnlineCount, row.LiveAnchorCount),
		})
	}
	return &resourcemetricdto.CMSResourceMetricOnlineTrendRes{
		Points:         points,
		RangeTruncated: len(rows) >= limit,
		TrendMaxPoints: limit,
	}, nil
}

func audienceOnlineCount(total, liveAnchors uint64) uint64 {
	if total <= liveAnchors {
		return 0
	}
	return total - liveAnchors
}

func toTrendQuery(startTime, endTime string, limit int, bypassDashboardCache bool) trendQuery {
	return trendQuery{
		StartTime:            startTime,
		EndTime:              endTime,
		Limit:                limit,
		BypassDashboardCache: bypassDashboardCache,
	}
}

func listTrendPoints(q trendQuery) []*trendPoint {
	// 显式时间窗不做缓存:仪表盘曾每轮询带「当前秒」的滑动 start/end,
	// 若再按完整时间串做永久 key 会每分钟新增大约数 MB 的 26k 点切片.
	if strings.TrimSpace(q.StartTime) != "" || strings.TrimSpace(q.EndTime) != "" {
		return loadTrendPointsFromDB(q)
	}
	key := cmsTrendQueryCacheKey(q)
	mgr := cmsTrendPointsCache()
	if q.BypassDashboardCache {
		rows := loadTrendPointsFromDB(q)
		if mgr != nil {
			_ = mgr.SetRow(gctx.New(), key, rows)
		}
		return rows
	}
	return mgr.MustGetRow(gctx.New(), key, func(ctx context.Context) ([]*trendPoint, error) {
		return loadTrendPointsFromDB(q), nil
	})
}

func loadTrendPointsFromDB(q trendQuery) []*trendPoint {
	start, end := resolveTrendTimeRange(q)
	limit := resolveTrendLimit(q.Limit)
	return fineToTrendPoints(resourcemetricdao.ListFineByTimeRange(start, end, limit))
}

// cmsTrendQueryCacheKey 仅默认查询(无 start/end)使用;按 5 分钟分桶,配合 TTL.
func cmsTrendQueryCacheKey(q trendQuery) string {
	limit := resolveTrendLimit(q.Limit)
	bucket := time.Now().Unix() / int64(cache.DashboardCacheTTL/time.Second)
	return fmt.Sprintf("default:%d:fine:%d", limit, bucket)
}

func fineToTrendPoints(rows []*entity.SysResourceMetric) []*trendPoint {
	list := make([]*trendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &trendPoint{
			RecordedAt:          row.RecordedAt,
			ProcMemMb:           row.ProcMemMb,
			ProcHeapAllocMb:     row.ProcHeapAllocMb,
			ProcHeapInuseMb:     row.ProcHeapInuseMb,
			ProcHeapSysMb:       row.ProcHeapSysMb,
			ProcHeapUsedPercent: row.ProcHeapUsedPercent,
			ProcHeapIdlePercent: row.ProcHeapIdlePercent,
			ProcCpuPercent:      row.ProcCpuPercent,
			SysMemUsedMb:        row.SysMemUsedMb,
			SysMemTotalMb:       row.SysMemTotalMb,
			SysMemUsedPercent:   row.SysMemUsedPercent,
			SysCpuPercent:       row.SysCpuPercent,
			OnlineCount:         row.OnlineCount,
			LiveAnchorCount:     row.LiveAnchorCount,
			GameOnlineCount:     row.GameOnlineCount,
			GameBetUserCount:    row.GameBetUserCount,
		})
	}
	return list
}

func resolveTrendTimeRange(q trendQuery) (start, end time.Time) {
	now := time.Now()
	end = now
	start = now.Add(-FineRetention)
	if t, ok := parseTrendTime(q.EndTime); ok {
		end = t
	}
	if t, ok := parseTrendTime(q.StartTime); ok {
		start = t
	}
	if end.Before(start) {
		start, end = end, start
	}
	fineStart := now.Add(-FineRetention)
	if start.Before(fineStart) {
		start = fineStart
	}
	return start, end
}

func parseTrendTime(raw string) (time.Time, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, false
	}
	t, err := time.ParseInLocation(trendTimeLayout, raw, time.Local)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

func resolveTrendLimit(limit int) int {
	if limit <= 0 {
		return maxTrendPoints
	}
	if limit > maxTrendPoints {
		return maxTrendPoints
	}
	return limit
}
