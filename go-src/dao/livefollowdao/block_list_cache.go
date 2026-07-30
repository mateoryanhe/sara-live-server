package livefollowdao

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var blockListCacheMgr *cache.CacheMgr

func blockListCacheKey(userId uint64) string {
	return fmt.Sprintf("live_follow_block_list:%d", userId)
}

func getBlockListCache(userId uint64) []*entity.LiveFollow {
	if blockListCacheMgr == nil || userId == 0 {
		return make([]*entity.LiveFollow, 0)
	}
	v := blockListCacheMgr.GetData(blockListCacheKey(userId), func(ctx context.Context) (interface{}, error) {
		return loadBlockedListFromDB(userId, 1, followingListCacheMaxSize), nil
	})
	list, _ := v.([]*entity.LiveFollow)
	if list == nil {
		return make([]*entity.LiveFollow, 0)
	}
	return list
}

func putBlockListCache(userId uint64, list []*entity.LiveFollow) {
	if blockListCacheMgr == nil || userId == 0 {
		return
	}
	if list == nil {
		list = make([]*entity.LiveFollow, 0)
	}
	blockListCacheMgr.FlushCache(blockListCacheKey(userId), list)
}

// PrependBlockedToListCache 拉黑成功后写入列表缓存头部
func PrependBlockedToListCache(f *entity.LiveFollow) {
	if blockListCacheMgr == nil || f == nil || f.UserId == 0 || f.AnchorId == 0 {
		return
	}
	if f.Status != entity.LiveFollowStatusBlock {
		return
	}
	list := getBlockListCache(f.UserId)
	newList := make([]*entity.LiveFollow, 0, len(list)+1)
	newList = append(newList, f)
	for _, row := range list {
		if row != nil && row.AnchorId != f.AnchorId {
			newList = append(newList, row)
		}
	}
	if len(newList) > followingListCacheMaxSize {
		newList = newList[:followingListCacheMaxSize]
	}
	putBlockListCache(f.UserId, newList)
}

// RemoveBlockedFromListCache 解除拉黑后从列表缓存移除
func RemoveBlockedFromListCache(userId, targetId uint64) {
	if blockListCacheMgr == nil || userId == 0 || targetId == 0 {
		return
	}
	if blockListCacheMgr.GetFromCache(blockListCacheKey(userId)) == nil {
		return
	}
	list := getBlockListCache(userId)
	newList := make([]*entity.LiveFollow, 0, len(list))
	for _, row := range list {
		if row != nil && row.AnchorId != targetId {
			newList = append(newList, row)
		}
	}
	putBlockListCache(userId, newList)
}

func loadBlockedListFromDB(userId uint64, page, pageSize int) []*entity.LiveFollow {
	list := make([]*entity.LiveFollow, 0)
	if userId == 0 {
		return list
	}
	if page <= 0 {
		page = 1
	}
	pageSize = FollowingListCachePageSize
	_ = g.Model(string(entity.TbLiveFollow)).
		Where("user_id = ? AND status = ?", userId, entity.LiveFollowStatusBlock).
		Order("updated_at desc").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&list)
	return list
}

// GetBlockedListByUser 分页获取用户当前拉黑列表(仅 Status == Block)
// 缓存5页数据,前4页且 pageSize=FollowingListCachePageSize 时走缓存,第5页起直接查库
func GetBlockedListByUser(userId uint64, page, pageSize int) []*entity.LiveFollow {
	list := make([]*entity.LiveFollow, 0)
	if userId == 0 {
		return list
	}
	if page <= 0 {
		page = 1
	}
	pageSize = FollowingListCachePageSize
	if page <= followingListCachedReadPages {
		cached := getBlockListCache(userId)
		start := (page - 1) * pageSize
		if start >= len(cached) {
			return list
		}
		end := start + pageSize
		if end > len(cached) {
			end = len(cached)
		}
		return cached[start:end]
	}
	return loadBlockedListFromDB(userId, page, pageSize)
}
