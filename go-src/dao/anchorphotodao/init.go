package anchorphotodao

import (
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
)

func Init() {
	rowCacheMgr = cache.NewPermanentRowCache[*entity.AnchorPhoto]()
	listCacheMgr = cache.NewPermanentListCache[*entity.AnchorPhoto]()
}
