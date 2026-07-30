package livefollowdao

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

const (
	FollowingListCachePageSize   = 50
	followingListCachedReadPages = 4
	followingListCacheMaxPages   = 5
	followingListCacheMaxSize    = FollowingListCachePageSize * followingListCacheMaxPages
)

var followingListCacheMgr *cache.CacheMgr

func followingListCacheKey(userId uint64) string {
	return fmt.Sprintf("live_follow_following_list:%d", userId)
}

func getFollowingListCache(userId uint64) []*entity.LiveFollow {
	if followingListCacheMgr == nil || userId == 0 {
		return make([]*entity.LiveFollow, 0)
	}
	v := followingListCacheMgr.GetData(followingListCacheKey(userId), func(ctx context.Context) (interface{}, error) {
		return loadFollowingsFromDB(userId, 1, followingListCacheMaxSize), nil
	})
	list, _ := v.([]*entity.LiveFollow)
	if list == nil {
		return make([]*entity.LiveFollow, 0)
	}
	return list
}

func putFollowingListCache(userId uint64, list []*entity.LiveFollow) {
	if followingListCacheMgr == nil || userId == 0 {
		return
	}
	if list == nil {
		list = make([]*entity.LiveFollow, 0)
	}
	followingListCacheMgr.FlushCache(followingListCacheKey(userId), list)
}

// PrependFollowingToListCache 关注成功后写入列表缓存头部
func PrependFollowingToListCache(f *entity.LiveFollow) {
	if followingListCacheMgr == nil || f == nil || f.UserId == 0 || f.AnchorId == 0 {
		return
	}
	if f.Status != entity.LiveFollowStatusFollow {
		return
	}
	list := getFollowingListCache(f.UserId)
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
	putFollowingListCache(f.UserId, newList)
}

// RemoveFollowingFromListCache 取消关注/拉黑后从列表缓存移除
func RemoveFollowingFromListCache(userId, anchorId uint64) {
	if followingListCacheMgr == nil || userId == 0 || anchorId == 0 {
		return
	}
	if followingListCacheMgr.GetFromCache(followingListCacheKey(userId)) == nil {
		return
	}
	list := getFollowingListCache(userId)
	newList := make([]*entity.LiveFollow, 0, len(list))
	for _, row := range list {
		if row != nil && row.AnchorId != anchorId {
			newList = append(newList, row)
		}
	}
	putFollowingListCache(userId, newList)
}

// GetFollowingsByUser 分页获取某用户当前已关注的记录(仅 Status == Follow)
// 缓存5页数据,前4页且 pageSize=FollowingListCachePageSize 时走缓存,第5页起直接查库
func GetFollowingsByUser(userId uint64, page, pageSize int) []*entity.LiveFollow {
	list := make([]*entity.LiveFollow, 0)
	if userId == 0 {
		return list
	}
	if page <= 0 {
		page = 1
	}
	pageSize = FollowingListCachePageSize
	if page <= followingListCachedReadPages {
		cached := getFollowingListCache(userId)
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
	return loadFollowingsFromDB(userId, page, pageSize)
}

func loadFollowingsFromDB(userId uint64, page, pageSize int) []*entity.LiveFollow {
	list := make([]*entity.LiveFollow, 0)
	if userId == 0 {
		return list
	}
	if page <= 0 {
		page = 1
	}
	pageSize = FollowingListCachePageSize
	_ = g.Model(string(entity.TbLiveFollow)).
		Where("user_id = ? AND status = ?", userId, entity.LiveFollowStatusFollow).
		Order("updated_at desc").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&list)
	return list
}
