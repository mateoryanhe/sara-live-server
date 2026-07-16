package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const (
	TbSysResourceMetric db.TbName = "sys_resource_metrics"
)

const (
	SysResourceMetricRecordedAt          db.TbCol = "recorded_at"
	SysResourceMetricProcMemMb           db.TbCol = "proc_mem_mb"
	SysResourceMetricProcHeapAllocMb     db.TbCol = "proc_heap_alloc_mb"
	SysResourceMetricProcHeapInuseMb     db.TbCol = "proc_heap_inuse_mb"
	SysResourceMetricProcHeapSysMb       db.TbCol = "proc_heap_sys_mb"
	SysResourceMetricProcHeapUsedPercent db.TbCol = "proc_heap_used_percent"
	SysResourceMetricProcHeapIdlePercent db.TbCol = "proc_heap_idle_percent"
	SysResourceMetricProcCpuPercent      db.TbCol = "proc_cpu_percent"
	SysResourceMetricSysMemUsedMb        db.TbCol = "sys_mem_used_mb"
	SysResourceMetricSysMemTotalMb       db.TbCol = "sys_mem_total_mb"
	SysResourceMetricSysMemUsedPercent   db.TbCol = "sys_mem_used_percent"
	SysResourceMetricSysCpuPercent       db.TbCol = "sys_cpu_percent"
)

// SysResourceMetric 系统资源采样记录(每10分钟一条,保留3天,懒缓冲异步入库)
type SysResourceMetric struct {
	migrate.OneModel
	RecordedAt          time.Time `gorm:"index;comment:采样时间" json:"recordedAt"`
	ProcMemMb           float64   `gorm:"type:decimal(12,2);default:0;comment:进程内存MB(RSS)" json:"procMemMb"`
	ProcHeapAllocMb     float64   `gorm:"type:decimal(12,2);default:0;comment:进程HeapAlloc MB" json:"procHeapAllocMb"`
	ProcHeapInuseMb     float64   `gorm:"type:decimal(12,2);default:0;comment:进程HeapInuse MB" json:"procHeapInuseMb"`
	ProcHeapSysMb       float64   `gorm:"type:decimal(12,2);default:0;comment:进程HeapSys MB" json:"procHeapSysMb"`
	ProcHeapUsedPercent float64   `gorm:"type:decimal(6,2);default:0;comment:堆使用比例" json:"procHeapUsedPercent"`
	ProcHeapIdlePercent float64   `gorm:"type:decimal(6,2);default:0;comment:堆空闲比例" json:"procHeapIdlePercent"`
	ProcCpuPercent      float64   `gorm:"type:decimal(6,2);default:0;comment:进程CPU使用率" json:"procCpuPercent"`
	SysMemUsedMb        float64   `gorm:"type:decimal(12,2);default:0;comment:系统已用内存MB" json:"sysMemUsedMb"`
	SysMemTotalMb       float64   `gorm:"type:decimal(12,2);default:0;comment:系统总内存MB" json:"sysMemTotalMb"`
	SysMemUsedPercent   float64   `gorm:"type:decimal(6,2);default:0;comment:系统内存使用率" json:"sysMemUsedPercent"`
	SysCpuPercent       float64   `gorm:"type:decimal(6,2);default:0;comment:系统CPU使用率" json:"sysCpuPercent"`
}

// NewSysResourceMetric 创建资源采样记录并推送到懒缓冲队列
func NewSysResourceMetric(
	recordedAt time.Time,
	procMemMb, procHeapAllocMb, procHeapInuseMb, procHeapSysMb, procHeapUsedPercent, procHeapIdlePercent, procCpuPercent float64,
	sysMemUsedMb, sysMemTotalMb, sysMemUsedPercent, sysCpuPercent float64,
) *SysResourceMetric {
	ret := &SysResourceMetric{}
	ret.ID = snowflake.GetId()
	ret.SetCreatedAt(recordedAt)
	ret.SetUpdatedAt(recordedAt)
	ret.SetRecordedAt(recordedAt)
	ret.SetProcMemMb(procMemMb)
	ret.SetProcHeapAllocMb(procHeapAllocMb)
	ret.SetProcHeapInuseMb(procHeapInuseMb)
	ret.SetProcHeapSysMb(procHeapSysMb)
	ret.SetProcHeapUsedPercent(procHeapUsedPercent)
	ret.SetProcHeapIdlePercent(procHeapIdlePercent)
	ret.SetProcCpuPercent(procCpuPercent)
	ret.SetSysMemUsedMb(sysMemUsedMb)
	ret.SetSysMemTotalMb(sysMemTotalMb)
	ret.SetSysMemUsedPercent(sysMemUsedPercent)
	ret.SetSysCpuPercent(sysCpuPercent)
	return ret
}

