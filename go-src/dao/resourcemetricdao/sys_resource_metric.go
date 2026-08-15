package resourcemetricdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity/sys"
)

// DeleteBefore 删除采样时间早于指定时刻的记录(保留最近数据)
func DeleteBefore(before time.Time) {
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

// ListSince 查询指定时间以来的采样记录(按时间升序)
func ListSince(since time.Time) []*entity.SysResourceMetric {
	ret := make([]*entity.SysResourceMetric, 0)
	_ = g.Model(string(entity.TbSysResourceMetric)).Ctx(gctx.New()).
		Where("recorded_at >= ?", since).
		OrderAsc("recorded_at").
		Scan(&ret)
	return ret
}

const maxTrendPoints = 1000

// ListByTimeRange 查询时间范围内的采样记录(按时间升序,最多 limit 条,取范围内最新数据)
func ListByTimeRange(start, end time.Time, limit int) []*entity.SysResourceMetric {
	if limit <= 0 {
		limit = maxTrendPoints
	}
	if limit > maxTrendPoints {
		limit = maxTrendPoints
	}
	ctx := gctx.New()
	m := g.Model(string(entity.TbSysResourceMetric)).Ctx(ctx)
	if !start.IsZero() {
		m = m.Where("recorded_at >= ?", start)
	}
	if !end.IsZero() {
		m = m.Where("recorded_at <= ?", end)
	}
	ret := make([]*entity.SysResourceMetric, 0)
	_ = m.OrderDesc("recorded_at").Limit(limit).Scan(&ret)
	for i, j := 0, len(ret)-1; i < j; i, j = i+1, j-1 {
		ret[i], ret[j] = ret[j], ret[i]
	}
	return ret
}
