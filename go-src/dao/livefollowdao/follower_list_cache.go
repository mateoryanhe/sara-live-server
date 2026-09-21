package livefollowdao

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
	userentity "xr-game-server/entity/user"
)

var followerListCacheMgr *cache.ListCache[*RelationUserListRow]

func followerListCacheKey(anchorId uint64) string {
	return fmt.Sprintf("live_follow_follower_list:%d", anchorId)
}

func getFollowerListCache(anchorId uint64) []*RelationUserListRow {
	if followerListCacheMgr == nil || anchorId == 0 {
		return make([]*RelationUserListRow, 0)
	}
	return followerListCacheMgr.MustGetList(gctx.New(), followerListCacheKey(anchorId), func(ctx context.Context) ([]*RelationUserListRow, error) {
		return loadFollowersFromDB(anchorId, 1, followingListCacheMaxSize), nil
	})
}

func putFollowerListCache(anchorId uint64, list []*RelationUserListRow) {
	if followerListCacheMgr == nil || anchorId == 0 {
		return
	}
	if list == nil {
		list = make([]*RelationUserListRow, 0)
	}
	followerListCacheMgr.PublishList(gctx.New(), followerListCacheKey(anchorId), list)
}

// PrependFollowerToListCache 关注成功后写入主播粉丝列表缓存头部
func PrependFollowerToListCache(f *entity.LiveFollow, follower *userentity.UserInfo) {
	if followerListCacheMgr == nil || f == nil || f.UserId == 0 || f.AnchorId == 0 {
		return
	}
	if f.Status != entity.LiveFollowStatusFollow {
		return
	}
	list := getFollowerListCache(f.AnchorId)
	newList := make([]*RelationUserListRow, 0, len(list)+1)
	newList = append(newList, newRelationUserListRow(f, follower))
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
	if _, ok := followerListCacheMgr.GetListCached(gctx.New(), followerListCacheKey(anchorId)); !ok {
		return
	}
	list := getFollowerListCache(anchorId)
	newList := make([]*RelationUserListRow, 0, len(list))
	for _, row := range list {
		if row != nil && row.UserId != userId {
			newList = append(newList, row)
		}
	}
	putFollowerListCache(anchorId, newList)
}

func loadFollowersFromDB(anchorId uint64, page, pageSize int) []*RelationUserListRow {
	list := make([]*RelationUserListRow, 0)
	if anchorId == 0 {
		return list
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > followingListCacheMaxSize {
		pageSize = FollowingListCachePageSize
	}
	const relationAlias = "lf"
	const userAlias = "ui"
	_ = g.DB().Model(string(entity.TbLiveFollow)+" "+relationAlias).Ctx(gctx.New()).
		LeftJoin(string(userentity.TbUserInfo)+" "+userAlias, userAlias+".id = "+relationAlias+".user_id").
		Fields(relationUserListFields(relationAlias, userAlias)...).
		Where(relationAlias+".anchor_id = ? AND "+relationAlias+".status = ?", anchorId, entity.LiveFollowStatusFollow).
		Order(relationAlias + ".updated_at desc").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&list)
	return list
}

// GetFollowersByAnchor 分页获取某主播的粉丝记录(仅 Status == Follow)
// 缓存8页数据,前7页走缓存,第8页起直接查库
func GetFollowersByAnchor(anchorId uint64, page, pageSize int) []*RelationUserListRow {
	list := make([]*RelationUserListRow, 0)
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
