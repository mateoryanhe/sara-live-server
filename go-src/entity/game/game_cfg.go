package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbGameCfg db.TbName = "game_cfgs"
)

// GameCfg 上架游戏配置(写库字段: gameCode / cover / nameEn / platform)
type GameCfg struct {
	migrate.OneModel
	GameCode string `gorm:"uniqueIndex;size:64;comment:游戏编码" json:"gameCode"`
	Cover    string `gorm:"size:512;default:'';comment:封面" json:"cover"`
	NameEn   string `gorm:"size:128;default:'';comment:英文名称" json:"nameEn"`
	Platform string `gorm:"size:32;default:'';comment:平台编码" json:"platform"`
}

func initGameCfg() {
	migrate.AutoMigrate(&GameCfg{})
}
