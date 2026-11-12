package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbGameActivityCfg db.TbName = "game_activity_cfgs"
)

// GameActivityCfg 游戏活动配置(CMS 前端配置; 通常仅一条).
// 先存配置,业务下发/跳转等逻辑后续再接.
type GameActivityCfg struct {
	migrate.OneModel
	ButtonIcon  string `gorm:"size:512;default:'';comment:游戏按钮icon URL" json:"buttonIcon"`
	GameListUrl string `gorm:"size:1024;default:'';comment:游戏列表URL" json:"gameListUrl"`
}

func (GameActivityCfg) TableName() string {
	return string(TbGameActivityCfg)
}

func initGameActivityCfg() {
	migrate.AutoMigrate(&GameActivityCfg{})
}
