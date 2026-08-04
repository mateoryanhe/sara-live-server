package preload

import "xr-game-server/dao/cfgdao"

func Init() {
	preloadRecentLoginUsers(cfgdao.GetRecentLoginPreloadLimit())
}
