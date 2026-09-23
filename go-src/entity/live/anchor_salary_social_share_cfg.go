package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbAnchorSalarySocialShareCfg db.TbName = "anchor_salary_social_share_cfgs"
)

// AnchorSalarySocialShareCfg 有底薪主播的社交钻石流水分佣分档配置。
type AnchorSalarySocialShareCfg struct {
	migrate.OneModel
	Level                     uint32  `gorm:"uniqueIndex;default:0;comment:等级" json:"level"`
	SocialTotalDiamondRevenue float64 `gorm:"type:decimal(20,4);default:0;comment:社交总钻石流水" json:"socialTotalDiamondRevenue"`
	AnchorSocialSharePercent  float64 `gorm:"column:social_share_percent;type:decimal(6,2);default:0;comment:主播社交提成比(%)" json:"anchorSocialSharePercent"`
	GuildSocialSharePercent   float64 `gorm:"type:decimal(6,2);default:0;comment:工会社交提成比(%)" json:"guildSocialSharePercent"`
}

func initAnchorSalarySocialShareCfg() {
	migrate.AutoMigrate(&AnchorSalarySocialShareCfg{})
}
