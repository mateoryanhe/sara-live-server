package resourcemonitor

import (
	"context"
	"time"
	"xr-game-server/dao/resourcemetricdao"
	"xr-game-server/dto/resourcemetricdto"
	"xr-game-server/entity"
)

// GetCMSResourceMetricTrend CMS获取系统资源趋势(最近3天)
func GetCMSResourceMetricTrend(_ context.Context, _ *resourcemetricdto.CMSResourceMetricTrendReq) (*resourcemetricdto.CMSResourceMetricTrendRes, error) {
	since := time.Now().Add(-MetricRetention)
	rows := resourcemetricdao.ListSince(since)
	return &resourcemetricdto.CMSResourceMetricTrendRes{
		Points: toResourceMetricPoints(rows),
	}, nil
}

func toResourceMetricPoints(rows []*entity.SysResourceMetric) []*resourcemetricdto.CMSResourceMetricPoint {
	list := make([]*resourcemetricdto.CMSResourceMetricPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &resourcemetricdto.CMSResourceMetricPoint{
			Time:                row.RecordedAt.Format("01-02 15:04"),
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
		})
	}
	return list
}
