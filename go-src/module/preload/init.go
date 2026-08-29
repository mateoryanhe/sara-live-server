package preload

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cfg"
	"xr-game-server/core/hotrestart"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dao/liveroomdao"
)

func Init() {
	hotrestart.RegisterHotRestartAuthProvider(cfgdao.GetHotRestartAuth)
	cfg.ApplyMemoryLimit(cfgdao.GetMemoryLimitM())
	initRegisterInitCurrency()
	preloadRecentLoginUsers(cfgdao.GetRecentLoginPreloadLimit())
	preloadLiveRoomGameRecommends()
}

func preloadLiveRoomGameRecommends() {
	ctx := gctx.New()
	roomCount := liveroomdao.PreloadLiveRoomGameRecommendToCache()
	if roomCount == 0 {
		g.Log().Info(ctx, "preload live room game recommends skipped: no rooms")
		return
	}
	g.Log().Infof(ctx, "preload live room game recommends done, roomCount=%d", roomCount)
}
