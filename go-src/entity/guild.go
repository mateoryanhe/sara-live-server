package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbLiveGuild db.TbName = "live_guilds"
)

// LiveGuild 直播工会
type LiveGuild struct {
	migrate.OneModel
	Name        string `gorm:"size:64;comment:工会名称" json:"name"`
	LeaderId    uint64 `gorm:"default:0;comment:会长/负责人ID" json:"leaderId"`
	Contact     string `gorm:"size:64;comment:联系方式" json:"contact"`
	Description string `gorm:"size:255;comment:工会简介" json:"description"`
	Status      uint8  `gorm:"default:1;comment:状态(0-禁用,1-启用)" json:"status"`
}

func InitLiveGuild() {
	migrate.AutoMigrate(&LiveGuild{})
}
