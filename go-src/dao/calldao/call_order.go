package calldao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
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

// GetOrderFromCache 仅从内存缓存读取通话订单,未命中不访问数据库
func GetOrderFromCache(id uint64) *entity.CallOrder {
	if id == 0 || orderCacheMgr == nil {
		return nil
	}
	v := orderCacheMgr.GetFromCache(id)
	if v == nil {
		return nil
	}
	o, _ := v.(*entity.CallOrder)
	return o
}

// LoadUnclosedOrders 启动时加载未关闭通话订单(仅初始化调用一次)
func LoadUnclosedOrders() []*entity.CallOrder {
	ctx := gctx.New()
	orders := make([]*entity.CallOrder, 0)
	err := g.Model(string(entity.TbCallOrder)).Ctx(ctx).
		WhereIn(string(entity.CallOrderStatusCol), g.Slice{
			entity.CallOrderStatusCalling,
			entity.CallOrderStatusAccepted,
			entity.CallOrderStatusInCall,
		}).
		Scan(&orders)
	if err != nil {
		g.Log().Errorf(ctx, "load unclosed call orders failed: %v", err)
		return nil
	}
	return orders
}
