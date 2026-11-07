package anchorreceivedgiftdao

import (
	"context"
	"fmt"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
)

const MaxReceivedGiftIds = 2000

var (
	rowCacheMgr  *cache.RowCache[*entity.AnchorReceivedGift]
	listCacheMgr *cache.ListCache[*entity.AnchorReceivedGift]
	anchorBagMu  sync.Map // anchorId -> *sync.Mutex
)

func listCacheKey(anchorId uint64) string {
	return fmt.Sprintf("anchor_received_gifts:%d", anchorId)
}

func lockAnchor(anchorId uint64) *sync.Mutex {
	v, _ := anchorBagMu.LoadOrStore(anchorId, &sync.Mutex{})
	return v.(*sync.Mutex)
}

func loadListCache(anchorId uint64) []*entity.AnchorReceivedGift {
	if listCacheMgr == nil || anchorId == 0 {
		return make([]*entity.AnchorReceivedGift, 0)
	}
	list := listCacheMgr.MustGetList(gctx.New(), listCacheKey(anchorId), func(ctx context.Context) ([]*entity.AnchorReceivedGift, error) {
		return loadFromDB(anchorId), nil
	})
	if list == nil {
		return make([]*entity.AnchorReceivedGift, 0)
	}
	return list
}

// GetReceivedGifts 读缓存返回主播收礼列表(含数量,最多 MaxReceivedGiftIds)
func GetReceivedGifts(anchorId uint64) []*entity.AnchorReceivedGift {
	return loadListCache(anchorId)
}

// GetReceivedGiftIds 读缓存返回主播收到过的礼物ID列表
func GetReceivedGiftIds(anchorId uint64) []uint64 {
	list := loadListCache(anchorId)
	out := make([]uint64, 0, len(list))
	for _, row := range list {
		if row == nil || row.GiftId == 0 {
			continue
		}
		out = append(out, row.GiftId)
	}
	return out
}

// RecordReceivedGift 累计主播收到的礼物数量;新礼物满 MaxReceivedGiftIds 后忽略
func RecordReceivedGift(anchorId, giftId, addCount uint64) bool {
	if rowCacheMgr == nil || listCacheMgr == nil || anchorId == 0 || giftId == 0 {
		return false
	}
	if addCount == 0 {
		addCount = 1
	}
	mu := lockAnchor(anchorId)
	mu.Lock()
	defer mu.Unlock()

	list := loadListCache(anchorId)
	for _, item := range list {
		if item == nil || item.GiftId != giftId {
			continue
		}
		item.AddCount(addCount)
		rowCacheMgr.PublishRow(gctx.New(), item.ID, item)
		listCacheMgr.PublishList(gctx.New(), listCacheKey(anchorId), list)
		return true
	}

	if len(list) >= MaxReceivedGiftIds {
		return false
	}

	rowId := entity.BuildAnchorReceivedGiftId(anchorId, giftId)
	if cached, ok := rowCacheMgr.GetRowCached(gctx.New(), rowId); ok && cached != nil && cached.GiftId == giftId {
		cached.AddCount(addCount)
		rowCacheMgr.PublishRow(gctx.New(), cached.ID, cached)
		newList := make([]*entity.AnchorReceivedGift, 0, len(list)+1)
		newList = append(newList, list...)
		newList = append(newList, cached)
		listCacheMgr.PublishList(gctx.New(), listCacheKey(anchorId), newList)
		return true
	}
	var existing *entity.AnchorReceivedGift
	_ = g.Model(string(entity.TbAnchorReceivedGift)).Where("id = ?", rowId).Scan(&existing)
	if existing != nil && existing.GiftId == giftId {
		existing.AddCount(addCount)
		rowCacheMgr.PublishRow(gctx.New(), existing.ID, existing)
		newList := make([]*entity.AnchorReceivedGift, 0, len(list)+1)
		newList = append(newList, list...)
		newList = append(newList, existing)
		listCacheMgr.PublishList(gctx.New(), listCacheKey(anchorId), newList)
		return true
	}

	row := entity.NewAnchorReceivedGift(anchorId, giftId, addCount)
	rowCacheMgr.PublishRow(gctx.New(), row.ID, row)
	newList := make([]*entity.AnchorReceivedGift, 0, len(list)+1)
	newList = append(newList, list...)
	newList = append(newList, row)
	listCacheMgr.PublishList(gctx.New(), listCacheKey(anchorId), newList)
	return true
}

// TryRecordReceivedGift 兼容旧调用:数量按 1 累计
func TryRecordReceivedGift(anchorId, giftId uint64) bool {
	return RecordReceivedGift(anchorId, giftId, 1)
}

func loadFromDB(anchorId uint64) []*entity.AnchorReceivedGift {
	out := make([]*entity.AnchorReceivedGift, 0)
	if anchorId == 0 {
		return out
	}
	var rows []*entity.AnchorReceivedGift
	_ = g.Model(string(entity.TbAnchorReceivedGift)).
		Where("anchor_id = ?", anchorId).
		Order("created_at asc, id asc").
		Limit(MaxReceivedGiftIds).
		Scan(&rows)
	for _, row := range rows {
		if row == nil || row.GiftId == 0 {
			continue
		}
		if row.Count == 0 {
			row.Count = 1
		}
		out = append(out, row)
	}
	return out
}
