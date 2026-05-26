package game

import "xr-game-server/dao/gamecfgdao"

// Init 服务启动时初始化游戏配置缓存
func Init() {
	gamecfgdao.InitGameCfgDao()
	gamecfgdao.ReloadCache()
}
