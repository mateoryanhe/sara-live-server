package livefollowdao

import (
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
)

// InitLiveFollowDao 初始化关注主播相关缓存
func InitLiveFollowDao() {
	followCacheMgr = cache.NewRowCache[*entity.LiveFollow]()
	followingListCacheMgr = cache.NewListCache[*entity.LiveFollow]()
	followerListCacheMgr = cache.NewListCache[*RelationUserListRow]()
	blockListCacheMgr = cache.NewListCache[*RelationUserListRow]()
}
