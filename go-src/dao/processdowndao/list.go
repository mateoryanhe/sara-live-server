package processdowndao

import (
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	sysentity "xr-game-server/entity/sys"
)

type ListFilter struct {
	Host      string
	StartTime int64
	EndTime   int64
	PageIndex int
	PageSize  int
}

func List(f *ListFilter) (int, []*sysentity.ProcessDownLog) {
	list := make([]*sysentity.ProcessDownLog, 0)
	if f == nil {
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
	m := g.Model(string(sysentity.TbProcessDownLog)).Ctx(gctx.New())
	if host := strings.TrimSpace(f.Host); host != "" {
		m = m.Where("host LIKE ?", "%"+host+"%")
	}
	if f.StartTime > 0 {
		m = m.Where("detected_at >= ?", time.Unix(f.StartTime, 0))
	}
	if f.EndTime > 0 {
		m = m.Where("detected_at <= ?", time.Unix(f.EndTime, 0))
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
