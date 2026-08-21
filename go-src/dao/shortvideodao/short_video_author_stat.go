package shortvideodao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/shortvideo"
)

var shortVideoAuthorStatCacheMgr *cache.CacheMgr

func initShortVideoAuthorStatDao() {
	shortVideoAuthorStatCacheMgr = cache.NewPermanentCacheMgr()
}

// GetAuthorStatByAuthorId 根据作者ID获取统计数据,不存在则新建内存对象
func GetAuthorStatByAuthorId(authorId uint64) *entity.ShortVideoAuthorStat {
	if authorId == 0 || shortVideoAuthorStatCacheMgr == nil {
		return nil
	}
	cacheData := shortVideoAuthorStatCacheMgr.GetData(authorId, func(ctx context.Context) (value interface{}, err error) {
		var row *entity.ShortVideoAuthorStat
		err = g.Model(string(entity.TbShortVideoAuthorStat)).Where(g.Map{
			string(db.IdName): authorId,
		}).Scan(&row)
		if row != nil && row.ID != 0 {
			return row, nil
		}
		return entity.NewShortVideoAuthorStat(authorId), nil
	})
	if cacheData == nil {
		return nil
	}
	stat, _ := cacheData.(*entity.ShortVideoAuthorStat)
	return stat
}

// GetAuthorStatFromCache 仅从内存缓存读取作者统计数据,未命中不访问数据库
func GetAuthorStatFromCache(authorId uint64) *entity.ShortVideoAuthorStat {
	if authorId == 0 || shortVideoAuthorStatCacheMgr == nil {
		return nil
	}
	v := shortVideoAuthorStatCacheMgr.GetFromCache(authorId)
	if v == nil {
		return nil
	}
	stat, _ := v.(*entity.ShortVideoAuthorStat)
	return stat
}
