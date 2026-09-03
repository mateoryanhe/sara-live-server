package shortvideodao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity/shortvideo"
)

type WatchListFilter struct {
	UserId    uint64
	StartTime int64
	EndTime   int64
	PageIndex int
	PageSize  int
	OnlyPaid  bool
}

func GetList(f *WatchListFilter) (int, []*entity.ShortVideoWatch) {
	if f == nil {
		return 0, make([]*entity.ShortVideoWatch, 0)
	}
	if f.PageIndex <= 0 {
		f.PageIndex = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 20
	}

	ctx := gctx.New()
	list := make([]*entity.ShortVideoWatch, 0)
	m := g.Model(string(entity.TbShortVideoWatch)).Ctx(ctx)

	if f.UserId > 0 {
		m = m.Where(string(entity.ShortVideoWatchUserId)+" = ?", f.UserId)
	}
	if f.OnlyPaid {
		m = m.Where(string(entity.ShortVideoWatchPaidTime) + " IS NOT NULL")
	}
	timeCol := "updated_at"
	orderBy := "updated_at desc, id desc"
	if f.OnlyPaid {
		timeCol = string(entity.ShortVideoWatchPaidTime)
		orderBy = string(entity.ShortVideoWatchPaidTime) + " desc, id desc"
	}
	if f.StartTime > 0 {
		m = m.Where(timeCol+" >= ?", time.Unix(f.StartTime, 0))
	}
	if f.EndTime > 0 {
		m = m.Where(timeCol+" <= ?", time.Unix(f.EndTime, 0))
	}

	total, err := m.Clone().Count()
	if err != nil {
		return 0, list
	}
	_ = m.Clone().Order(orderBy).
		Limit(f.PageSize).Offset((f.PageIndex - 1) * f.PageSize).
		Scan(&list)
	return total, list
}
