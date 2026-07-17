package resourcemetricdto

import "github.com/gogf/gf/v2/frame/g"

type CMSResourceMetricTrendReq struct {
	g.Meta `path:"/getResourceMetricTrend" method:"post" summary:"获取系统资源趋势" tags:"资源监控"`
}

type CMSResourceMetricPoint struct {
	Time                string  `json:"time"`
	ProcMemMb           float64 `json:"procMemMb"`
	ProcHeapAllocMb     float64 `json:"procHeapAllocMb"`
	ProcHeapInuseMb     float64 `json:"procHeapInuseMb"`
	ProcHeapSysMb       float64 `json:"procHeapSysMb"`
	ProcHeapUsedPercent float64 `json:"procHeapUsedPercent"`
	ProcHeapIdlePercent float64 `json:"procHeapIdlePercent"`
	ProcCpuPercent      float64 `json:"procCpuPercent"`
	SysMemUsedMb        float64 `json:"sysMemUsedMb"`
	SysMemTotalMb       float64 `json:"sysMemTotalMb"`
	SysMemUsedPercent   float64 `json:"sysMemUsedPercent"`
	SysCpuPercent       float64 `json:"sysCpuPercent"`
	OnlineCount         uint64  `json:"onlineCount,string"`
}

type CMSResourceMetricTrendRes struct {
	Points []*CMSResourceMetricPoint `json:"points"`
}
