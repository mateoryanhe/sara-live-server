package shortvideodao

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var shortVideoStatCacheMgr *cache.CacheMgr

func initShortVideoStatDao() {
	shortVideoStatCacheMgr = cache.NewCacheMgr()
}

// GetStatByVideoId 根据视频ID获取统计数据,不存在则新建内存对象
func GetStatByVideoId(videoId uint64) *entity.ShortVideoStat {
	if videoId == 0 || shortVideoStatCacheMgr == nil {
		return nil
	}
	cacheData := shortVideoStatCacheMgr.GetData(videoId, func(ctx context.Context) (value interface{}, err error) {
		var row *entity.ShortVideoStat
		err = g.Model(string(entity.TbShortVideoStat)).Where(g.Map{
			string(db.IdName): videoId,
		}).Scan(&row)
		if row != nil && row.ID != 0 {
			return row, nil
		}
		return entity.NewShortVideoStat(videoId), nil
	})
	if cacheData == nil {
		return nil
	}
	stat, _ := cacheData.(*entity.ShortVideoStat)
	return stat
}

func DeleteByVideoId(videoId uint64) error {
	if videoId == 0 {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbShortVideoStat)).WherePri(videoId).Delete()
	return err
}
