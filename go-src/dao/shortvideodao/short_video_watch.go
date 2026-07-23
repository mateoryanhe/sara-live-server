package shortvideodao

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

const (
	WatchListCachePageSize = 30
	watchListCachedPages   = 2
	watchListCacheMaxSize  = WatchListCachePageSize * watchListCachedPages
)

var (
	watchOneCacheMgr  *cache.CacheMgr
	watchListCacheMgr *cache.CacheMgr
)

func initShortVideoWatchDao() {
	watchOneCacheMgr = cache.NewCacheMgr()
	watchListCacheMgr = cache.NewCacheMgr()
}

func watchCacheKey(userId, videoId uint64) string {
	return fmt.Sprintf("%v_%v", userId, videoId)
}

func watchListCacheKey(userId uint64) string {
	return fmt.Sprintf("watch_list_%v", userId)
}

func GetOneShortVideoWatch(userId, videoId uint64) *entity.ShortVideoWatch {
	id := watchCacheKey(userId, videoId)
	v := watchOneCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var watch entity.ShortVideoWatch
		if scanErr := g.DB().Model(string(entity.TbShortVideoWatch)).
			Where("user_id = ? AND video_id = ?", userId, videoId).
			Scan(&watch); scanErr == nil && watch.ID != "" {
			return &watch, nil
		}
		return entity.NewShortVideoWatch(userId, videoId), nil
	})
	return v.(*entity.ShortVideoWatch)
}

// GetShortVideoWatch 分页查询用户可继续观看的记录(已上架且免费/已付费/作者本人/仍有试看)
// 前2页且 pageSize=30 时走缓存(首次加载2页共60条), 第3页起直接查库
func GetShortVideoWatch(userId uint64, pageIndex, pageSize int) []*entity.ShortVideoWatch {
	if userId == 0 {
		return make([]*entity.ShortVideoWatch, 0)
	}
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = WatchListCachePageSize
	}

	if pageIndex <= watchListCachedPages && pageSize == WatchListCachePageSize {
		cached := getWatchListCache(userId)
		start := (pageIndex - 1) * pageSize
		if start >= len(cached) {
			return make([]*entity.ShortVideoWatch, 0)
		}
		end := start + pageSize
		if end > len(cached) {
			end = len(cached)
		}
		return cached[start:end]
	}
	return loadShortVideoWatchFromDB(userId, pageIndex, pageSize)
}

func getWatchListCache(userId uint64) []*entity.ShortVideoWatch {
	key := watchListCacheKey(userId)
	v := watchListCacheMgr.GetData(key, func(ctx context.Context) (value interface{}, err error) {
		return loadShortVideoWatchFromDB(userId, 1, WatchListCachePageSize*watchListCachedPages), nil
	})
	list, _ := v.([]*entity.ShortVideoWatch)
	if list == nil {
		return make([]*entity.ShortVideoWatch, 0)
	}
	return list
}

// PrependWatchToWatchListCache 将观看记录插入列表缓存头部; 无缓存时先从库加载再写入,最多保留60条
func PrependWatchToWatchListCache(userId uint64, watch *entity.ShortVideoWatch) {
	if watchListCacheMgr == nil || userId == 0 || watch == nil || watch.VideoId == 0 {
		return
	}
	key := watchListCacheKey(userId)
	list := getWatchListCacheForUpdate(userId)
	newList := prependWatchToWatchListCache(list, watch)
	watchListCacheMgr.FlushCache(key, newList)
}

func getWatchListCacheForUpdate(userId uint64) []*entity.ShortVideoWatch {
	key := watchListCacheKey(userId)
	if v := watchListCacheMgr.GetFromCache(key); v != nil {
		if list, ok := v.([]*entity.ShortVideoWatch); ok && list != nil {
			return list
		}
	}
	return loadShortVideoWatchFromDB(userId, 1, watchListCacheMaxSize)
}

func prependWatchToWatchListCache(list []*entity.ShortVideoWatch, watch *entity.ShortVideoWatch) []*entity.ShortVideoWatch {
	newList := make([]*entity.ShortVideoWatch, 0, len(list)+1)
	newList = append(newList, watch)
	for _, row := range list {
		if row == nil || row.VideoId == watch.VideoId {
			continue
		}
		newList = append(newList, row)
	}
	if len(newList) > watchListCacheMaxSize {
		newList = newList[:watchListCacheMaxSize]
	}
	return newList
}

func loadShortVideoWatchFromDB(userId uint64, pageIndex, pageSize int) []*entity.ShortVideoWatch {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = WatchListCachePageSize
	}

	ctx := gctx.New()
	whereSQL := `FROM short_video_watches w
INNER JOIN short_videos v ON v.id = w.video_id
WHERE w.user_id = ?
AND v.status = ?
AND (
	v.is_paid != ?
	OR v.author_id = ?
	OR w.paid_time IS NOT NULL
	OR w.free_time > 0
)`
	args := []any{
		userId,
		entity.ShortVideoStatusOnShelf,
		entity.ShortVideoPaidYes,
		userId,
	}

	list := make([]*entity.ShortVideoWatch, 0)
	querySQL := `SELECT w.* ` + whereSQL + `
ORDER BY w.updated_at DESC, w.id DESC
LIMIT ? OFFSET ?`
	queryArgs := append(append([]any{}, args...), pageSize, (pageIndex-1)*pageSize)
	_ = g.DB().GetScan(ctx, &list, querySQL, queryArgs...)
	return list
}

// DeleteWatchByVideoId 删除指定视频的观看/点赞记录,并清理相关缓存
func DeleteWatchByVideoId(videoId uint64) error {
	if videoId == 0 {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbShortVideoWatch)).Where("video_id = ?", videoId).Delete()
	if err != nil {
		return err
	}
	removeWatchCacheByVideoId(videoId)
	removeWatchListCacheByVideoId(videoId)
	return nil
}

func removeWatchCacheByVideoId(videoId uint64) {
	if watchOneCacheMgr == nil || videoId == 0 {
		return
	}
	ctx := gctx.New()
	keys, err := watchOneCacheMgr.Cache.Keys(ctx)
	if err != nil || len(keys) == 0 {
		return
	}
	videoIdStr := strconv.FormatUint(videoId, 10)
	for _, key := range keys {
		keyStr, ok := key.(string)
		if !ok {
			continue
		}
		idx := strings.LastIndex(keyStr, "_")
		if idx < 0 || keyStr[idx+1:] != videoIdStr {
			continue
		}
		_, _ = watchOneCacheMgr.Cache.Remove(ctx, key)
	}
}

func removeWatchListCacheByVideoId(videoId uint64) {
	if watchListCacheMgr == nil || videoId == 0 {
		return
	}
	ctx := gctx.New()
	keys, err := watchListCacheMgr.Cache.Keys(ctx)
	if err != nil || len(keys) == 0 {
		return
	}
	prefix := "watch_list_"
	for _, key := range keys {
		keyStr, ok := key.(string)
		if !ok || !strings.HasPrefix(keyStr, prefix) {
			continue
		}
		value, _ := watchListCacheMgr.Cache.Get(ctx, key)
		if value == nil || value.IsNil() {
			continue
		}
		list, ok := value.Val().([]*entity.ShortVideoWatch)
		if !ok {
			continue
		}
		for _, row := range list {
			if row != nil && row.VideoId == videoId {
				_, _ = watchListCacheMgr.Cache.Remove(ctx, key)
				break
			}
		}
	}
}
