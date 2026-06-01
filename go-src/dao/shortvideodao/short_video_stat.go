package shortvideodao

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var shortVideoStatCacheMgr *cache.CacheMgr

func InitShortVideoStatDao() {
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

// GetMapByVideoIds 批量获取统计数据,key为视频ID,缺失项视为0
func GetMapByVideoIds(videoIds []uint64) map[uint64]*entity.ShortVideoStat {
	ret := make(map[uint64]*entity.ShortVideoStat, len(videoIds))
	if len(videoIds) == 0 {
		return ret
	}
	rows := make([]*entity.ShortVideoStat, 0)
	_ = g.Model(string(entity.TbShortVideoStat)).WhereIn(string(db.IdName), videoIds).Scan(&rows)
	for _, row := range rows {
		if row != nil && row.ID != 0 {
			ret[row.ID] = row
		}
	}
	for _, id := range videoIds {
		if id == 0 {
			continue
		}
		if _, ok := ret[id]; !ok {
			ret[id] = entity.NewShortVideoStat(id)
		}
	}
	return ret
}

func DeleteByVideoId(videoId uint64) error {
	if videoId == 0 {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbShortVideoStat)).WherePri(videoId).Delete()
	return err
}
