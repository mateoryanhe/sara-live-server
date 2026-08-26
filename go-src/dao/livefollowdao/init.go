package livefollowdao

import (
	"xr-game-server/entity/live"
	"xr-game-server/core/cache"
)

// InitLiveFollowDao 初始化关注主播相关缓存
func InitLiveFollowDao() {
	followCacheMgr = cache.NewRowCache[*entity.LiveFollow]()
	followingListCacheMgr = cache.NewListCache[*entity.LiveFollow]()
	followerListCacheMgr = cache.NewListCache[*entity.LiveFollow]()
	blockListCacheMgr = cache.NewListCache[*entity.LiveFollow]()
}
