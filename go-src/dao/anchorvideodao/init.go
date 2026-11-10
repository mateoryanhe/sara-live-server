package anchorvideodao

import (
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
)

func Init() {
	rowCacheMgr = cache.NewPermanentRowCache[*entity.AnchorVideo]()
	listCacheMgr = cache.NewPermanentListCache[*entity.AnchorVideo]()
}
