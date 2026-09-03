package resourcemetricdto

import "github.com/gogf/gf/v2/frame/g"

type CMSResourceMetricMemoryTrendReq struct {
	g.Meta    `path:"/getResourceMetricMemoryTrend" method:"post" summary:"获取内存资源趋势" tags:"资源监控"`
	StartTime string `json:"startTime" dc:"开始时间 YYYY-MM-DD HH:mm:ss,空则默认最近3天(粗采样)"`
	EndTime   string `json:"endTime" dc:"结束时间 YYYY-MM-DD HH:mm:ss,空则默认当前时间"`
	Limit     int    `json:"limit" dc:"最多返回条数,默认10000,最大10000"`
}

type CMSResourceMetricHeapTrendReq struct {
	g.Meta    `path:"/getResourceMetricHeapTrend" method:"post" summary:"获取堆内存资源趋势" tags:"资源监控"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Limit     int    `json:"limit"`
}

type CMSResourceMetricRatioTrendReq struct {
	g.Meta    `path:"/getResourceMetricRatioTrend" method:"post" summary:"获取堆比例资源趋势" tags:"资源监控"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Limit     int    `json:"limit"`
}

type CMSResourceMetricCpuTrendReq struct {
	g.Meta    `path:"/getResourceMetricCpuTrend" method:"post" summary:"获取CPU资源趋势" tags:"资源监控"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Limit     int    `json:"limit"`
}

type CMSResourceMetricOnlineTrendReq struct {
	g.Meta    `path:"/getResourceMetricOnlineTrend" method:"post" summary:"获取在线人数资源趋势" tags:"资源监控"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Limit     int    `json:"limit"`
}

type CMSResourceMetricMemoryPoint struct {
	Time              string  `json:"time"`
	ProcMemMb         float64 `json:"procMemMb"`
	SysMemUsedMb      float64 `json:"sysMemUsedMb"`
	SysMemTotalMb     float64 `json:"sysMemTotalMb"`
	SysMemUsedPercent float64 `json:"sysMemUsedPercent"`
}

type CMSResourceMetricHeapPoint struct {
	Time            string  `json:"time"`
	ProcHeapAllocMb float64 `json:"procHeapAllocMb"`
	ProcHeapInuseMb float64 `json:"procHeapInuseMb"`
	ProcHeapSysMb   float64 `json:"procHeapSysMb"`
}

type CMSResourceMetricRatioPoint struct {
	Time                string  `json:"time"`
	ProcHeapUsedPercent float64 `json:"procHeapUsedPercent"`
	ProcHeapIdlePercent float64 `json:"procHeapIdlePercent"`
}

type CMSResourceMetricCpuPoint struct {
	Time           string  `json:"time"`
	ProcCpuPercent float64 `json:"procCpuPercent"`
	SysCpuPercent  float64 `json:"sysCpuPercent"`
}

type CMSResourceMetricOnlinePoint struct {
	Time        string `json:"time"`
	OnlineCount uint64 `json:"onlineCount,string"`
}

type CMSResourceMetricMemoryTrendRes struct {
	Points []*CMSResourceMetricMemoryPoint `json:"points"`
}

type CMSResourceMetricHeapTrendRes struct {
	Points []*CMSResourceMetricHeapPoint `json:"points"`
}

type CMSResourceMetricRatioTrendRes struct {
	Points []*CMSResourceMetricRatioPoint `json:"points"`
}

type CMSResourceMetricCpuTrendRes struct {
	Points []*CMSResourceMetricCpuPoint `json:"points"`
}

type CMSResourceMetricOnlineTrendRes struct {
	Points []*CMSResourceMetricOnlinePoint `json:"points"`
}
