package livefollowdao

import (
	"xr-game-server/core/cache"
)

// InitLiveFollowDao 初始化关注主播相关缓存
func InitLiveFollowDao() {
	followCacheMgr = cache.NewCacheMgr()
	followingListCacheMgr = cache.NewCacheMgr()
	followerListCacheMgr = cache.NewCacheMgr()
	blockListCacheMgr = cache.NewCacheMgr()
}
