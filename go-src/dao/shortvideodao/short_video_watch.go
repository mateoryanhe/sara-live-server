package shortvideodao

import (
	"context"
	"sort"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/shortvideo"
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

// GetUserShortVideoWatchAll 返回用户全部观看记录(按更新时间降序,不做业务过滤)
func GetUserShortVideoWatchAll(userId uint64) []*entity.ShortVideoWatch {
	if userId == 0 {
		return make([]*entity.ShortVideoWatch, 0)
	}
	list := getWatchMap(userId).Values()
	sortWatchList(list)
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
	vals, _ := watchCacheMgr.Cache.Values(context.Background())
	for _, v := range vals {
		m, ok := v.(*userWatchMap)
		if !ok || m == nil {
			continue
		}
		m.Remove(videoId)
	}
}
