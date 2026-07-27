package game

import (
	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity"
)

// GetAllGameCfgFromCache 获取全部游戏配置(供其它模块使用)
func GetAllGameCfgFromCache() []*entity.GameCfg {
	return cfgdao.GetAllGameCfgCached()
}
