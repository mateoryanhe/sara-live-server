package resourcemonitor

import (
	"context"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"xr-game-server/core/event"
	"xr-game-server/core/push"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dao/resourcemetricdao"
	"xr-game-server/entity/sys"
	"xr-game-server/gameevent"
)

var currentProcess *process.Process

type cpuBaseline struct {
	at        time.Time
	sysTotal  float64
	sysIdle   float64
	procBusy  float64
	hasSys    bool
	hasProc   bool
}

var (
	cpuMu       sync.Mutex
	cpuPrev     cpuBaseline
	cpuNumLogical = runtime.NumCPU()
)

func initMonitor() {
	var err error
	currentProcess, err = process.NewProcess(int32(os.Getpid()))
	if err != nil {
		currentProcess = nil
	}
	if n, err := cpu.Counts(true); err == nil && n > 0 {
		cpuNumLogical = n
	}
	// 先采一次基线(不阻塞等待),再按细间隔写入
	xrtimer.AddOnce(gctx.New(), 2*time.Second, func(ctx context.Context) {
		recordResourceMetric()
	})
	xrtimer.AddOnce(gctx.New(), 5*time.Second, func(ctx context.Context) {
		backfillMissingAggs()
	})
	xrtimer.AddSingleton(gctx.New(), FineInterval, func(ctx context.Context) {
		recordResourceMetric()
	})
	xrtimer.AddSingleton(gctx.New(), CoarseInterval, func(ctx context.Context) {
		rollupLastCompletedBucket()
		cleanupExpiredResourceMetrics()
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
	now := time.Now()
	resourcemetricdao.DeleteFineBefore(now.Add(-FineRetention))
	resourcemetricdao.DeleteAggBefore(now.Add(-CoarseRetention))
}

func enqueueResourceMetric(now time.Time) {
	procMemMb, procHeapAllocMb, procHeapInuseMb, procHeapSysMb, procHeapUsedPercent, procHeapIdlePercent := collectProcessMemMetrics()
	sysMemUsedMb, sysMemTotalMb, sysMemUsedPercent := collectSystemMemMetrics()
	procCpuPercent, sysCpuPercent := collectCpuPercents(now)
	entity.NewSysResourceMetric(
		now,
		procMemMb, procHeapAllocMb, procHeapInuseMb, procHeapSysMb, procHeapUsedPercent, procHeapIdlePercent, procCpuPercent,
		sysMemUsedMb, sysMemTotalMb, sysMemUsedPercent, sysCpuPercent,
		uint64(push.OnlineCount()),
	)
}

func collectProcessMemMetrics() (procMemMb, procHeapAllocMb, procHeapInuseMb, procHeapSysMb, procHeapUsedPercent, procHeapIdlePercent float64) {
	if currentProcess != nil {
		if memInfo, err := currentProcess.MemoryInfo(); err == nil && memInfo != nil {
			procMemMb = bytesToMb(memInfo.RSS)
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

func collectSystemMemMetrics() (sysMemUsedMb, sysMemTotalMb, sysMemUsedPercent float64) {
	if vm, err := mem.VirtualMemory(); err == nil && vm != nil {
		sysMemUsedMb = bytesToMb(vm.Used)
		sysMemTotalMb = bytesToMb(vm.Total)
		sysMemUsedPercent = vm.UsedPercent
	}
	return
}

// collectCpuPercents 基于两次瞬时计数差分计算占用率,不 sleep
func collectCpuPercents(now time.Time) (procCpuPercent, sysCpuPercent float64) {
	cur := readCpuBaseline(now)

	cpuMu.Lock()
	prev := cpuPrev
	cpuPrev = cur
	cpuMu.Unlock()

	if prev.at.IsZero() {
		return 0, 0
	}
	wallSec := now.Sub(prev.at).Seconds()
	if wallSec <= 0 {
		return 0, 0
	}

	if prev.hasSys && cur.hasSys {
		totalDelta := cur.sysTotal - prev.sysTotal
		idleDelta := cur.sysIdle - prev.sysIdle
		if totalDelta > 0 {
			busy := totalDelta - idleDelta
			if busy < 0 {
				busy = 0
			}
			sysCpuPercent = busy / totalDelta * 100
		}
	}
	if prev.hasProc && cur.hasProc && cpuNumLogical > 0 {
		busyDelta := cur.procBusy - prev.procBusy
		if busyDelta < 0 {
			busyDelta = 0
		}
		procCpuPercent = busyDelta / wallSec / float64(cpuNumLogical) * 100
	}
	return clampPercent(procCpuPercent), clampPercent(sysCpuPercent)
}

func readCpuBaseline(now time.Time) cpuBaseline {
	cur := cpuBaseline{at: now}
	if times, err := cpu.Times(false); err == nil && len(times) > 0 {
		t := times[0]
		cur.sysTotal = t.User + t.System + t.Nice + t.Iowait + t.Irq + t.Softirq + t.Steal + t.Idle + t.Guest + t.GuestNice
		cur.sysIdle = t.Idle
		cur.hasSys = cur.sysTotal > 0
	}
	if currentProcess != nil {
		if t, err := currentProcess.Times(); err == nil && t != nil {
			cur.procBusy = t.User + t.System + t.Iowait + t.Irq + t.Softirq + t.Nice + t.Steal
			cur.hasProc = true
		}
	}
	return cur
}

func clampPercent(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
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
