package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbLiveRevenueShareCfg db.TbName = "live_revenue_share_cfgs"
)

const (
	LiveRevenueShareCfgAnchorSharePercent db.TbCol = "anchor_share_percent"
	LiveRevenueShareCfgGuildSharePercent  db.TbCol = "guild_share_percent"
)

// LiveRevenueShareCfg 流水分佣配置(CMS 管理,通常仅一条)
type LiveRevenueShareCfg struct {
	migrate.OneModel
	AnchorSharePercent float64 `gorm:"type:decimal(6,2);default:30;comment:主播流水分佣比例(%)" json:"anchorSharePercent"`
	GuildSharePercent  float64 `gorm:"type:decimal(6,2);default:10;comment:工会流水分佣比例(%)" json:"guildSharePercent"`
}

func initLiveRevenueShareCfg() {
	migrate.AutoMigrate(&LiveRevenueShareCfg{})
}
