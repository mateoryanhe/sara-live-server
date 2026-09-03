package resourcemetricdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity/sys"
)

const maxTrendPoints = 10000

// DeleteFineBefore 删除细采样过期记录
func DeleteFineBefore(before time.Time) {
	deleteBefore(string(entity.TbSysResourceMetric), before, "细采样")
}

// DeleteAggBefore 删除粗采样过期记录
func DeleteAggBefore(before time.Time) {
	deleteBefore(string(entity.TbSysResourceMetricAgg), before, "粗采样")
}

// DeleteBefore 兼容旧调用,等同删除细采样
func DeleteBefore(before time.Time) {
	DeleteFineBefore(before)
}

func deleteBefore(table string, before time.Time, label string) {
	result, err := g.Model(table).Ctx(gctx.New()).
		Where("recorded_at < ?", before).
		Delete()
	if err != nil {
		g.Log().Errorf(gctx.New(), "删除过期资源%s记录失败,before=%v,err=%v", label, before, err)
		return
	}
	rows, _ := result.RowsAffected()
	if rows > 0 {
		g.Log().Infof(gctx.New(), "删除过期资源%s记录成功,before=%v,rows=%v", label, before, rows)
	}
}

// ListSince 查询指定时间以来的细采样记录(按时间升序)
func ListSince(since time.Time) []*entity.SysResourceMetric {
	ret := make([]*entity.SysResourceMetric, 0)
	_ = g.Model(string(entity.TbSysResourceMetric)).Ctx(gctx.New()).
		Where("recorded_at >= ?", since).
		OrderAsc("recorded_at").
		Scan(&ret)
	return ret
}

// ListFineByTimeRange 查询细采样时间范围(按时间升序,最多 limit 条,取范围内最新数据)
func ListFineByTimeRange(start, end time.Time, limit int) []*entity.SysResourceMetric {
	return listByTimeRange[entity.SysResourceMetric](string(entity.TbSysResourceMetric), start, end, limit)
}

// ListAggByTimeRange 查询粗采样时间范围
func ListAggByTimeRange(start, end time.Time, limit int) []*entity.SysResourceMetricAgg {
	return listByTimeRange[entity.SysResourceMetricAgg](string(entity.TbSysResourceMetricAgg), start, end, limit)
}

// ListByTimeRange 兼容旧调用,查细采样
func ListByTimeRange(start, end time.Time, limit int) []*entity.SysResourceMetric {
	return ListFineByTimeRange(start, end, limit)
}

func listByTimeRange[T any](table string, start, end time.Time, limit int) []*T {
	if limit <= 0 {
		limit = maxTrendPoints
	}
	if limit > maxTrendPoints {
		limit = maxTrendPoints
	}
	ctx := gctx.New()
	m := g.Model(table).Ctx(ctx)
	if !start.IsZero() {
		m = m.Where("recorded_at >= ?", start)
	}
	if !end.IsZero() {
		m = m.Where("recorded_at <= ?", end)
	}
	ret := make([]*T, 0)
	_ = m.OrderDesc("recorded_at").Limit(limit).Scan(&ret)
	for i, j := 0, len(ret)-1; i < j; i, j = i+1, j-1 {
		ret[i], ret[j] = ret[j], ret[i]
	}
	return ret
}

// HasAggAt 粗采样桶是否已存在
func HasAggAt(bucketStart time.Time) bool {
	count, err := g.Model(string(entity.TbSysResourceMetricAgg)).Ctx(gctx.New()).
		Where("recorded_at = ?", bucketStart).
		Count()
	return err == nil && count > 0
}

// InsertAgg 直写粗采样行
func InsertAgg(row *entity.SysResourceMetricAgg) error {
	if row == nil {
		return nil
	}
	_, err := g.Model(string(entity.TbSysResourceMetricAgg)).Ctx(gctx.New()).Data(row).Insert()
	return err
}

// ListFineInBucket 取某聚合桶内全部细采样(桶为左闭右开)
func ListFineInBucket(bucketStart, bucketEnd time.Time) []*entity.SysResourceMetric {
	ret := make([]*entity.SysResourceMetric, 0)
	_ = g.Model(string(entity.TbSysResourceMetric)).Ctx(gctx.New()).
		Where("recorded_at >= ?", bucketStart).
		Where("recorded_at < ?", bucketEnd).
		OrderAsc("recorded_at").
		Scan(&ret)
	return ret
}
