package resourcemonitor

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/dao/resourcemetricdao"
	"xr-game-server/entity/sys"
)

// rollupLastCompletedBucket 将刚结束的 5 分钟细采样桶聚合成一条粗采样
func rollupLastCompletedBucket() {
	now := time.Now()
	bucketEnd := now.Truncate(CoarseInterval)
	bucketStart := bucketEnd.Add(-CoarseInterval)
	rollupBucket(bucketStart, bucketEnd)
}

// backfillMissingAggs 启动时补写保留期内缺失的粗采样桶(含旧细采样)
func backfillMissingAggs() {
	now := time.Now()
	end := now.Truncate(CoarseInterval)
	start := end.Add(-CoarseRetention)
	for bucketStart := start; bucketStart.Before(end); bucketStart = bucketStart.Add(CoarseInterval) {
		rollupBucket(bucketStart, bucketStart.Add(CoarseInterval))
	}
}

func rollupBucket(bucketStart, bucketEnd time.Time) {
	if !bucketEnd.After(bucketStart) {
		return
	}
	if resourcemetricdao.HasAggAt(bucketStart) {
		return
	}
	rows := resourcemetricdao.ListFineInBucket(bucketStart, bucketEnd)
	if len(rows) == 0 {
		return
	}
	agg := averageFineRows(bucketStart, rows)
	if err := resourcemetricdao.InsertAgg(agg); err != nil {
		g.Log().Errorf(gctx.New(), "写入资源粗采样失败,bucket=%v,err=%v", bucketStart, err)
	}
}

func averageFineRows(bucketStart time.Time, rows []*entity.SysResourceMetric) *entity.SysResourceMetricAgg {
	n := float64(len(rows))
	var (
		procMemMb, procHeapAllocMb, procHeapInuseMb, procHeapSysMb        float64
		procHeapUsedPercent, procHeapIdlePercent, procCpuPercent          float64
		sysMemUsedMb, sysMemTotalMb, sysMemUsedPercent, sysCpuPercent float64
		onlineMax                                                         uint64
	)
	for _, row := range rows {
		if row == nil {
			continue
		}
		procMemMb += row.ProcMemMb
		procHeapAllocMb += row.ProcHeapAllocMb
		procHeapInuseMb += row.ProcHeapInuseMb
		procHeapSysMb += row.ProcHeapSysMb
		procHeapUsedPercent += row.ProcHeapUsedPercent
		procHeapIdlePercent += row.ProcHeapIdlePercent
		procCpuPercent += row.ProcCpuPercent
		sysMemUsedMb += row.SysMemUsedMb
		sysMemTotalMb += row.SysMemTotalMb
		sysMemUsedPercent += row.SysMemUsedPercent
		sysCpuPercent += row.SysCpuPercent
		if row.OnlineCount > onlineMax {
			onlineMax = row.OnlineCount
		}
	}
	return entity.NewSysResourceMetricAggFromAvg(
		bucketStart,
		procMemMb/n, procHeapAllocMb/n, procHeapInuseMb/n, procHeapSysMb/n,
		procHeapUsedPercent/n, procHeapIdlePercent/n, procCpuPercent/n,
		sysMemUsedMb/n, sysMemTotalMb/n, sysMemUsedPercent/n, sysCpuPercent/n,
		onlineMax,
	)
}
