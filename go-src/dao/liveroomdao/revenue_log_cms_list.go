package liveroomdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity"
)

// RevenueLogCMSListFilter CMS直播收益流水查询条件
type RevenueLogCMSListFilter struct {
	ReceiverId  uint64
	RevenueType uint8
	StartTime   int64
	EndTime     int64
	PageIndex   int
	PageSize    int
}

// RevenueLogCMSList CMS分页查询直播收益流水(按ID倒序)
func RevenueLogCMSList(f *RevenueLogCMSListFilter) (int, []*entity.LiveRevenueLog) {
	list := make([]*entity.LiveRevenueLog, 0)
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
	m := g.Model(string(entity.TbLiveRevenueLog)).Ctx(ctx)
	if f.ReceiverId > 0 {
		m = m.Where(string(entity.LiveRevenueLogReceiverId)+" = ?", f.ReceiverId)
	}
	if f.RevenueType > 0 {
		m = m.Where(string(entity.LiveRevenueLogRevenueType)+" = ?", f.RevenueType)
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
