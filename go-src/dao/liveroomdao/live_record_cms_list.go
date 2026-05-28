package liveroomdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity"
)

// LiveRecordCMSListFilter CMS直播记录查询条件
type LiveRecordCMSListFilter struct {
	AnchorId  uint64
	StartTime int64
	EndTime   int64
	PageIndex int
	PageSize  int
}

// LiveRecordCMSList CMS分页查询直播记录(按ID倒序)
func LiveRecordCMSList(f *LiveRecordCMSListFilter) (int, []*entity.LiveRecord) {
	list := make([]*entity.LiveRecord, 0)
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
	m := g.Model(string(entity.TbLiveRecord)).Ctx(ctx)
	if f.AnchorId > 0 {
		m = m.Where(string(entity.LiveRecordAnchorId)+" = ?", f.AnchorId)
	}
	if f.StartTime > 0 {
		m = m.Where(string(entity.LiveRecordStartTime)+" >= ?", time.Unix(f.StartTime, 0))
	}
	if f.EndTime > 0 {
		m = m.Where(string(entity.LiveRecordStartTime)+" <= ?", time.Unix(f.EndTime, 0))
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
