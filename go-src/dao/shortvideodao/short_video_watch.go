package shortvideodao

import (
	"context"
	"sort"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

const (
	WatchListCachePageSize = 30
	watchCacheMaxSize      = 1000
)

var watchCacheMgr *cache.CacheMgr

func initShortVideoWatchDao() {
	watchCacheMgr = cache.NewCacheMgr()
}

type userWatchMap = gmap.KVMap[uint64, *entity.ShortVideoWatch]

// getWatchMap 统一入口: key=userId, val=videoId->观看记录
func getWatchMap(userId uint64) *userWatchMap {
	if watchCacheMgr == nil || userId == 0 {
		return gmap.NewKVMap[uint64, *entity.ShortVideoWatch](true)
	}
	v := watchCacheMgr.GetData(userId, func(ctx context.Context) (interface{}, error) {
		return loadAllShortVideoWatchFromDB(userId), nil
	})
	m, _ := v.(*userWatchMap)
	if m == nil {
		return gmap.NewKVMap[uint64, *entity.ShortVideoWatch](true)
	}
	return m
}

func loadAllShortVideoWatchFromDB(userId uint64) *userWatchMap {
	m := gmap.NewKVMap[uint64, *entity.ShortVideoWatch](true)
	if userId == 0 {
		return m
	}
	list := make([]*entity.ShortVideoWatch, 0)
	_ = g.DB().Model(string(entity.TbShortVideoWatch)).
		Where("user_id = ?", userId).
		Order("updated_at desc, id desc").
		Limit(watchCacheMaxSize).
		Scan(&list)
	for _, watch := range list {
		if watch != nil && watch.VideoId > 0 {
			m.Set(watch.VideoId, watch)
		}
	}
	return m
}

func GetOneShortVideoWatch(userId, videoId uint64) *entity.ShortVideoWatch {
	if userId == 0 || videoId == 0 {
		return entity.NewShortVideoWatch(userId, videoId)
	}
	m := getWatchMap(userId)
	if watch, ok := m.Search(videoId); ok {
		return watch
	}
	watch := entity.NewShortVideoWatch(userId, videoId)
	m.Set(videoId, watch)
	return watch
}

// GetShortVideoWatch 分页查询用户可继续观看的记录(已上架且免费/已付费,不含作者本人视频)
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

	filtered := filterWatchList(getWatchMap(userId).Values(), userId)
	sortWatchList(filtered)

	start := (pageIndex - 1) * pageSize
	if start >= len(filtered) {
		return make([]*entity.ShortVideoWatch, 0)
	}
	end := start + pageSize
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[start:end]
}

func filterWatchList(rows []*entity.ShortVideoWatch, userId uint64) []*entity.ShortVideoWatch {
	list := make([]*entity.ShortVideoWatch, 0, len(rows))
	for _, watch := range rows {
		if watch == nil || watch.VideoId == 0 {
			continue
		}
		video := GetShortVideoById(watch.VideoId)
		if video == nil || video.Status != entity.ShortVideoStatusOnShelf {
			continue
		}
		if video.AuthorId == userId {
			continue
		}
		if video.IsPaid != entity.ShortVideoPaidYes || watch.PaidTime != nil {
			list = append(list, watch)
		}
	}
	return list
}

func sortWatchList(list []*entity.ShortVideoWatch) {
	sort.Slice(list, func(i, j int) bool {
		a, b := list[i], list[j]
		if a == nil || b == nil {
			return a != nil
		}
		if !a.UpdatedAt.Equal(b.UpdatedAt) {
			return a.UpdatedAt.After(b.UpdatedAt)
		}
		return a.ID > b.ID
	})
}

// DeleteWatchByVideoId 物理删除指定视频的观看/点赞记录,并清理相关缓存
func DeleteWatchByVideoId(videoId uint64) error {
	if videoId == 0 {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbShortVideoWatch)).Where("video_id = ?", videoId).Delete()
	if err != nil {
		return err
	}
	removeWatchFromCacheByVideoId(videoId)
	return nil
}

func removeWatchFromCacheByVideoId(videoId uint64) {
	if watchCacheMgr == nil || videoId == 0 {
		return
	}
	for _, v := range watchCacheMgr.Cache.MustValues(context.Background()) {
		m, ok := v.(*userWatchMap)
		if !ok || m == nil {
			continue
		}
		m.Remove(videoId)
	}
}
