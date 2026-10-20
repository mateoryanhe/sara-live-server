package resourcemonitor

import "time"

const (
	// FineInterval 细采样间隔(写入 sys_resource_metrics)
	FineInterval = 10 * time.Second
	// FineRetention 细采样保留时长(CMS 在线/资源趋势展示最近 3 天)
	FineRetention = 3 * 24 * time.Hour
	// maxTrendPoints 3 天 × 86400s / 10s ≈ 25920，略留余量
	maxTrendPoints = 26000
)
