package shortvideodao

import (
	"context"
	"xr-game-server/core/lambda"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
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
		newTemp := make([]*entity.ShortVideoWatch, 0)
		newTemp = append(newTemp, one)
		all = append(newTemp, all...)
		watchCacheMgr.FlushCache(userId, all)
		return one
	}
	return one
}

func GetShortVideoWatch(userId uint64) []*entity.ShortVideoWatch {
	v := watchCacheMgr.GetData(userId, func(ctx context.Context) (value interface{}, err error) {
		return loadShortVideoWatchByUserId(userId), nil
	})
	return v.([]*entity.ShortVideoWatch)
}

func loadShortVideoWatchByUserId(userId uint64) []*entity.ShortVideoWatch {
	watch := make([]*entity.ShortVideoWatch, 0)
	_ = g.Model(string(entity.TbShortVideoWatch)).Where("user_id = ?", userId).OrderDesc("updated_at").Scan(&watch)
	return watch
}

func removeVideoFromWatchCache(videoId uint64) {
	if videoId == 0 || watchCacheMgr == nil {
		return
	}
	ctx := gctx.New()
	keys, err := watchCacheMgr.Cache.Keys(ctx)
	if err != nil || len(keys) == 0 {
		return
	}
	for _, key := range keys {
		value, _ := watchCacheMgr.Cache.Get(ctx, key)
		if value == nil || value.IsNil() {
			continue
		}
		all, ok := value.Val().([]*entity.ShortVideoWatch)
		if !ok || len(all) == 0 {
			continue
		}
		filtered := lambda.Filter(all, func(watch *entity.ShortVideoWatch) bool {
			return watch.VideoId != videoId
		})
		if len(filtered) == len(all) {
			continue
		}
		watchCacheMgr.FlushCache(key, filtered)
	}
}

// DeleteWatchByVideoId 删除指定视频的观看/点赞记录,并刷新已加载的用户观看缓存
func DeleteWatchByVideoId(videoId uint64) error {
	if videoId == 0 {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbShortVideoWatch)).Where("video_id = ?", videoId).Delete()
	if err != nil {
		return err
	}
	removeVideoFromWatchCache(videoId)
	return nil
}
