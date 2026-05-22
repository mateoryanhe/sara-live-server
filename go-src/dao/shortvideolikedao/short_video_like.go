package shortvideolikedao

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var likeCacheMgr *cache.CacheMgr

func InitShortVideoLikeDao() {
	likeCacheMgr = cache.NewCacheMgr()
}

func GetById(id string) *entity.ShortVideoLike {
	v := likeCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var like *entity.ShortVideoLike
		_ = g.Model(string(entity.TbShortVideoLike)).Where("id = ?", id).Scan(&like)
		return like, nil
	})
	if v == nil {
		return nil
	}
	like, _ := v.(*entity.ShortVideoLike)
	return like
}

func GetByUserVideo(userId, videoId uint64) *entity.ShortVideoLike {
	return GetById(entity.BuildShortVideoLikeId(userId, videoId))
}

func AddLikeToCache(like *entity.ShortVideoLike) {
	if like == nil {
		return
	}
	likeCacheMgr.FlushCache(like.ID, like)
}
