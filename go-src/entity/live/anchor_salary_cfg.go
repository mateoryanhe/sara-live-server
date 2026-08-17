package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbAnchorSalaryCfg db.TbName = "anchor_salary_cfgs"
)

// AnchorSalaryCfg 主播结算薪资分档配置(CMS 多行管理)
type AnchorSalaryCfg struct {
	migrate.OneModel
	WeeklyWorkDays           uint64  `gorm:"default:0;comment:每周工作天数门槛" json:"weeklyWorkDays"`
	DailyLiveDurationMinutes uint64  `gorm:"default:0;comment:每天直播时长门槛(分钟)" json:"dailyLiveDurationMinutes"`
	SalaryAmount             float64 `gorm:"type:decimal(16,4);default:0;comment:薪资金额" json:"salaryAmount"`
	Sort                     int     `gorm:"default:0;comment:排序值(越大越靠前)" json:"sort"`
}

func initAnchorSalaryCfg() {
	migrate.AutoMigrate(&AnchorSalaryCfg{})
}
