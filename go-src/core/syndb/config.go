package syndb

import (
	"time"

	"xr-game-server/core/cfg"
)

const (
	defaultTickIntervalMs = 1000
	defaultCpuIdlePercent = 70
	defaultMaxPendingWait = 5000
	defaultBatchSize      = 100
)

type runtimeConfig struct {
	tickInterval   time.Duration
	cpuIdlePercent float64
	maxPendingWait time.Duration
	batchSize      int
}

func loadRuntimeConfig() runtimeConfig {
	c := cfg.SyndbBufferCfgVar
	tickMs := c.TickIntervalMs
	if tickMs <= 0 {
		tickMs = defaultTickIntervalMs
	}
	idlePercent := c.CpuIdlePercent
	if idlePercent <= 0 {
		idlePercent = defaultCpuIdlePercent
	}
	if idlePercent > 100 {
		idlePercent = 100
	}
	waitMs := c.MaxPendingWaitMs
	if waitMs <= 0 {
		waitMs = defaultMaxPendingWait
	}
	batchSize := c.BatchSize
	if batchSize <= 0 {
		batchSize = defaultBatchSize
	}
	return runtimeConfig{
		tickInterval:   time.Duration(tickMs) * time.Millisecond,
		cpuIdlePercent: float64(idlePercent),
		maxPendingWait: time.Duration(waitMs) * time.Millisecond,
		batchSize:      batchSize,
	}
}
