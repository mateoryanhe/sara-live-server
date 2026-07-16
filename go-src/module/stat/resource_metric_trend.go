package stat

import (
	"context"
	"time"
	"xr-game-server/dao/statdao"
	"xr-game-server/dto/statdto"
	"xr-game-server/entity"
)

// GetCMSResourceMetricTrend CMS获取系统资源趋势(最近3天)
func GetCMSResourceMetricTrend(_ context.Context, _ *statdto.CMSResourceMetricTrendReq) (*statdto.CMSResourceMetricTrendRes, error) {
	since := time.Now().Add(-resourceMetricRetention)
	rows := statdao.ListSysResourceMetricsSince(since)
	return &statdto.CMSResourceMetricTrendRes{
		Points: toResourceMetricPoints(rows),
	}, nil
}

func toResourceMetricPoints(rows []*entity.SysResourceMetric) []*statdto.CMSResourceMetricPoint {
	list := make([]*statdto.CMSResourceMetricPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &statdto.CMSResourceMetricPoint{
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
