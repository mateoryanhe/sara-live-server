package game

import (
	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/game"
)

// GetAllGameCfgFromCache 获取全部上架游戏配置(供其它模块使用).
func GetAllGameCfgFromCache() []*entity.GameCfg {
	return cfgdao.GetAllGameCfgFromMemory()
}

// GetAllOnShelfGamesFromMemory 获取全部已上架游戏(读 game_cfgs 永久缓存).
func GetAllOnShelfGamesFromMemory() []*entity.GameCfg {
	return cfgdao.GetAllGameCfgFromMemory()
}
