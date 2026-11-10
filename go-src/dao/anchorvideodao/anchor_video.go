package anchorvideodao

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

const DefaultMaxAnchorVideos = 100

var (
	rowCacheMgr       *cache.RowCache[*entity.AnchorVideo]
	listCacheMgr      *cache.ListCache[*entity.AnchorVideo]
	anchorBagMu       sync.Map // anchorId -> *sync.Mutex
	maxVideosProvider func() int
)

// SetMaxVideosProvider 注入视频上限读取;未注入或返回<=0 时用 DefaultMaxAnchorVideos
func SetMaxVideosProvider(fn func() int) {
	maxVideosProvider = fn
}

// MaxVideos 当前主播视频上限
func MaxVideos() int {
	if maxVideosProvider != nil {
		if n := maxVideosProvider(); n > 0 {
			return n
		}
	}
	return DefaultMaxAnchorVideos
}

func listCacheKey(anchorId uint64) string {
	return fmt.Sprintf("anchor_videos:%d", anchorId)
}

func lockAnchor(anchorId uint64) *sync.Mutex {
	v, _ := anchorBagMu.LoadOrStore(anchorId, &sync.Mutex{})
	return v.(*sync.Mutex)
}

func loadVideoListCache(anchorId uint64) []*entity.AnchorVideo {
	if listCacheMgr == nil || anchorId == 0 {
		return make([]*entity.AnchorVideo, 0)
	}
	list := listCacheMgr.MustGetList(gctx.New(), listCacheKey(anchorId), func(ctx context.Context) ([]*entity.AnchorVideo, error) {
		return loadVideosFromDB(anchorId), nil
	})
	if list == nil {
		return make([]*entity.AnchorVideo, 0)
	}
	return list
}

// GetVideoList 读缓存返回主播视频(按创建时间倒序;条数受当前上限截断)
func GetVideoList(anchorId uint64) []*entity.AnchorVideo {
	list := loadVideoListCache(anchorId)
	maxN := MaxVideos()
	if len(list) > maxN {
		return list[:maxN]
	}
	return list
}

// CountVideos 当前有效视频数(读缓存,受上限截断)
func CountVideos(anchorId uint64) int {
	return len(GetVideoList(anchorId))
}

// GetById 按视频ID读行缓存
func GetById(id uint64) *entity.AnchorVideo {
	if rowCacheMgr == nil || id == 0 {
		return nil
	}
	return rowCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.AnchorVideo, error) {
		var row *entity.AnchorVideo
		_ = g.Model(string(entity.TbAnchorVideo)).Where(g.Map{
			string(db.IdName): id,
		}).Scan(&row)
		return row, nil
	})
}

// TryAddVideo 主播新增视频;已满当前上限返回 nil
func TryAddVideo(anchorId uint64, video string) *entity.AnchorVideo {
	if rowCacheMgr == nil || listCacheMgr == nil || anchorId == 0 || video == "" {
		return nil
	}
	mu := lockAnchor(anchorId)
	mu.Lock()
	defer mu.Unlock()

	maxN := MaxVideos()
	list := loadVideoListCache(anchorId)
	if len(list) >= maxN {
		return nil
	}
	row := entity.NewAnchorVideo(anchorId, video)
	rowCacheMgr.PublishRow(gctx.New(), row.ID, row)

	newList := make([]*entity.AnchorVideo, 0, len(list)+1)
	newList = append(newList, row)
	newList = append(newList, list...)
	listCacheMgr.PublishList(gctx.New(), listCacheKey(anchorId), newList)
	return row
}

// SoftDeleteVideos 批量软删自己的视频并刷新列表缓存;返回实际删除成功的行
func SoftDeleteVideos(anchorId uint64, videoIds []uint64) []*entity.AnchorVideo {
	if rowCacheMgr == nil || listCacheMgr == nil || anchorId == 0 || len(videoIds) == 0 {
		return nil
	}
	want := make(map[uint64]struct{}, len(videoIds))
	for _, id := range videoIds {
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

	deleted := make([]*entity.AnchorVideo, 0, len(want))
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
	list := loadVideoListCache(anchorId)
	newList := make([]*entity.AnchorVideo, 0, len(list))
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

func loadVideosFromDB(anchorId uint64) []*entity.AnchorVideo {
	out := make([]*entity.AnchorVideo, 0)
	if anchorId == 0 {
		return out
	}
	var rows []*entity.AnchorVideo
	_ = g.Model(string(entity.TbAnchorVideo)).
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
