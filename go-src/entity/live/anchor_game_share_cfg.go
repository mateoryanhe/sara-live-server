package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbAnchorGameShareCfg db.TbName = "anchor_game_share_cfgs"
)

const (
	AnchorGameShareSalaryTypeWithSalary uint32 = 1
	AnchorGameShareSalaryTypeNoSalary   uint32 = 2
)

// AnchorGameShareCfg 游戏金币流水档位分佣配置。
// SalaryType 区分有底薪与无底薪主播。
type AnchorGameShareCfg struct {
	migrate.OneModel
	SalaryType             uint32  `gorm:"uniqueIndex:uk_anchor_game_share_salary_level;default:0;comment:薪资类型(1有底薪 2无底薪)" json:"salaryType"`
	Level                  uint32  `gorm:"uniqueIndex:uk_anchor_game_share_salary_level;default:0;comment:等级" json:"level"`
	GameTotalGoldRevenue   float64 `gorm:"type:decimal(20,4);default:0;comment:游戏总金币流水" json:"gameTotalGoldRevenue"`
	AnchorGameSharePercent float64 `gorm:"type:decimal(6,2);default:0;comment:主播游戏提成比(%)" json:"anchorGameSharePercent"`
	GuildGameSharePercent  float64 `gorm:"type:decimal(6,2);default:0;comment:工会游戏提成比(%)" json:"guildGameSharePercent"`
}

func initAnchorGameShareCfg() {
	migrate.AutoMigrate(&AnchorGameShareCfg{})
}
