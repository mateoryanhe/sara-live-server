package liveroomdao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
)

var liveRecordCacheMgr *cache.RowCache[*entity.LiveRecord]

// InitLiveRecordDao 初始化直播记录缓存
func initLiveRecordDao() {
	liveRecordCacheMgr = cache.NewRowCache[*entity.LiveRecord]()
	initAppLiveRecordListCache()
	initLiveRecordUserDao()
	initRevenueLogDao()
}

// GetLiveRecordById 按主键查询直播记录(走缓存)
func GetLiveRecordById(id uint64) *entity.LiveRecord {
	if id == 0 || liveRecordCacheMgr == nil {
		return nil
	}
	v := liveRecordCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.LiveRecord, error) {
		var ret *entity.LiveRecord
		err := g.DB().Model(string(entity.TbLiveRecord)).WherePri(id).Scan(&ret)
		if err != nil || ret == nil {
			return entity.NewLiveRecord(id), nil
		}
		return ret, nil
	})
	return v
}

// GetDataFromCache 仅从内存缓存读取直播记录,未命中返回 nil(不查库)
func GetDataFromCache(id uint64) *entity.LiveRecord {
	if id == 0 || liveRecordCacheMgr == nil {
		return nil
	}
	v, _ := liveRecordCacheMgr.GetRowCached(gctx.New(), id)
	return v
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
