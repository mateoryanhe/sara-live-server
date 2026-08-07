package syndb

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

const cpuSampleInterval = 200 * time.Millisecond

// sampleSystemCPU 采样系统 CPU,返回使用率与空闲率(0~100).
func sampleSystemCPU() (usedPercent, idlePercent float64) {
	percents, err := cpu.Percent(cpuSampleInterval, false)
	if err != nil || len(percents) == 0 {
		return 0, 100
	}
	usedPercent = percents[0]
	if usedPercent < 0 {
		usedPercent = 0
	}
	if usedPercent > 100 {
		usedPercent = 100
	}
	idlePercent = 100 - usedPercent
	return usedPercent, idlePercent
}
