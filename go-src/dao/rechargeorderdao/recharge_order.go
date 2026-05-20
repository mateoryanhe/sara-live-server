package rechargeorderdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"time"
	"xr-game-server/entity"
)

// GetById 按主键查询充值订单
func GetById(id uint64) *entity.RechargeOrder {
	var ret entity.RechargeOrder
	err := g.DB().Model(string(entity.TbRechargeOrder)).WherePri(id).Scan(&ret)
	if err != nil {
		return nil
	}
	if ret.ID == 0 {
		return nil
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
	return total, list
}
