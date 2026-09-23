package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbAnchorNoSalaryShareCfg db.TbName = "anchor_no_salary_share_cfgs"
)

// AnchorNoSalaryShareCfg 无底薪主播社交钻石流水提成配置。
type AnchorNoSalaryShareCfg struct {
	migrate.OneModel
	AnchorSocialSharePercent float64 `gorm:"type:decimal(6,2);default:10;comment:主播社交提成比(%)" json:"anchorSocialSharePercent"`
	GuildSocialSharePercent  float64 `gorm:"type:decimal(6,2);default:10;comment:工会社交提成比(%)" json:"guildSocialSharePercent"`
}

func initAnchorNoSalaryShareCfg() {
	migrate.AutoMigrate(&AnchorNoSalaryShareCfg{})
}
