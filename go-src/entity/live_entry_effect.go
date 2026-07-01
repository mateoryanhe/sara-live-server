package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbLiveEntryEffect db.TbName = "live_entry_effects"
)

const (
	LiveEntryEffectStatusOffShelf uint8 = 0
	LiveEntryEffectStatusOnShelf  uint8 = 1
)

// LiveEntryEffect 进场特效配置(CMS 管理)
type LiveEntryEffect struct {
	migrate.OneModel
	Name       string `gorm:"size:64;comment:名称" json:"name"`
	LevelStart int    `gorm:"default:0;comment:等级开始" json:"levelStart"`
	LevelEnd   int    `gorm:"default:0;comment:等级结束" json:"levelEnd"`
	Animation  string `gorm:"size:255;default:'';comment:动画资源" json:"animation"`
	Status     uint8  `gorm:"default:0;comment:状态(0-下架,1-上架)" json:"status"`
}

func initLiveEntryEffect() {
	migrate.AutoMigrate(&LiveEntryEffect{})
}
