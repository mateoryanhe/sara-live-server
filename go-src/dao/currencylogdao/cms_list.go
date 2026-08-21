package currencylogdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity/user"
)

// CMSListFilter CMS货币流水查询条件
type CMSListFilter struct {
	UserId       uint64
	CurrencyType uint8
	StartTime    int64
	EndTime      int64
	PageIndex    int
	PageSize     int
}

// CMSList CMS分页查询货币流水(按ID倒序)
func CMSList(f *CMSListFilter) (int, []*entity.CurrencyLog) {
	list := make([]*entity.CurrencyLog, 0)
	if f == nil || f.CurrencyType == 0 {
		return 0, list
	}
	pageIndex := f.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}
	pageSize := f.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	ctx := gctx.New()
	m := g.Model(string(entity.TbCurrencyLog)).Ctx(ctx).
		Where(string(entity.CurrencyLogType)+" = ?", f.CurrencyType)
	if f.UserId > 0 {
		m = m.Where(string(entity.CurrencyLogUserId)+" = ?", f.UserId)
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
		Limit(pageSize).Offset((pageIndex - 1) * pageSize).
		Scan(&list)
	return total, list
}
