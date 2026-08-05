package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbGameCfg db.TbName = "game_cfgs"
)

// GameCfg 游戏配置(CMS管理)
type GameCfg struct {
	migrate.OneModel
	Name      string `gorm:"size:64;comment:游戏名称" json:"name"`
	Code      string `gorm:"uniqueIndex;size:64;comment:游戏编码" json:"code"`
	LiveCover string `gorm:"size:255;default:'';comment:直播间游戏封面" json:"liveCover"`
	Link      string `gorm:"size:512;default:'';comment:跳转链接" json:"link"`
	Sort      int    `gorm:"default:0;comment:排序值(越大越靠前)" json:"sort"`
}

func initGameCfg() {
	migrate.AutoMigrate(&GameCfg{})
}
