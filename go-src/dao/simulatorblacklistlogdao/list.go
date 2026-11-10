package simulatorblacklistlogdao

import (
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	userentity "xr-game-server/entity/user"
)

type ListFilter struct {
	Key       string
	LoginType string
	StartTime int64
	EndTime   int64
	PageIndex int
	PageSize  int
}

func List(f *ListFilter) (int, []*userentity.SimulatorBlacklistLog) {
	list := make([]*userentity.SimulatorBlacklistLog, 0)
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
	m := g.Model(string(userentity.TbSimulatorBlacklistLog)).Ctx(gctx.New())
	if key := strings.TrimSpace(f.Key); key != "" {
		like := "%" + key + "%"
		m = m.Where("(device_id LIKE ? OR ip LIKE ? OR email LIKE ? OR cpu_model LIKE ?)", like, like, like, like)
	}
	if loginType := strings.TrimSpace(f.LoginType); loginType != "" {
		m = m.Where("login_type = ?", loginType)
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
