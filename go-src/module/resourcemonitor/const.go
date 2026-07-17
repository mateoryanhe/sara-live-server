package resourcemonitor

import "time"

const (
	MetricInterval  = 10 * time.Minute
	MetricRetention = 3 * 24 * time.Hour
	cpuSample       = 500 * time.Millisecond
)
