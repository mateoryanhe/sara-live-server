package liveroomdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var revenueLogCacheMgr *cache.CacheMgr

func initRevenueLogDao() {
	revenueLogCacheMgr = cache.NewCacheMgr()
}

// GetRevenueLogById 按主键查询直播收益流水(走缓存);数据库不存在则返回新实例
func GetRevenueLogById(id uint64) *entity.LiveRevenueLog {
	if id == 0 || revenueLogCacheMgr == nil {
		return nil
	}
	v := revenueLogCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var ret *entity.LiveRevenueLog
		err = g.DB().Model(string(entity.TbLiveRevenueLog)).WherePri(id).Scan(&ret)
		if err != nil || ret == nil {
			return entity.NewLiveRevenueLog(id), nil
		}
		return ret, nil
	})
	if v == nil {
		return nil
	}
	r, _ := v.(*entity.LiveRevenueLog)
	return r
}
