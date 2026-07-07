package calldao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var orderCacheMgr *cache.CacheMgr

// Init 初始化通话相关 DAO 缓存
func Init() {
	orderCacheMgr = cache.NewCacheMgr()
	userCacheMgr = cache.NewCacheMgr()
}

// GetOrderById 按主键查询通话订单(走缓存)
func GetOrderById(id uint64) *entity.CallOrder {
	if id == 0 || orderCacheMgr == nil {
		return nil
	}
	v := orderCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var ret entity.CallOrder
		err = g.DB().Model(string(entity.TbCallOrder)).WherePri(id).Scan(&ret)
		if err != nil || ret.ID == 0 {
			return nil, err
		}
		return &ret, nil
	})
	if v == nil {
		return nil
	}
	o, _ := v.(*entity.CallOrder)
	return o
}

// AddOrderToCache 新建通话订单后写入缓存
func AddOrderToCache(o *entity.CallOrder) {
	if o == nil || orderCacheMgr == nil {
		return
	}
	orderCacheMgr.FlushCache(o.ID, o)
}

// FlushOrderCache 通话订单变更后刷新缓存
func FlushOrderCache(o *entity.CallOrder) {
	AddOrderToCache(o)
}
