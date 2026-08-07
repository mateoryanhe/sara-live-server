package cfg

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

const SyndbBufferStr = "bufferSize.db"

// SyndbBufferCfg syndb 统一缓冲落库配置
type SyndbBufferCfg struct {
	// TickIntervalMs 调度周期(毫秒),默认 1000
	TickIntervalMs int
	// CpuIdlePercent 系统 CPU 空闲比例阈值(0~100),达到则允许批量落库,默认 70
	CpuIdlePercent int
	// MaxPendingWaitMs Pending 最长等待(毫秒),超时强制落库,默认 5000
	MaxPendingWaitMs int
	// BatchSize 每轮全局批量落库上限(所有 queue 共用),默认 100
	BatchSize int
}

var SyndbBufferCfgVar = &SyndbBufferCfg{}

func initDbBufferCfg() {
	initSyndbBuffer()
}

func initSyndbBuffer() {
	data, _ := g.Cfg().GetWithCmd(gctx.New(), SyndbBufferStr)
	if err := data.Scan(SyndbBufferCfgVar); err != nil {
		g.Log().Error(gctx.New(), "无法加载 syndb 缓冲配置")
	}
}
