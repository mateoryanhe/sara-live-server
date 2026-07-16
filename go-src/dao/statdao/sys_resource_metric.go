package statdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity"
)

// DeleteSysResourceMetricsBefore 删除采样时间早于指定时刻的记录(保留最近数据)
func DeleteSysResourceMetricsBefore(before time.Time) {
	result, err := g.Model(string(entity.TbSysResourceMetric)).Ctx(gctx.New()).
		Where("recorded_at < ?", before).
		Delete()
	if err != nil {
		g.Log().Errorf(gctx.New(), "删除过期资源采样记录失败,before=%v,err=%v", before, err)
		return
	}
	rows, _ := result.RowsAffected()
	if rows > 0 {
		g.Log().Infof(gctx.New(), "删除过期资源采样记录成功,before=%v,rows=%v", before, rows)
	}
}

// ListSysResourceMetricsSince 查询指定时间以来的采样记录(按时间升序)
func ListSysResourceMetricsSince(since time.Time) []*entity.SysResourceMetric {
	ret := make([]*entity.SysResourceMetric, 0)
	_ = g.Model(string(entity.TbSysResourceMetric)).Ctx(gctx.New()).
		Where("recorded_at >= ?", since).
		OrderAsc("recorded_at").
		Scan(&ret)
	return ret
}
