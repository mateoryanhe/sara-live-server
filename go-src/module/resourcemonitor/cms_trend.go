package resourcemonitor

import (
	"context"
	"strings"
	"time"

	"xr-game-server/dao/resourcemetricdao"
	"xr-game-server/dto/resourcemetricdto"
	"xr-game-server/entity"
)

const trendTimeLayout = "2006-01-02 15:04:05"

type trendQuery struct {
	StartTime string
	EndTime   string
	Limit     int
}

func GetCMSResourceMetricMemoryTrend(_ context.Context, req *resourcemetricdto.CMSResourceMetricMemoryTrendReq) (*resourcemetricdto.CMSResourceMetricMemoryTrendRes, error) {
	rows := listTrendRows(toTrendQuery(req.StartTime, req.EndTime, req.Limit))
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
	rows := listTrendRows(toTrendQuery(req.StartTime, req.EndTime, req.Limit))
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
	rows := listTrendRows(toTrendQuery(req.StartTime, req.EndTime, req.Limit))
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
	rows := listTrendRows(toTrendQuery(req.StartTime, req.EndTime, req.Limit))
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
	rows := listTrendRows(toTrendQuery(req.StartTime, req.EndTime, req.Limit))
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

func listTrendRows(q trendQuery) []*entity.SysResourceMetric {
	start, end := resolveTrendTimeRange(q)
	limit := resolveTrendLimit(q.Limit)
	return resourcemetricdao.ListByTimeRange(start, end, limit)
}

func resolveTrendTimeRange(q trendQuery) (start, end time.Time) {
	now := time.Now()
	end = now
	start = now.Add(-MetricRetention)
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
		return 1000
	}
	if limit > 1000 {
		return 1000
	}
	return limit
}