func (receiver *SysResourceMetric) SetRecordedAt(val time.Time) {
	receiver.RecordedAt = val
	syndb.AddDataToLazyChan(TbSysResourceMetric, SysResourceMetricRecordedAt, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *SysResourceMetric) SetProcMemMb(val float64) {
	receiver.ProcMemMb = val
	syndb.AddDataToLazyChan(TbSysResourceMetric, SysResourceMetricProcMemMb, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *SysResourceMetric) SetProcHeapAllocMb(val float64) {
	receiver.ProcHeapAllocMb = val
	syndb.AddDataToLazyChan(TbSysResourceMetric, SysResourceMetricProcHeapAllocMb, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *SysResourceMetric) SetProcHeapInuseMb(val float64) {
	receiver.ProcHeapInuseMb = val
	syndb.AddDataToLazyChan(TbSysResourceMetric, SysResourceMetricProcHeapInuseMb, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *SysResourceMetric) SetProcHeapSysMb(val float64) {
	receiver.ProcHeapSysMb = val
	syndb.AddDataToLazyChan(TbSysResourceMetric, SysResourceMetricProcHeapSysMb, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *SysResourceMetric) SetProcHeapUsedPercent(val float64) {
	receiver.ProcHeapUsedPercent = val
	syndb.AddDataToLazyChan(TbSysResourceMetric, SysResourceMetricProcHeapUsedPercent, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *SysResourceMetric) SetProcHeapIdlePercent(val float64) {
	receiver.ProcHeapIdlePercent = val
	syndb.AddDataToLazyChan(TbSysResourceMetric, SysResourceMetricProcHeapIdlePercent, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *SysResourceMetric) SetProcCpuPercent(val float64) {
	receiver.ProcCpuPercent = val
	syndb.AddDataToLazyChan(TbSysResourceMetric, SysResourceMetricProcCpuPercent, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *SysResourceMetric) SetSysMemUsedMb(val float64) {
	receiver.SysMemUsedMb = val
	syndb.AddDataToLazyChan(TbSysResourceMetric, SysResourceMetricSysMemUsedMb, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *SysResourceMetric) SetSysMemTotalMb(val float64) {
	receiver.SysMemTotalMb = val
	syndb.AddDataToLazyChan(TbSysResourceMetric, SysResourceMetricSysMemTotalMb, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *SysResourceMetric) SetSysMemUsedPercent(val float64) {
	receiver.SysMemUsedPercent = val
	syndb.AddDataToLazyChan(TbSysResourceMetric, SysResourceMetricSysMemUsedPercent, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *SysResourceMetric) SetSysCpuPercent(val float64) {
	receiver.SysCpuPercent = val
	syndb.AddDataToLazyChan(TbSysResourceMetric, SysResourceMetricSysCpuPercent, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *SysResourceMetric) SetCreatedAt(val time.Time) {
	receiver.CreatedAt = val
	syndb.AddDataToLazyChan(TbSysResourceMetric, db.CreatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *SysResourceMetric) SetUpdatedAt(val time.Time) {
	receiver.UpdatedAt = val
	syndb.AddDataToLazyChan(TbSysResourceMetric, db.UpdatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func initSysResourceMetric() {
	syndb.RegLazy(TbSysResourceMetric, db.CreatedAtName)
	syndb.RegLazy(TbSysResourceMetric, db.UpdatedAtName)
	syndb.RegLazy(TbSysResourceMetric, SysResourceMetricRecordedAt)
	syndb.RegLazy(TbSysResourceMetric, SysResourceMetricProcMemMb)
	syndb.RegLazy(TbSysResourceMetric, SysResourceMetricProcHeapAllocMb)
	syndb.RegLazy(TbSysResourceMetric, SysResourceMetricProcHeapInuseMb)
	syndb.RegLazy(TbSysResourceMetric, SysResourceMetricProcHeapSysMb)
	syndb.RegLazy(TbSysResourceMetric, SysResourceMetricProcHeapUsedPercent)
	syndb.RegLazy(TbSysResourceMetric, SysResourceMetricProcHeapIdlePercent)
	syndb.RegLazy(TbSysResourceMetric, SysResourceMetricProcCpuPercent)
	syndb.RegLazy(TbSysResourceMetric, SysResourceMetricSysMemUsedMb)
	syndb.RegLazy(TbSysResourceMetric, SysResourceMetricSysMemTotalMb)
	syndb.RegLazy(TbSysResourceMetric, SysResourceMetricSysMemUsedPercent)
	syndb.RegLazy(TbSysResourceMetric, SysResourceMetricSysCpuPercent)

	migrate.AutoMigrate(&SysResourceMetric{})
}
