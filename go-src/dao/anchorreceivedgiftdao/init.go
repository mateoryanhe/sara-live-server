package anchorreceivedgiftdao

import (
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
)

func Init() {
	rowCacheMgr = cache.NewPermanentRowCache[*entity.AnchorReceivedGift]()
	listCacheMgr = cache.NewPermanentListCache[*entity.AnchorReceivedGift]()
}
