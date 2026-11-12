package anchorvideopurchasedao

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	entity "xr-game-server/entity/live"
)

const purchaseListCacheMaxSize = 2000

var purchaseListCacheMgr *cache.RowCache[[]*entity.AnchorVideoPurchase]

func Init() {
	purchaseListCacheMgr = cache.NewPermanentRowCache[[]*entity.AnchorVideoPurchase]()
}

func purchaseCacheKey(userId, anchorId uint64) string {
	return fmt.Sprintf("%d_%d", userId, anchorId)
}

func getPurchaseList(userId, anchorId uint64) []*entity.AnchorVideoPurchase {
	if purchaseListCacheMgr == nil || userId == 0 || anchorId == 0 {
		return make([]*entity.AnchorVideoPurchase, 0)
	}
	key := purchaseCacheKey(userId, anchorId)
	list := purchaseListCacheMgr.MustGetRow(gctx.New(), key, func(ctx context.Context) ([]*entity.AnchorVideoPurchase, error) {
		return loadPurchasesFromDB(userId, anchorId), nil
	})
	if list == nil {
		return make([]*entity.AnchorVideoPurchase, 0)
	}
	return list
}

func loadPurchasesFromDB(userId, anchorId uint64) []*entity.AnchorVideoPurchase {
	out := make([]*entity.AnchorVideoPurchase, 0)
	if userId == 0 || anchorId == 0 {
		return out
	}
	var rows []*entity.AnchorVideoPurchase
	_ = g.DB().Model(string(entity.TbAnchorVideoPurchase)).
		Where("user_id = ? AND anchor_id = ?", userId, anchorId).
		Order("paid_at desc, id desc").
		Limit(purchaseListCacheMaxSize).
		Scan(&rows)
	for _, row := range rows {
		if row != nil && row.VideoId > 0 {
			out = append(out, row)
		}
	}
	return out
}

func publishPurchaseList(userId, anchorId uint64, list []*entity.AnchorVideoPurchase) {
	if purchaseListCacheMgr == nil || userId == 0 || anchorId == 0 {
		return
	}
	if list == nil {
		list = make([]*entity.AnchorVideoPurchase, 0)
	}
	purchaseListCacheMgr.PublishRow(gctx.New(), purchaseCacheKey(userId, anchorId), list)
}

// GetPurchasedVideoIds 用户对该主播已付费视频ID集合(缓存 key=userId_anchorId).
func GetPurchasedVideoIds(userId, anchorId uint64) map[uint64]struct{} {
	out := make(map[uint64]struct{})
	for _, row := range getPurchaseList(userId, anchorId) {
		if row != nil && row.VideoId > 0 {
			out[row.VideoId] = struct{}{}
		}
	}
	return out
}

// HasPurchased 是否已付费解锁指定视频(读 userId+anchorId 缓存).
func HasPurchased(userId, anchorId, videoId uint64) bool {
	if userId == 0 || anchorId == 0 || videoId == 0 {
		return false
	}
	_, ok := GetPurchasedVideoIds(userId, anchorId)[videoId]
	return ok
}

// SavePurchase 写入付费记录并刷新 userId+anchorId 购买列表缓存.
func SavePurchase(row *entity.AnchorVideoPurchase) {
	if row == nil || row.UserId == 0 || row.VideoId == 0 || row.AnchorId == 0 {
		return
	}
	list := getPurchaseList(row.UserId, row.AnchorId)
	newList := make([]*entity.AnchorVideoPurchase, 0, len(list)+1)
	newList = append(newList, row)
	for _, old := range list {
		if old == nil || old.VideoId == 0 || old.VideoId == row.VideoId {
			continue
		}
		newList = append(newList, old)
	}
	if len(newList) > purchaseListCacheMaxSize {
		newList = newList[:purchaseListCacheMaxSize]
	}
	publishPurchaseList(row.UserId, row.AnchorId, newList)
}
