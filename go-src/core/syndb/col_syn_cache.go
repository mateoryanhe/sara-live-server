package syndb

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/gogf/gf/v2/container/gqueue"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gutil"
	"xr-game-server/constants/common"
	"xr-game-server/constants/db"
	"xr-game-server/core/xrlog"
)

const (
	// Max 单表单批最大落库行数(内存换批量,减少 SQL 次数)
	Max = 500
	// maxQuickFlushPerTick quick 每轮最多落库次数(独立配额,保障约 1 秒内入库)
	maxQuickFlushPerTick = 10
	// maxLazyFlushPerTick lazy 每轮最多落库次数(独立配额,不抢占 quick)
	maxLazyFlushPerTick = 4
	// CloseTime 服务器关闭时缩短同步周期
	CloseTime = 50 * time.Millisecond
	// syndbPushSlowLogMs Push 超过该阈值(毫秒)时记录慢日志
	syndbPushSlowLogMs = int64(5)
)

type ColData struct {
	IdVal  any
	ColVal any
}

// ColSynCache 内存同步数据库工具
type ColSynCache struct {
	// 无界队列,Push 写入链表不阻塞
	DataQueue *gqueue.TQueue[*ColData]
	// Pending 待落库数据(按主键去重,只保留最新值)
	Pending map[any]*ColData
	// 同步频率
	Period time.Duration
	// 本批 Pending 首次写入时间
	LastTime time.Time
	// 表名
	TbName string
	// 列名
	ColName string
	// 主键列名
	IdName string
}

func newColSynCache(tbName string, tbCol string, period time.Duration) *ColSynCache {
	return &ColSynCache{
		TbName:    tbName,
		ColName:   tbCol,
		DataQueue: gqueue.NewTQueue[*ColData](),
		Pending:   make(map[any]*ColData),
		Period:    period,
		LastTime:  time.Now(),
		IdName:    string(db.IdName),
	}
}

// AddDataToLazyChan 加入变更数据到延迟缓存区
func AddDataToLazyChan(tbName db.TbName, tbCol db.TbCol, colData *ColData) {
	addDataToQueue(lazyMap, tbName, tbCol, colData)
}

// AddDataToQuickChan 加入变更数据到快速缓存区
func AddDataToQuickChan(tbName db.TbName, tbCol db.TbCol, colData *ColData) {
	addDataToQueue(quickMap, tbName, tbCol, colData)
}

func addDataToQueue(cacheMap map[string]*ColSynCache, tbName db.TbName, tbCol db.TbCol, colData *ColData) {
	key := string(tbName) + ":" + string(tbCol)
	start := time.Now()
	cacheMap[key].DataQueue.Push(colData)
	if costMs := time.Since(start).Milliseconds(); costMs >= syndbPushSlowLogMs {
		g.Log().Warningf(gctx.New(), "syndb Push慢,table=%v,col=%v,耗时=%vms,id=%v", tbName, tbCol, costMs, colData.IdVal)
	}
}

// SysExit 调整全部数据库同步组件同步时间到最小
func SysExit(sig os.Signal) {
	g.Log().Warning(gctx.New(), "准备关机,开始强制同步内存数据到数据库")

	for true {
		for _, val := range lazyMap {
			gutil.TryCatch(gctx.New(), func(ctx context.Context) {
				val.ChangeSynTime()
			}, func(ctx context.Context, exception error) {
			})
		}
		for _, quick := range quickMap {
			gutil.TryCatch(gctx.New(), func(ctx context.Context) {
				quick.ChangeSynTime()
			}, func(ctx context.Context, exception error) {
			})
		}
	}
}

func consume(ctx context.Context) {
	pullCaches(quickMap)
	pullCaches(lazyMap)
	start := time.Now()
	quickRows := flushCaches(quickMap, maxQuickFlushPerTick)
	lazyRows := flushCaches(lazyMap, maxLazyFlushPerTick)
	totalRows := quickRows + lazyRows
	if totalRows > 0 {
		g.Log().Infof(ctx, "syndb刷盘成功,quick=%v,lazy=%v,total=%v,costMs=%v",
			quickRows, lazyRows, totalRows, time.Since(start).Milliseconds())
	}
}

func pullCaches(cacheMap map[string]*ColSynCache) {
	for _, colCache := range cacheMap {
		if colCache.isIdle() {
			continue
		}
		safePullData(colCache)
	}
}

