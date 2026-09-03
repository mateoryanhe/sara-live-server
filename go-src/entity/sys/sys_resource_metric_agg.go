package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
)

const (
	TbSysResourceMetricAgg db.TbName = "sys_resource_metric_aggs"
)

// SysResourceMetricAgg 系统资源粗采样(5分钟聚合,保留3天,直写DB)
type SysResourceMetricAgg struct {
	migrate.OneModel
	RecordedAt          time.Time `gorm:"uniqueIndex;comment:聚合桶起始时间" json:"recordedAt"`
	ProcMemMb           float64   `gorm:"type:decimal(12,2);default:0;comment:进程内存MB(RSS)均值" json:"procMemMb"`
	ProcHeapAllocMb     float64   `gorm:"type:decimal(12,2);default:0;comment:进程HeapAlloc MB均值" json:"procHeapAllocMb"`
	ProcHeapInuseMb     float64   `gorm:"type:decimal(12,2);default:0;comment:进程HeapInuse MB均值" json:"procHeapInuseMb"`
	ProcHeapSysMb       float64   `gorm:"type:decimal(12,2);default:0;comment:进程HeapSys MB均值" json:"procHeapSysMb"`
	ProcHeapUsedPercent float64   `gorm:"type:decimal(6,2);default:0;comment:堆使用比例均值" json:"procHeapUsedPercent"`
	ProcHeapIdlePercent float64   `gorm:"type:decimal(6,2);default:0;comment:堆空闲比例均值" json:"procHeapIdlePercent"`
	ProcCpuPercent      float64   `gorm:"type:decimal(6,2);default:0;comment:进程CPU使用率均值" json:"procCpuPercent"`
	SysMemUsedMb        float64   `gorm:"type:decimal(12,2);default:0;comment:系统已用内存MB均值" json:"sysMemUsedMb"`
	SysMemTotalMb       float64   `gorm:"type:decimal(12,2);default:0;comment:系统总内存MB均值" json:"sysMemTotalMb"`
	SysMemUsedPercent   float64   `gorm:"type:decimal(6,2);default:0;comment:系统内存使用率均值" json:"sysMemUsedPercent"`
	SysCpuPercent       float64   `gorm:"type:decimal(6,2);default:0;comment:系统CPU使用率均值" json:"sysCpuPercent"`
	OnlineCount         uint64    `gorm:"default:0;comment:在线人数峰值" json:"onlineCount"`
}

// NewSysResourceMetricAggFromAvg 由细采样聚合结果构造粗采样行
func NewSysResourceMetricAggFromAvg(
	bucketStart time.Time,
	procMemMb, procHeapAllocMb, procHeapInuseMb, procHeapSysMb, procHeapUsedPercent, procHeapIdlePercent, procCpuPercent float64,
	sysMemUsedMb, sysMemTotalMb, sysMemUsedPercent, sysCpuPercent float64,
	onlineCountMax uint64,
) *SysResourceMetricAgg {
	now := time.Now()
	return &SysResourceMetricAgg{
		OneModel: migrate.OneModel{
			ID:        snowflake.GetId(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		RecordedAt:          bucketStart,
		ProcMemMb:           procMemMb,
		ProcHeapAllocMb:     procHeapAllocMb,
		ProcHeapInuseMb:     procHeapInuseMb,
		ProcHeapSysMb:       procHeapSysMb,
		ProcHeapUsedPercent: procHeapUsedPercent,
		ProcHeapIdlePercent: procHeapIdlePercent,
		ProcCpuPercent:      procCpuPercent,
		SysMemUsedMb:        sysMemUsedMb,
		SysMemTotalMb:       sysMemTotalMb,
		SysMemUsedPercent:   sysMemUsedPercent,
		SysCpuPercent:       sysCpuPercent,
		OnlineCount:         onlineCountMax,
	}
}

func initSysResourceMetricAgg() {
	migrate.AutoMigrate(&SysResourceMetricAgg{})
}
