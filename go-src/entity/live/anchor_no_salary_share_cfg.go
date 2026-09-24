package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbAnchorNoSalaryShareCfg db.TbName = "anchor_no_salary_share_cfgs"
)

// AnchorNoSalaryShareCfg 无底薪主播的社交钻石流水分佣分档配置。
type AnchorNoSalaryShareCfg struct {
	migrate.OneModel
	Level                     uint32  `gorm:"uniqueIndex;default:0;comment:等级" json:"level"`
	SocialTotalDiamondRevenue float64 `gorm:"type:decimal(20,4);default:0;comment:社交总钻石流水" json:"socialTotalDiamondRevenue"`
	AnchorSocialSharePercent  float64 `gorm:"type:decimal(6,2);default:0;comment:主播社交提成比(%)" json:"anchorSocialSharePercent"`
	GuildSocialSharePercent   float64 `gorm:"type:decimal(6,2);default:0;comment:工会社交提成比(%)" json:"guildSocialSharePercent"`
}

func initAnchorNoSalaryShareCfg() {
	migrate.AutoMigrate(&AnchorNoSalaryShareCfg{})
}
