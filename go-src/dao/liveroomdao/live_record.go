package liveroomdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
)

var liveRecordCacheMgr *cache.CacheMgr

// InitLiveRecordDao 初始化直播记录缓存
func initLiveRecordDao() {
	liveRecordCacheMgr = cache.NewCacheMgr()
	initAppLiveRecordListCache()
	initLiveRecordUserDao()
	initRevenueLogDao()
}

// GetLiveRecordById 按主键查询直播记录(走缓存)
func GetLiveRecordById(id uint64) *entity.LiveRecord {
	if id == 0 || liveRecordCacheMgr == nil {
		return nil
	}
	v := liveRecordCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var ret *entity.LiveRecord
		err = g.DB().Model(string(entity.TbLiveRecord)).WherePri(id).Scan(&ret)
		if err != nil || ret == nil {
			return entity.NewLiveRecord(id), nil
		}
		return ret, nil
	})
	if v == nil {
		return nil
	}
	r, _ := v.(*entity.LiveRecord)
	return r
}

// GetDataFromCache 仅从内存缓存读取直播记录,未命中返回 nil(不查库)
func GetDataFromCache(id uint64) *entity.LiveRecord {
	if id == 0 || liveRecordCacheMgr == nil {
		return nil
	}
	v := liveRecordCacheMgr.GetFromCache(id)
	if v == nil {
		return nil
	}
	r, _ := v.(*entity.LiveRecord)
	return r
}

func resolveLiveRecordTarget(row *entity.LiveRecord) *entity.LiveRecord {
	if row == nil {
		return nil
	}
	if cached := GetDataFromCache(row.ID); cached != nil {
		return cached
	}
	return row
}
