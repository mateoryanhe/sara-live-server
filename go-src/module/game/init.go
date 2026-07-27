package game

import "xr-game-server/dao/cfgdao"

// Init 服务启动时初始化游戏配置缓存
func Init() {
	cfgdao.InitGameCfgDao()
	cfgdao.ReloadGameCfgCache()
}
