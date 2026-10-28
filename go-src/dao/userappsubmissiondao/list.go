package userappsubmissiondao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	userentity "xr-game-server/entity/user"
)

type ListFilter struct {
	UserId    uint64
	StartTime int64
	EndTime   int64
	PageIndex int
	PageSize  int
}

func pageArgs(f *ListFilter) (pageIndex, pageSize int) {
	pageIndex = 1
	pageSize = 20
	if f == nil {
		return pageIndex, pageSize
	}
	if f.PageIndex > 0 {
		pageIndex = f.PageIndex
	}
	if f.PageSize > 0 {
		pageSize = f.PageSize
	}
	return pageIndex, pageSize
}

func ListFeedback(f *ListFilter) (int, []*userentity.UserFeedbackRecord) {
	list := make([]*userentity.UserFeedbackRecord, 0)
	if f == nil {
		return 0, list
	}
	pageIndex, pageSize := pageArgs(f)
	m := g.Model(string(userentity.TbUserFeedbackRecord)).Ctx(gctx.New())
	if f.UserId > 0 {
		m = m.Where("user_id = ?", f.UserId)
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

func ListReport(f *ListFilter) (int, []*userentity.UserReportRecord) {
	list := make([]*userentity.UserReportRecord, 0)
	if f == nil {
		return 0, list
	}
	pageIndex, pageSize := pageArgs(f)
	m := g.Model(string(userentity.TbUserReportRecord)).Ctx(gctx.New())
	if f.UserId > 0 {
		m = m.Where("user_id = ?", f.UserId)
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
