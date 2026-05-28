package liveroomdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var giftLogCacheMgr *cache.CacheMgr

func initGiftLogDao() {
	giftLogCacheMgr = cache.NewCacheMgr()
}

// GetGiftLogById 按主键查询礼物流水(走缓存);数据库不存在则返回新实例
func GetGiftLogById(id uint64) *entity.LiveGiftLog {
	if id == 0 || giftLogCacheMgr == nil {
		return nil
	}
	v := giftLogCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var ret *entity.LiveGiftLog
		err = g.DB().Model(string(entity.TbLiveGiftLog)).WherePri(id).Scan(&ret)
		if err != nil || ret == nil {
			return entity.NewLiveGiftLog(id), nil
		}
		return ret, nil
	})
	if v == nil {
		return nil
	}
	r, _ := v.(*entity.LiveGiftLog)
	return r
}
