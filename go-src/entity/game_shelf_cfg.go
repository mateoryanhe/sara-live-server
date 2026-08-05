package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const TbGameShelfCfg db.TbName = "game_shelf_cfgs"

// GameShelfCfg 上架游戏(仅存游戏编码)
type GameShelfCfg struct {
	migrate.OneModel
	GameCode string `gorm:"uniqueIndex;size:64;comment:游戏编码" json:"gameCode"`
}

func initGameShelfCfg() {
	migrate.AutoMigrate(&GameShelfCfg{})
}
