package rechargeorderdao

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"strings"
	"time"
	"xr-game-server/core/cache"
	"xr-game-server/entity/recharge"
)

var orderCacheMgr *cache.RowCache[*entity.RechargeOrder]

func InitRechargeOrderDao() {
	orderCacheMgr = cache.NewRowCache[*entity.RechargeOrder]()
}

// GetById 按主键查询充值订单(走缓存)
func GetById(id uint64) *entity.RechargeOrder {
	if id == 0 {
		return nil
	}
	v := orderCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.RechargeOrder, error) {
		var ret entity.RechargeOrder
		err := g.DB().Model(string(entity.TbRechargeOrder)).WherePri(id).Scan(&ret)
		if err != nil || ret.ID == 0 {
			return nil, err
		}
		return &ret, nil
	})
	return v
}

// AddOrderToCache 新建订单后写入缓存
func AddOrderToCache(o *entity.RechargeOrder) {
	if o == nil {
		return
	}
	orderCacheMgr.PublishRow(gctx.New(), o.ID, o)
}

// FlushOrderCache 订单变更后刷新缓存
func FlushOrderCache(o *entity.RechargeOrder) {
	if o == nil {
		return
	}
	orderCacheMgr.PublishRow(gctx.New(), o.ID, o)
}

// GetFromCacheById 仅从内存缓存读取订单,未命中返回 nil(不查库)
func GetFromCacheById(id uint64) *entity.RechargeOrder {
	if id == 0 || orderCacheMgr == nil {
		return nil
	}
	v, _ := orderCacheMgr.GetRowCached(gctx.New(), id)
	return v
}

// mergeOrderFromCache 列表项优先使用内存缓存中的最新订单数据
func mergeOrderFromCache(o *entity.RechargeOrder) *entity.RechargeOrder {
	if o == nil {
		return nil
	}
	if cached := GetFromCacheById(o.ID); cached != nil {
		return cached
	}
	return o
}

// RemoveOrderCache 移除指定订单缓存
func RemoveOrderCache(orderId uint64) {
	if orderCacheMgr == nil || orderId == 0 {
		return
	}
	orderCacheMgr.RemoveRow(gctx.New(), orderId)
}

// GetByThirdOrderId 按第三方订单号查询(Google orderId 幂等)
func GetByThirdOrderId(thirdOrderId string) *entity.RechargeOrder {
	thirdOrderId = strings.TrimSpace(thirdOrderId)
	if thirdOrderId == "" {
		return nil
	}
	ctx := gctx.New()
	var ret entity.RechargeOrder
	err := g.DB().Model(string(entity.TbRechargeOrder)).Ctx(ctx).
		Where(string(entity.RechargeOrderThirdOrderId)+" = ?", thirdOrderId).
		Scan(&ret)
	if err != nil || ret.ID == 0 {
		return nil
	}
	if cached := GetFromCacheById(ret.ID); cached != nil {
		return cached
	}
	return &ret
}

// ListByUserId App 端按用户分页查询(按 id 倒序,可按状态过滤)
// statusVal < 0 表示不过滤状态;其他值为数据库中的实际 status 枚举值
func ListByUserId(userId uint64, statusVal int, pageIndex, pageSize int) (int, []*entity.RechargeOrder) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	ctx := gctx.New()
	list := make([]*entity.RechargeOrder, 0)
	m := g.Model(string(entity.TbRechargeOrder)).Ctx(ctx).
		Where(string(entity.RechargeOrderUserId)+" = ?", userId)
	if statusVal >= 0 {
		m = m.Where(string(entity.RechargeOrderStatus)+" = ?", statusVal)
	}
	total, err := m.Clone().Count()
	if err != nil {
		return 0, list
	}
	_ = m.Clone().Order("id desc").
		Limit(pageSize).Offset((pageIndex - 1) * pageSize).
		Scan(&list)
	return total, list
}

// CMSListFilter CMS 列表过滤条件
type CMSListFilter struct {
	UserId    uint64
	OrderId   uint64
	StatusVal int   // <0=不过滤;>=0 则为数据库中的实际 status 枚举值
	Source    int   // 0=不过滤
	StartTime int64 // 秒,0=不过滤
	EndTime   int64 // 秒,0=不过滤
	PageIndex int
	PageSize  int
}

// CMSList CMS 端分页查询(按 id 倒序,多维度过滤)
func CMSList(f *CMSListFilter) (int, []*entity.RechargeOrder) {
	if f.PageIndex <= 0 {
		f.PageIndex = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 20
	}
	ctx := gctx.New()
	list := make([]*entity.RechargeOrder, 0)
	m := g.Model(string(entity.TbRechargeOrder)).Ctx(ctx)

	if f.OrderId > 0 {
		m = m.WherePri(f.OrderId)
	}
	if f.UserId > 0 {
		m = m.Where(string(entity.RechargeOrderUserId)+" = ?", f.UserId)
	}
	if f.StatusVal >= 0 {
		m = m.Where(string(entity.RechargeOrderStatus)+" = ?", f.StatusVal)
	}
	if f.Source > 0 {
		m = m.Where(string(entity.RechargeOrderSource)+" = ?", f.Source)
	}
	if f.StartTime > 0 {
		m = m.Where("created_at >= ?", time.Unix(f.StartTime, 0))
	}
	if f.EndTime > 0 {
		m = m.Where("created_at <= ?", time.Unix(f.EndTime, 0))
	}

	total, err := m.Clone().Count()
	if err != nil {
		return 0, list
	}
	_ = m.Clone().Order("id desc").
		Limit(f.PageSize).Offset((f.PageIndex - 1) * f.PageSize).
		Scan(&list)
	for i, row := range list {
		list[i] = mergeOrderFromCache(row)
	}
	return total, list
}

// LoadPendingOrders 启动时加载待支付充值订单
func LoadPendingOrders() []*entity.RechargeOrder {
	ctx := gctx.New()
	orders := make([]*entity.RechargeOrder, 0)
	err := g.Model(string(entity.TbRechargeOrder)).Ctx(ctx).
		Where(string(entity.RechargeOrderStatus)+" = ?", entity.RechargeOrderStatusPending).
		Scan(&orders)
	if err != nil {
		g.Log().Errorf(ctx, "load pending recharge orders failed: %v", err)
		return nil
	}
	return orders
}

// SumCompletedRechargeGoldByUserId 统计用户已完成充值订单累计到账金币
func SumCompletedRechargeGoldByUserId(userId uint64) float64 {
	if userId == 0 {
		return 0
	}
	ctx := gctx.New()
	val, err := g.Model(string(entity.TbRechargeOrder)).Ctx(ctx).
		Where(string(entity.RechargeOrderUserId)+" = ?", userId).
		Where(string(entity.RechargeOrderStatus)+" = ?", entity.RechargeOrderStatusCompleted).
		Sum(string(entity.RechargeOrderGold))
	if err != nil {
		g.Log().Errorf(ctx, "sum completed recharge gold failed userId=%d: %v", userId, err)
		return 0
	}
	return val
}
