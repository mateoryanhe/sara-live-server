package resourcemonitor

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/dao/resourcemetricdao"
	"xr-game-server/dto/resourcemetricdto"
	"xr-game-server/entity/sys"
)

var cmsTrendPointsCacheMgr *cache.RowCache[[]*trendPoint]

func cmsTrendPointsCache() *cache.RowCache[[]*trendPoint] {
	if cmsTrendPointsCacheMgr == nil {
		cmsTrendPointsCacheMgr = cache.NewPermanentRowCache[[]*trendPoint]()
	}
	return cmsTrendPointsCacheMgr
}

const trendTimeLayout = "2006-01-02 15:04:05"

type trendQuery struct {
	StartTime string
	EndTime   string
	Limit     int
}

// trendPoint 细/粗采样统一视图,供 CMS 映射
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
	rows := listTrendPoints(toTrendQuery(req.StartTime, req.EndTime, req.Limit))
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
	rows := listTrendPoints(toTrendQuery(req.StartTime, req.EndTime, req.Limit))
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
	rows := listTrendPoints(toTrendQuery(req.StartTime, req.EndTime, req.Limit))
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
	rows := listTrendPoints(toTrendQuery(req.StartTime, req.EndTime, req.Limit))
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
	rows := listTrendPoints(toTrendQuery(req.StartTime, req.EndTime, req.Limit))
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
	rows := listTrendPoints(toTrendQuery(req.StartTime, req.EndTime, limit))
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

func toTrendQuery(startTime, endTime string, limit int) trendQuery {
	return trendQuery{StartTime: startTime, EndTime: endTime, Limit: limit}
}

func listTrendPoints(q trendQuery) []*trendPoint {
	key := cmsTrendQueryCacheKey(q)
	return cmsTrendPointsCache().MustGetRow(gctx.New(), key, func(ctx context.Context) ([]*trendPoint, error) {
		return loadTrendPointsFromDB(q), nil
	})
}

func loadTrendPointsFromDB(q trendQuery) []*trendPoint {
	start, end := resolveTrendTimeRange(q)
	limit := resolveTrendLimit(q.Limit)
	if useFineTrend(start, end) {
		return fineToTrendPoints(resourcemetricdao.ListFineByTimeRange(start, end, limit))
	}
	now := time.Now()
	fineStart := now.Add(-FineRetention)
	if end.After(fineStart) && start.Before(fineStart) {
		coarseEnd := end
		if coarseEnd.After(fineStart) {
			coarseEnd = fineStart
		}
		coarse := aggToTrendPoints(resourcemetricdao.ListAggByTimeRange(start, coarseEnd, limit))
		remaining := limit - len(coarse)
		if remaining < 1 {
			remaining = 1
		}
		fine := fineToTrendPoints(resourcemetricdao.ListFineByTimeRange(fineStart, end, remaining))
		return capTrendPoints(mergeTrendPointsByTime(coarse, fine), limit)
	}
	return aggToTrendPoints(resourcemetricdao.ListAggByTimeRange(start, end, limit))
}

func mergeTrendPointsByTime(parts ...[]*trendPoint) []*trendPoint {
	total := 0
	for _, p := range parts {
		total += len(p)
	}
	if total == 0 {
		return nil
	}
	merged := make([]*trendPoint, 0, total)
	for _, p := range parts {
		for _, row := range p {
			if row != nil {
				merged = append(merged, row)
			}
		}
	}
	sort.Slice(merged, func(i, j int) bool {
		return merged[i].RecordedAt.Before(merged[j].RecordedAt)
	})
	if len(merged) <= 1 {
		return merged
	}
	out := make([]*trendPoint, 0, len(merged))
	var lastAt time.Time
	for _, row := range merged {
		if len(out) > 0 && row.RecordedAt.Equal(lastAt) {
			out[len(out)-1] = row
			continue
		}
		out = append(out, row)
		lastAt = row.RecordedAt
	}
	return out
}

// capTrendPoints 超出 limit 时保留较早一段，与 DAO 升序 limit 策略一致
func capTrendPoints(rows []*trendPoint, limit int) []*trendPoint {
	if limit <= 0 || len(rows) <= limit {
		return rows
	}
	return rows[:limit]
}

// cmsTrendQueryCacheKey 仪表盘默认查询按 10 分钟分桶(与细采样间隔一致); 自定义时间范围按参数固定缓存.
func cmsTrendQueryCacheKey(q trendQuery) string {
	limit := resolveTrendLimit(q.Limit)
	if strings.TrimSpace(q.StartTime) != "" || strings.TrimSpace(q.EndTime) != "" {
		return fmt.Sprintf("custom:%d:%s:%s", limit, strings.TrimSpace(q.StartTime), strings.TrimSpace(q.EndTime))
	}
	start, end := resolveTrendTimeRange(q)
	mode := "agg"
	if useFineTrend(start, end) {
		mode = "fine"
	}
	bucket := time.Now().Unix() / (10 * 60)
	return fmt.Sprintf("default:%d:%s:%d", limit, mode, bucket)
}

// useFineTrend 查询窗口完全落在细采样保留期内时用细数据,否则用粗数据
func useFineTrend(start, end time.Time) bool {
	now := time.Now()
	fineStart := now.Add(-FineRetention)
	if start.Before(fineStart) {
		return false
	}
	if end.Before(fineStart) {
		return false
	}
	return true
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

func aggToTrendPoints(rows []*entity.SysResourceMetricAgg) []*trendPoint {
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
	start = now.Add(-CoarseRetention)
	if t, ok := parseTrendTime(q.EndTime); ok {
		end = t
	}
	if t, ok := parseTrendTime(q.StartTime); ok {
		start = t
	}
	if end.Before(start) {
		start, end = end, start
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
