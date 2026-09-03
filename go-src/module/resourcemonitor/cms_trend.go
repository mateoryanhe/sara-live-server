package resourcemonitor

import (
	"context"
	"strings"
	"time"

	"xr-game-server/dao/resourcemetricdao"
	"xr-game-server/dto/resourcemetricdto"
	"xr-game-server/entity/sys"
)

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

func GetCMSResourceMetricOnlineTrend(_ context.Context, req *resourcemetricdto.CMSResourceMetricOnlineTrendReq) (*resourcemetricdto.CMSResourceMetricOnlineTrendRes, error) {
	rows := listTrendPoints(toTrendQuery(req.StartTime, req.EndTime, req.Limit))
	points := make([]*resourcemetricdto.CMSResourceMetricOnlinePoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		points = append(points, &resourcemetricdto.CMSResourceMetricOnlinePoint{
			Time:        row.RecordedAt.Format(trendTimeLayout),
			OnlineCount: row.OnlineCount,
		})
	}
	return &resourcemetricdto.CMSResourceMetricOnlineTrendRes{Points: points}, nil
}

func toTrendQuery(startTime, endTime string, limit int) trendQuery {
	return trendQuery{StartTime: startTime, EndTime: endTime, Limit: limit}
}

func listTrendPoints(q trendQuery) []*trendPoint {
	start, end := resolveTrendTimeRange(q)
	limit := resolveTrendLimit(q.Limit)
	if useFineTrend(start, end) {
		return fineToTrendPoints(resourcemetricdao.ListFineByTimeRange(start, end, limit))
	}
	return aggToTrendPoints(resourcemetricdao.ListAggByTimeRange(start, end, limit))
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