func flushCaches(cacheMap map[string]*ColSynCache, flushBudget int) int {
	flushedRows := 0
	for _, colCache := range cacheMap {
		if flushBudget <= 0 {
			break
		}
		if len(colCache.Pending) == 0 {
			continue
		}
		if rows := safeSyn(colCache); rows > 0 {
			flushedRows += rows
			flushBudget--
		}
	}
	return flushedRows
}

func safePullData(colCache *ColSynCache) {
	gutil.TryCatch(gctx.New(), func(ctx context.Context) {
		colCache.PullData()
	}, func(ctx context.Context, exception error) {
		xrlog.ErrorWithErr(ctx, "SynDb", "PullData,table="+colCache.TbName+",col="+colCache.ColName, exception)
	})
}

func safeSyn(colCache *ColSynCache) int {
	var flushedRows int
	gutil.TryCatch(gctx.New(), func(ctx context.Context) {
		var flushed bool
		flushedRows, flushed = colCache.Syn()
		if !flushed {
			flushedRows = 0
		}
	}, func(ctx context.Context, exception error) {
		xrlog.ErrorWithErr(ctx, "SynDb", "Syn,table="+colCache.TbName+",col="+colCache.ColName, exception)
	})
	return flushedRows
}

func (colCache *ColSynCache) isIdle() bool {
	return colCache.DataQueue.Len() == 0 && len(colCache.Pending) == 0
}

// PullData 批量拉取队列数据到 Pending(内存按主键去重)
func (colCache *ColSynCache) PullData() {
	for len(colCache.Pending) < Max {
		select {
		case data := <-colCache.DataQueue.C:
			if len(colCache.Pending) == 0 {
				colCache.LastTime = time.Now()
			}
			colCache.Pending[data.IdVal] = data
		default:
			return
		}
	}
}

func (colCache *ColSynCache) shouldFlush(now time.Time) bool {
	if len(colCache.Pending) == 0 {
		return false
	}
	if len(colCache.Pending) >= Max {
		return true
	}
	return !now.Before(colCache.LastTime.Add(colCache.Period))
}

func (colCache *ColSynCache) Syn() (int, bool) {
	if !colCache.shouldFlush(time.Now()) {
		return 0, false
	}
	return colCache.batchSave()
}

func (colCache *ColSynCache) batchSave() (int, bool) {
	if len(colCache.Pending) == common.Zero {
		return 0, false
	}
	dataMap := make([]map[string]interface{}, 0, len(colCache.Pending))
	for idVal, val := range colCache.Pending {
		dataMap = append(dataMap, g.Map{
			colCache.IdName:  idVal,
			colCache.ColName: val.ColVal,
		})
	}
	rows := len(dataMap)
	_, err := g.DB().Model(colCache.TbName).Data(dataMap).Batch(rows).Save()
	if err == nil {
		colCache.Pending = make(map[any]*ColData)
		return rows, true
	}

	saved := colCache.saveRowsDiscardOnError(dataMap)
	colCache.Pending = make(map[any]*ColData)
	if saved == 0 {
		g.Log().Warningf(gctx.New(), "syndb batchSave失败已全部丢弃,table=%v,col=%v,rows=%v,err=%v",
			colCache.TbName, colCache.ColName, rows, err)
		return 0, true
	}
	if saved < rows {
		g.Log().Warningf(gctx.New(), "syndb batchSave部分落库,table=%v,col=%v,saved=%v,discarded=%v,err=%v",
			colCache.TbName, colCache.ColName, saved, rows-saved, err)
	}
	return saved, true
}

func (colCache *ColSynCache) saveRowsDiscardOnError(dataMap []map[string]interface{}) int {
	saved := 0
	for _, row := range dataMap {
		_, err := g.DB().Model(colCache.TbName).Data(row).Save()
		if err != nil {
			idVal := row[colCache.IdName]
			if isDuplicateKeyErr(err) {
				g.Log().Warningf(gctx.New(), "syndb落库冲突已丢弃,table=%v,col=%v,id=%v,err=%v",
					colCache.TbName, colCache.ColName, idVal, err)
			} else {
				g.Log().Warningf(gctx.New(), "syndb落库失败已丢弃,table=%v,col=%v,id=%v,err=%v",
					colCache.TbName, colCache.ColName, idVal, err)
			}
			continue
		}
		saved++
	}
	return saved
}

func isDuplicateKeyErr(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "1062") || strings.Contains(msg, "Duplicate entry")
}

// ChangeSynTime 变更同步时间
func (colCache *ColSynCache) ChangeSynTime() {
	colCache.Period = CloseTime
}
