package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbSimulatorCpuKeyword db.TbName = "simulator_cpu_keywords"
)

const (
	SimulatorCpuKeywordKeyword db.TbCol = "keyword"
	SimulatorCpuKeywordRemark  db.TbCol = "remark"
)

// SimulatorCpuKeyword 模拟器拦截 CPU 型号关键词(模糊匹配,CMS 管理)
type SimulatorCpuKeyword struct {
	migrate.OneModel
	Keyword string `gorm:"size:128;uniqueIndex;default:'';comment:CPU关键词(模糊匹配,存小写)" json:"keyword"`
	Remark  string `gorm:"size:255;default:'';comment:备注" json:"remark"`
}

func initSimulatorCpuKeyword() {
	migrate.AutoMigrate(&SimulatorCpuKeyword{})
}
