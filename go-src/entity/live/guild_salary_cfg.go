package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbGuildSalaryCfg db.TbName = "guild_salary_cfgs"
)

// GuildSalaryCfg 工会结算薪资分档配置(CMS 多行管理,结算逻辑后续再接)
type GuildSalaryCfg struct {
	migrate.OneModel
	DailyEffectiveLiveCount  uint64  `gorm:"default:0;comment:每天有效直播次数门槛" json:"dailyEffectiveLiveCount"`
	WeeklyEffectiveLiveCount uint64  `gorm:"default:0;comment:每周有效直播次数门槛" json:"weeklyEffectiveLiveCount"`
	SalaryAmount             float64 `gorm:"type:decimal(16,4);default:0;comment:薪资金额" json:"salaryAmount"`
	Sort                     int     `gorm:"default:0;comment:排序值(越大越靠前)" json:"sort"`
}

func initGuildSalaryCfg() {
	migrate.AutoMigrate(&GuildSalaryCfg{})
}
