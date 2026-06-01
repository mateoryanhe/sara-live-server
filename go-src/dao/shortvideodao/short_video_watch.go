package shortvideodao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var watchCacheMgr *cache.CacheMgr

func InitShortVideoWatchDao() {
	watchCacheMgr = cache.NewCacheMgr()
}

func GetShortVideoWatchByUserVideo(userId, videoId uint64) *entity.ShortVideoWatch {
	if userId == 0 || videoId == 0 || watchCacheMgr == nil {
		return nil
	}
	id := entity.BuildShortVideoWatchId(userId, videoId)
	v := watchCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var watch *entity.ShortVideoWatch
		_ = g.Model(string(entity.TbShortVideoWatch)).Where("id = ?", id).Scan(&watch)
		if watch == nil || watch.ID == "" {
			return nil, nil
		}
		return watch, nil
	})
	if v == nil {
		return nil
	}
	watch, _ := v.(*entity.ShortVideoWatch)
	if watch == nil || watch.ID == "" {
		return nil
	}
	return watch
}

func SaveToCache(watch *entity.ShortVideoWatch) {
	if watch == nil || watchCacheMgr == nil || watch.ID == "" {
		return
	}
	watchCacheMgr.FlushCache(watch.ID, watch)
}
