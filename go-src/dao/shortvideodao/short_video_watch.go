package shortvideodao

import (
	"context"
	"xr-game-server/core/lambda"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var watchCacheMgr *cache.CacheMgr

func initShortVideoWatchDao() {
	watchCacheMgr = cache.NewCacheMgr()
}

func GetOneShortVideoWatch(userId, videoId uint64) *entity.ShortVideoWatch {
	all := GetShortVideoWatch(userId)
	one, ok := lambda.Find(all, func(watch *entity.ShortVideoWatch) bool {
		return watch.VideoId == videoId
	})
	if !ok {
		one = entity.NewShortVideoWatch(userId, videoId)
		all = append(all, one)
		watchCacheMgr.FlushCache(userId, one)
		return one
	}
	return one
}

func GetShortVideoWatch(userId uint64) []*entity.ShortVideoWatch {
	v := watchCacheMgr.GetData(userId, func(ctx context.Context) (value interface{}, err error) {
		watch := make([]*entity.ShortVideoWatch, 0)
		_ = g.Model(string(entity.TbShortVideoWatch)).Where("user_id = ?", userId).Scan(&watch)
		return watch, nil
	})
	return v.([]*entity.ShortVideoWatch)
}
