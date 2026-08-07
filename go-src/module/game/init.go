package game

import (
	"xr-game-server/dao/cfgdao"
)

// Init 服务启动时初始化游戏配置缓存(不拉取第三方).
func Init() {
	cfgdao.InitGameCfgDao()
	cfgdao.ReloadGameCfgCache()
	initVendorGameCache()
	initGamePlatformCfg()
	initGameConsumeRank()
}
