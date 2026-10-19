package resourcemonitor

import "time"

const (
	// FineInterval 细采样间隔(写入 sys_resource_metrics)
	FineInterval = 10 * time.Second
	// FineRetention 细采样保留时长
	FineRetention = 24 * time.Hour
	// CoarseInterval 粗采样聚合桶大小(写入 sys_resource_metric_aggs)
	CoarseInterval = 5 * time.Minute
	// CoarseRetention 粗采样保留时长（与 CMS 仪表盘在线趋势「近 30 天」对齐，略留余量）
	CoarseRetention = 35 * 24 * time.Hour

	maxTrendPoints = 10000
)
