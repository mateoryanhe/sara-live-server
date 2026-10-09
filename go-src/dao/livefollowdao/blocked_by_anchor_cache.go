package livefollowdao

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
)

// blockedByAnchorListCacheMgr 观众维度:哪些用户(含主播)拉黑了我(anchor_id=观众, user_id=拉黑方)
var blockedByAnchorListCacheMgr *cache.ListCache[uint64]

func blockedByAnchorCacheKey(viewerId uint64) string {
	return fmt.Sprintf("live_follow_blocked_by_anchors:%d", viewerId)
}

func loadBlockedByAnchorIDsFromDB(viewerId uint64) []uint64 {
	ids := make([]uint64, 0)
	if viewerId == 0 {
		return ids
	}
	const relationAlias = "lf"
	_ = g.DB().Model(string(entity.TbLiveFollow)+" "+relationAlias).Ctx(gctx.New()).
		Fields(relationAlias+".user_id").
		Where(relationAlias+".anchor_id = ? AND "+relationAlias+".status = ?", viewerId, entity.LiveFollowStatusBlock).
		Order(relationAlias + ".updated_at desc").
		Scan(&ids)
	return ids
}

func getBlockedByAnchorIDs(viewerId uint64) []uint64 {
	if blockedByAnchorListCacheMgr == nil || viewerId == 0 {
		return make([]uint64, 0)
	}
	list := blockedByAnchorListCacheMgr.MustGetList(gctx.New(), blockedByAnchorCacheKey(viewerId), func(ctx context.Context) ([]uint64, error) {
		return loadBlockedByAnchorIDsFromDB(viewerId), nil
	})
	if list == nil {
		return make([]uint64, 0)
	}
	return list
}

func putBlockedByAnchorCache(viewerId uint64, list []uint64) {
	if blockedByAnchorListCacheMgr == nil || viewerId == 0 {
		return
	}
	if list == nil {
		list = make([]uint64, 0)
	}
	blockedByAnchorListCacheMgr.PublishList(gctx.New(), blockedByAnchorCacheKey(viewerId), list)
}

// GetBlockedByAnchorIDSet 查询「拉黑过该观众的用户 ID」(含主播 roomId); miss 时一次 DB 加载并缓存
func GetBlockedByAnchorIDSet(viewerId uint64) map[uint64]struct{} {
	set := make(map[uint64]struct{})
	if viewerId == 0 {
		return set
	}
	for _, id := range getBlockedByAnchorIDs(viewerId) {
		if id == 0 {
			continue
		}
		set[id] = struct{}{}
	}
	return set
}

// IsViewerBlockedByAnchor 观众是否被该主播(或用户)拉黑
func IsViewerBlockedByAnchor(viewerId, blockerId uint64) bool {
	if viewerId == 0 || blockerId == 0 || viewerId == blockerId {
		return false
	}
	set := GetBlockedByAnchorIDSet(viewerId)
	_, ok := set[blockerId]
	return ok
}

// AddBlockedByAnchorToViewerCache 主播(拉黑方)拉黑观众后,写入观众侧反向列表缓存
func AddBlockedByAnchorToViewerCache(blockerId, viewerId uint64) {
	if blockedByAnchorListCacheMgr == nil || blockerId == 0 || viewerId == 0 || blockerId == viewerId {
		return
	}
	list := getBlockedByAnchorIDs(viewerId)
	for _, id := range list {
		if id == blockerId {
			return
		}
	}
	newList := make([]uint64, 0, len(list)+1)
	newList = append(newList, blockerId)
	newList = append(newList, list...)
	putBlockedByAnchorCache(viewerId, newList)
}

// RemoveBlockedByAnchorFromViewerCache 解除拉黑后,从观众侧反向列表缓存移除拉黑方
func RemoveBlockedByAnchorFromViewerCache(blockerId, viewerId uint64) {
	if blockedByAnchorListCacheMgr == nil || blockerId == 0 || viewerId == 0 {
		return
	}
	if _, ok := blockedByAnchorListCacheMgr.GetListCached(gctx.New(), blockedByAnchorCacheKey(viewerId)); !ok {
		return
	}
	list := getBlockedByAnchorIDs(viewerId)
	newList := make([]uint64, 0, len(list))
	for _, id := range list {
		if id != blockerId {
			newList = append(newList, id)
		}
	}
	putBlockedByAnchorCache(viewerId, newList)
}

func preloadBlockedByAnchorCache(viewerId uint64) {
	putBlockedByAnchorCache(viewerId, loadBlockedByAnchorIDsFromDB(viewerId))
}
