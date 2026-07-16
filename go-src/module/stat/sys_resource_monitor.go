package stat

import (
	"context"
	"os"
	"runtime"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"xr-game-server/core/event"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dao/statdao"
	"xr-game-server/entity"
	"xr-game-server/gameevent"
)

const (
	resourceMetricInterval  = 10 * time.Minute
	resourceMetricRetention = 3 * 24 * time.Hour
	resourceCpuSample       = 500 * time.Millisecond
)

var currentProcess *process.Process

func initResourceMonitor() {
	var err error
	currentProcess, err = process.NewProcess(int32(os.Getpid()))
	if err != nil {
		currentProcess = nil
	}
	xrtimer.AddOnce(gctx.New(), time.Minute, func(ctx context.Context) {
		recordResourceMetric()
	})
	xrtimer.AddSingleton(gctx.New(), resourceMetricInterval, func(ctx context.Context) {
		recordResourceMetric()
	})
	event.Sub(gameevent.DayEvent, onDayCleanupResourceMetrics)
}

func recordResourceMetric() {
	enqueueResourceMetric(time.Now())
}

func onDayCleanupResourceMetrics(_ any) {
	cleanupExpiredResourceMetrics()
}

func cleanupExpiredResourceMetrics() {
	statdao.DeleteSysResourceMetricsBefore(time.Now().Add(-resourceMetricRetention))
}

func enqueueResourceMetric(now time.Time) {
	procMemMb, procHeapAllocMb, procHeapInuseMb, procHeapSysMb, procHeapUsedPercent, procHeapIdlePercent, procCpuPercent := collectProcessMetrics()
	sysMemUsedMb, sysMemTotalMb, sysMemUsedPercent, sysCpuPercent := collectSystemMetrics()
	entity.NewSysResourceMetric(
		now,
		procMemMb, procHeapAllocMb, procHeapInuseMb, procHeapSysMb, procHeapUsedPercent, procHeapIdlePercent, procCpuPercent,
		sysMemUsedMb, sysMemTotalMb, sysMemUsedPercent, sysCpuPercent,
	)
}

func collectProcessMetrics() (procMemMb, procHeapAllocMb, procHeapInuseMb, procHeapSysMb, procHeapUsedPercent, procHeapIdlePercent, procCpuPercent float64) {
	if currentProcess != nil {
		if memInfo, err := currentProcess.MemoryInfo(); err == nil && memInfo != nil {
			procMemMb = bytesToMb(memInfo.RSS)
		}
		if cpuPct, err := currentProcess.Percent(resourceCpuSample); err == nil {
			procCpuPercent = cpuPct
		}
	}

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	procHeapAllocMb = bytesToMb(ms.HeapAlloc)
	procHeapInuseMb = bytesToMb(ms.HeapInuse)
	procHeapSysMb = bytesToMb(ms.HeapSys)
	procHeapUsedPercent, procHeapIdlePercent = calcHeapPercent(ms.HeapInuse, ms.HeapIdle, ms.HeapSys)
	if procMemMb <= 0 {
		procMemMb = bytesToMb(ms.Alloc)
	}
	return
}

func collectSystemMetrics() (sysMemUsedMb, sysMemTotalMb, sysMemUsedPercent, sysCpuPercent float64) {
	if vm, err := mem.VirtualMemory(); err == nil && vm != nil {
		sysMemUsedMb = bytesToMb(vm.Used)
		sysMemTotalMb = bytesToMb(vm.Total)
		sysMemUsedPercent = vm.UsedPercent
	}
	if percents, err := cpu.Percent(resourceCpuSample, false); err == nil && len(percents) > 0 {
		sysCpuPercent = percents[0]
	}
	return
}

func bytesToMb(val uint64) float64 {
	return float64(val) / 1024 / 1024
}

func calcHeapPercent(heapInuse, heapIdle, heapSys uint64) (usedPercent, idlePercent float64) {
	if heapSys == 0 {
		return 0, 0
	}
	usedPercent = float64(heapInuse) / float64(heapSys) * 100
	idlePercent = float64(heapIdle) / float64(heapSys) * 100
	return
}
