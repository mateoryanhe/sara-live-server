package calldao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity/call"
)

// CallOrderCMSListFilter CMS通话订单查询条件
type CallOrderCMSListFilter struct {
	CallerId    uint64
	ReceiverId  uint64
	ReceiverIds []uint64
	Source      uint8
	Status      uint8
	CallType    uint8
	StartTime   int64
	EndTime     int64
	PageIndex   int
	PageSize    int
}

// CallOrderCMSList CMS分页查询通话订单(按ID倒序)
func CallOrderCMSList(f *CallOrderCMSListFilter) (int, []*entity.CallOrder) {
	list := make([]*entity.CallOrder, 0)
	if f == nil {
		return 0, list
	}
	if f.PageIndex <= 0 {
		f.PageIndex = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 20
	}
	ctx := gctx.New()
	m := g.Model(string(entity.TbCallOrder)).Ctx(ctx)
	if f.CallerId > 0 {
		m = m.Where(string(entity.CallOrderCallerId)+" = ?", f.CallerId)
	}
	if len(f.ReceiverIds) > 0 {
		m = m.Where(string(entity.CallOrderReceiverId)+" IN (?)", f.ReceiverIds)
	} else if f.ReceiverId > 0 {
		m = m.Where(string(entity.CallOrderReceiverId)+" = ?", f.ReceiverId)
	}
	if f.Source > 0 {
		m = m.Where(string(entity.CallOrderSource)+" = ?", f.Source)
	}
	if f.Status > 0 {
		m = m.Where(string(entity.CallOrderStatusCol)+" = ?", f.Status)
	}
	if f.CallType > 0 {
		m = m.Where(string(entity.CallOrderCallType)+" = ?", f.CallType)
	}
	if f.StartTime > 0 {
		m = m.Where(string(entity.CallOrderCallStartTime)+" >= ?", time.Unix(f.StartTime, 0))
	}
	if f.EndTime > 0 {
		m = m.Where(string(entity.CallOrderCallStartTime)+" <= ?", time.Unix(f.EndTime, 0))
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
