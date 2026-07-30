package livefollowdao

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var followerListCacheMgr *cache.CacheMgr

func followerListCacheKey(anchorId uint64) string {
	return fmt.Sprintf("live_follow_follower_list:%d", anchorId)
}

func getFollowerListCache(anchorId uint64) []*entity.LiveFollow {
	if followerListCacheMgr == nil || anchorId == 0 {
		return make([]*entity.LiveFollow, 0)
	}
	v := followerListCacheMgr.GetData(followerListCacheKey(anchorId), func(ctx context.Context) (interface{}, error) {
		return loadFollowersFromDB(anchorId, 1, followingListCacheMaxSize), nil
	})
	list, _ := v.([]*entity.LiveFollow)
	if list == nil {
		return make([]*entity.LiveFollow, 0)
	}
	return list
}

func putFollowerListCache(anchorId uint64, list []*entity.LiveFollow) {
	if followerListCacheMgr == nil || anchorId == 0 {
		return
	}
	if list == nil {
		list = make([]*entity.LiveFollow, 0)
	}
	followerListCacheMgr.FlushCache(followerListCacheKey(anchorId), list)
}

// PrependFollowerToListCache 关注成功后写入主播粉丝列表缓存头部
func PrependFollowerToListCache(f *entity.LiveFollow) {
	if followerListCacheMgr == nil || f == nil || f.UserId == 0 || f.AnchorId == 0 {
		return
	}
	if f.Status != entity.LiveFollowStatusFollow {
		return
	}
	list := getFollowerListCache(f.AnchorId)
	newList := make([]*entity.LiveFollow, 0, len(list)+1)
	newList = append(newList, f)
	for _, row := range list {
		if row != nil && row.UserId != f.UserId {
			newList = append(newList, row)
		}
	}
	if len(newList) > followingListCacheMaxSize {
		newList = newList[:followingListCacheMaxSize]
	}
	putFollowerListCache(f.AnchorId, newList)
}

// RemoveFollowerFromListCache 取消关注/拉黑后从主播粉丝列表缓存移除
func RemoveFollowerFromListCache(anchorId, userId uint64) {
	if followerListCacheMgr == nil || anchorId == 0 || userId == 0 {
		return
	}
	if followerListCacheMgr.GetFromCache(followerListCacheKey(anchorId)) == nil {
		return
	}
	list := getFollowerListCache(anchorId)
	newList := make([]*entity.LiveFollow, 0, len(list))
	for _, row := range list {
		if row != nil && row.UserId != userId {
			newList = append(newList, row)
		}
	}
	putFollowerListCache(anchorId, newList)
}

func loadFollowersFromDB(anchorId uint64, page, pageSize int) []*entity.LiveFollow {
	list := make([]*entity.LiveFollow, 0)
	if anchorId == 0 {
		return list
	}
	if page <= 0 {
		page = 1
	}
	pageSize = FollowingListCachePageSize
	_ = g.Model(string(entity.TbLiveFollow)).
		Where("anchor_id = ? AND status = ?", anchorId, entity.LiveFollowStatusFollow).
		Order("updated_at desc").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&list)
	return list
}

// GetFollowersByAnchor 分页获取某主播的粉丝记录(仅 Status == Follow)
// 缓存5页数据,前4页且 pageSize=FollowingListCachePageSize 时走缓存,第5页起直接查库
func GetFollowersByAnchor(anchorId uint64, page, pageSize int) []*entity.LiveFollow {
	list := make([]*entity.LiveFollow, 0)
	if anchorId == 0 {
		return list
	}
	if page <= 0 {
		page = 1
	}
	pageSize = FollowingListCachePageSize
	if page <= followingListCachedReadPages {
		cached := getFollowerListCache(anchorId)
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
	return loadFollowersFromDB(anchorId, page, pageSize)
}
