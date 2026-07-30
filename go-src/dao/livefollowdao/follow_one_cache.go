package livefollowdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var followCacheMgr *cache.CacheMgr

// GetById 按复合ID获取(走缓存)
func GetById(id string) *entity.LiveFollow {
	v := followCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var f *entity.LiveFollow
		_ = g.Model(string(entity.TbLiveFollow)).Ctx(ctx).Where("id = ?", id).Scan(&f)
		return f, nil
	})
	if v == nil {
		return nil
	}
	f, _ := v.(*entity.LiveFollow)
	return f
}

// GetByUserAnchor 按 (userId, anchorId) 获取(走缓存)
func GetByUserAnchor(userId, anchorId uint64) *entity.LiveFollow {
	return GetById(entity.BuildLiveFollowId(userId, anchorId))
}

// AddFollowToCache 关注成功后刷新单条缓存
func AddFollowToCache(f *entity.LiveFollow) {
	if f == nil || followCacheMgr == nil {
		return
	}
	followCacheMgr.FlushCache(f.ID, f)
}

// IsBlocked 查询 userId 是否已拉黑 targetId
func IsBlocked(userId, targetId uint64) bool {
	if userId == 0 || targetId == 0 || userId == targetId {
		return false
	}
	if blocked, ok := isBlockedFromListCache(userId, targetId); ok {
		return blocked
	}
	existing := GetByUserAnchor(userId, targetId)
	return existing != nil && existing.Status == entity.LiveFollowStatusBlock
}

// IsFollowing 查询 userId 是否已关注 anchorId
func IsFollowing(userId, anchorId uint64) bool {
	if userId == 0 || anchorId == 0 || userId == anchorId {
		return false
	}
	if following, ok := isFollowingFromListCache(userId, anchorId); ok {
		return following
	}
	existing := GetByUserAnchor(userId, anchorId)
	return existing != nil && existing.Status == entity.LiveFollowStatusFollow
}

// isFollowingFromListCache 从关注列表缓存判断是否已关注;(结果, 是否已由列表缓存得出结论)
func isFollowingFromListCache(userId, anchorId uint64) (bool, bool) {
	if followingListCacheMgr == nil || userId == 0 {
		return false, false
	}
	list := getFollowingListCache(userId)
	for _, f := range list {
		if f != nil && f.AnchorId == anchorId {
			return true, true
		}
	}
	if len(list) < followingListCacheMaxSize {
		return false, true
	}
	return false, false
}

// isBlockedFromListCache 从拉黑列表缓存判断是否已拉黑;(结果, 是否已由列表缓存得出结论)
func isBlockedFromListCache(userId, targetId uint64) (bool, bool) {
	if blockListCacheMgr == nil || userId == 0 {
		return false, false
	}
	list := getBlockListCache(userId)
	for _, f := range list {
		if f != nil && f.AnchorId == targetId {
			return true, true
		}
	}
	if len(list) < followingListCacheMaxSize {
		return false, true
	}
	return false, false
}
