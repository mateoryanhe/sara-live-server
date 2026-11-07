package anchorphotodao

import (
	"context"
	"fmt"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
)

const DefaultMaxAnchorPhotos = 100

var (
	rowCacheMgr       *cache.RowCache[*entity.AnchorPhoto]
	listCacheMgr      *cache.ListCache[*entity.AnchorPhoto]
	anchorBagMu       sync.Map // anchorId -> *sync.Mutex
	maxPhotosProvider func() int
)

// SetMaxPhotosProvider 注入相册上限读取(如 livecfg);未注入或返回<=0 时用 DefaultMaxAnchorPhotos
func SetMaxPhotosProvider(fn func() int) {
	maxPhotosProvider = fn
}

// MaxPhotos 当前主播相册上限
func MaxPhotos() int {
	if maxPhotosProvider != nil {
		if n := maxPhotosProvider(); n > 0 {
			return n
		}
	}
	return DefaultMaxAnchorPhotos
}

// MaxAnchorPhotos 兼容旧常量名
const MaxAnchorPhotos = DefaultMaxAnchorPhotos

func listCacheKey(anchorId uint64) string {
	return fmt.Sprintf("anchor_photos:%d", anchorId)
}

func lockAnchor(anchorId uint64) *sync.Mutex {
	v, _ := anchorBagMu.LoadOrStore(anchorId, &sync.Mutex{})
	return v.(*sync.Mutex)
}

func loadPhotoListCache(anchorId uint64) []*entity.AnchorPhoto {
	if listCacheMgr == nil || anchorId == 0 {
		return make([]*entity.AnchorPhoto, 0)
	}
	list := listCacheMgr.MustGetList(gctx.New(), listCacheKey(anchorId), func(ctx context.Context) ([]*entity.AnchorPhoto, error) {
		return loadPhotosFromDB(anchorId), nil
	})
	if list == nil {
		return make([]*entity.AnchorPhoto, 0)
	}
	return list
}

// GetPhotoList 读缓存返回主播相册(按创建时间倒序;条数受当前上限截断)
func GetPhotoList(anchorId uint64) []*entity.AnchorPhoto {
	list := loadPhotoListCache(anchorId)
	maxN := MaxPhotos()
	if len(list) > maxN {
		return list[:maxN]
	}
	return list
}

// CountPhotos 当前有效照片数(读缓存,受上限截断)
func CountPhotos(anchorId uint64) int {
	return len(GetPhotoList(anchorId))
}

// GetById 按照片ID读行缓存
func GetById(id uint64) *entity.AnchorPhoto {
	if rowCacheMgr == nil || id == 0 {
		return nil
	}
	return rowCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.AnchorPhoto, error) {
		var row *entity.AnchorPhoto
		_ = g.Model(string(entity.TbAnchorPhoto)).Where(g.Map{
			string(db.IdName): id,
		}).Scan(&row)
		return row, nil
	})
}

// TryAddPhoto 主播新增照片;已满当前上限返回 nil
func TryAddPhoto(anchorId uint64, photo string) *entity.AnchorPhoto {
	if rowCacheMgr == nil || listCacheMgr == nil || anchorId == 0 || photo == "" {
		return nil
	}
	mu := lockAnchor(anchorId)
	mu.Lock()
	defer mu.Unlock()

	maxN := MaxPhotos()
	list := loadPhotoListCache(anchorId)
	if len(list) >= maxN {
		return nil
	}
	row := entity.NewAnchorPhoto(anchorId, photo)
	rowCacheMgr.PublishRow(gctx.New(), row.ID, row)

	newList := make([]*entity.AnchorPhoto, 0, len(list)+1)
	newList = append(newList, row)
	newList = append(newList, list...)
	listCacheMgr.PublishList(gctx.New(), listCacheKey(anchorId), newList)
	return row
}

// SoftDeletePhoto 软删自己的一张照片并刷新列表缓存;成功返回被删行
func SoftDeletePhoto(anchorId, photoId uint64) *entity.AnchorPhoto {
	rows := SoftDeletePhotos(anchorId, []uint64{photoId})
	if len(rows) == 0 {
		return nil
	}
	return rows[0]
}

// SoftDeletePhotos 批量软删自己的照片并刷新列表缓存;返回实际删除成功的行
func SoftDeletePhotos(anchorId uint64, photoIds []uint64) []*entity.AnchorPhoto {
	if rowCacheMgr == nil || listCacheMgr == nil || anchorId == 0 || len(photoIds) == 0 {
		return nil
	}
	want := make(map[uint64]struct{}, len(photoIds))
	for _, id := range photoIds {
		if id > 0 {
			want[id] = struct{}{}
		}
	}
	if len(want) == 0 {
		return nil
	}

	mu := lockAnchor(anchorId)
	mu.Lock()
	defer mu.Unlock()

	deleted := make([]*entity.AnchorPhoto, 0, len(want))
	for id := range want {
		row := GetById(id)
		if row == nil || row.ID == 0 || row.IsDeleted || row.AnchorId != anchorId {
			continue
		}
		row.MarkDeleted()
		rowCacheMgr.PublishRow(gctx.New(), row.ID, row)
		deleted = append(deleted, row)
	}
	if len(deleted) == 0 {
		return nil
	}

	delSet := make(map[uint64]struct{}, len(deleted))
	for _, row := range deleted {
		delSet[row.ID] = struct{}{}
	}
	list := loadPhotoListCache(anchorId)
	newList := make([]*entity.AnchorPhoto, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}
		if _, ok := delSet[item.ID]; ok {
			continue
		}
		newList = append(newList, item)
	}
	listCacheMgr.PublishList(gctx.New(), listCacheKey(anchorId), newList)
	return deleted
}

func loadPhotosFromDB(anchorId uint64) []*entity.AnchorPhoto {
	out := make([]*entity.AnchorPhoto, 0)
	if anchorId == 0 {
		return out
	}
	var rows []*entity.AnchorPhoto
	_ = g.Model(string(entity.TbAnchorPhoto)).
		Where("anchor_id = ? AND is_deleted = 0", anchorId).
		Order("created_at desc, id desc").
		Limit(2000).
		Scan(&rows)
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		out = append(out, row)
	}
	return out
}
