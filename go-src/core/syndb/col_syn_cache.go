package syndb

import (
	"context"
	"fmt"
	"os"
	"sort"
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
	// CloseTime 服务器关闭时同步间隔
	CloseTime = 50 * time.Millisecond
	// syndbPushSlowLogMs Push 超过该阈值(毫秒)时记录慢日志
	syndbPushSlowLogMs = int64(5)
)

const (
	flushReasonIdle     = "cpu_idle"
	flushReasonForce    = "force_wait"
	flushReasonShutdown = "shutdown"
)

type ColData struct {
	IdVal  any
	ColVal any
}

// ColSynCache 内存同步数据库工具
type ColSynCache struct {
	DataQueue        *gqueue.TQueue[*ColData]
	Pending          map[any]*ColData
	firstPendingTime time.Time
	lastFlushTime    time.Time
	TbName           string
	ColName          string
	IdName           string
}

type flushCandidate struct {
	cache  *ColSynCache
	reason string
	waitMs int64
}

type queueFlushDetail struct {
	table  string
	col    string
	rows   int
	reason string
	waitMs int64
}

func newColSynCache(tbName string, tbCol string) *ColSynCache {
	now := time.Now()
	return &ColSynCache{
		TbName:           tbName,
		ColName:          tbCol,
		DataQueue:        gqueue.NewTQueue[*ColData](),
		Pending:          make(map[any]*ColData),
		firstPendingTime: time.Time{},
		lastFlushTime:    now,
		IdName:           string(db.IdName),
	}
}

// AddData 加入变更数据到缓冲队列.
func AddData(tbName db.TbName, tbCol db.TbCol, colData *ColData) {
	key := cacheKey(tbName, tbCol)
	cache, ok := synCacheMap[key]
	if !ok || cache == nil {
		return
	}
	start := time.Now()
	cache.DataQueue.Push(colData)
	if costMs := time.Since(start).Milliseconds(); costMs >= syndbPushSlowLogMs {
		xrlog.DetailLog.Warningf(gctx.New(), "syndb Push慢,table=%v,col=%v,耗时=%vms,id=%v", tbName, tbCol, costMs, colData.IdVal)
	}
}

// SysExit 关机前同步刷盘,必须在关机 handler 内同步执行.
func SysExit(sig os.Signal) {
	_ = sig
	xrlog.DetailLog.Warning(gctx.New(), "准备关机,开始强制同步内存数据到数据库")
	runShutdownSynLoop()
}

func runShutdownSynLoop() {
	ctx := gctx.New()
	for {
		rows, _ := consumeFlush(ctx, flushReasonShutdown, 0, 0, true)
		if rows == 0 && allSynCachesIdle() {
			xrlog.DetailLog.Info(ctx, "syndb关机刷盘完成")
			return
		}
		time.Sleep(CloseTime)
	}
}

func allSynCachesIdle() bool {
	for _, colCache := range synCacheMap {
		if !colCache.isIdle() {
			return false
		}
	}
	return true
}

func consume(ctx context.Context) {
	sysCpu, cpuIdle := sampleSystemCPU()
	consumeFlush(ctx, flushReasonIdle, sysCpu, cpuIdle, false)
}

func consumeFlush(ctx context.Context, defaultReason string, sysCpu, cpuIdle float64, shutdown bool) (int, string) {
	now := time.Now()
	pullAllCaches()

	candidates := collectFlushCandidates(now, cpuIdle, shutdown)
	if len(candidates) == 0 {
		return 0, ""
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].cache.lastFlushTime.Before(candidates[j].cache.lastFlushTime)
	})

	start := time.Now()
	details := make([]queueFlushDetail, 0, len(candidates))
	totalRows := 0
	idleCount := 0
	forceCount := 0
	remainingBudget := synCfg.batchSize
	for _, cand := range candidates {
		if remainingBudget <= 0 {
			break
		}
		rows := safeBatchSave(cand.cache, remainingBudget)
		if rows <= 0 {
			continue
		}
		cand.cache.lastFlushTime = now
		totalRows += rows
		remainingBudget -= rows
		details = append(details, queueFlushDetail{
			table:  cand.cache.TbName,
			col:    cand.cache.ColName,
			rows:   rows,
			reason: cand.reason,
			waitMs: cand.waitMs,
		})
		switch cand.reason {
		case flushReasonForce:
			forceCount++
		case flushReasonIdle:
			idleCount++
		}
	}

	if totalRows <= 0 {
		return 0, ""
	}

	primaryReason := defaultReason
	if !shutdown {
		switch {
		case forceCount > 0 && idleCount > 0:
			primaryReason = flushReasonForce + "+" + flushReasonIdle
		case forceCount > 0:
			primaryReason = flushReasonForce
		default:
			primaryReason = flushReasonIdle
		}
	}

	logLine := formatFlushLog(primaryReason, sysCpu, cpuIdle, start, totalRows, idleCount, forceCount, len(candidates), details)
	xrlog.DetailLog.Info(gctx.New(), logLine)
	return totalRows, logLine
}

func collectFlushCandidates(now time.Time, cpuIdle float64, shutdown bool) []*flushCandidate {
	candidates := make([]*flushCandidate, 0)
	for _, colCache := range synCacheMap {
		if len(colCache.Pending) == 0 {
			continue
		}
		wait := pendingWaitDuration(colCache, now)
		waitMs := wait.Milliseconds()
		switch {
		case shutdown:
			candidates = append(candidates, &flushCandidate{
				cache:  colCache,
				reason: flushReasonShutdown,
				waitMs: waitMs,
			})
		case wait >= synCfg.maxPendingWait:
			candidates = append(candidates, &flushCandidate{
				cache:  colCache,
				reason: flushReasonForce,
				waitMs: waitMs,
			})
		case cpuIdle >= synCfg.cpuIdlePercent:
			candidates = append(candidates, &flushCandidate{
				cache:  colCache,
				reason: flushReasonIdle,
				waitMs: waitMs,
			})
		}
	}
	return candidates
}

func pendingWaitDuration(colCache *ColSynCache, now time.Time) time.Duration {
	if colCache.firstPendingTime.IsZero() {
		return 0
	}
	if now.Before(colCache.firstPendingTime) {
		return 0
	}
	return now.Sub(colCache.firstPendingTime)
}

func formatFlushLog(primaryReason string, sysCpu, cpuIdle float64, start time.Time, totalRows, idleCount, forceCount, queueCount int, details []queueFlushDetail) string {
	var b strings.Builder
	b.WriteString("syndb刷盘")
	b.WriteString(",reason=")
	b.WriteString(primaryReason)
	if primaryReason != flushReasonShutdown {
		b.WriteString(fmt.Sprintf(",sysCpu=%.1f%%,cpuIdle=%.1f%%,idleThreshold=%.0f%%", sysCpu, cpuIdle, synCfg.cpuIdlePercent))
	}
	b.WriteString(fmt.Sprintf(",batchLimit=%d,queues=%d,rows=%d,idleQueues=%d,forceQueues=%d,costMs=%d",
		synCfg.batchSize, queueCount, totalRows, idleCount, forceCount, time.Since(start).Milliseconds()))
	b.WriteString(",detail=[")
	for i, item := range details {
		if i > 0 {
			b.WriteString("; ")
		}
		b.WriteString(fmt.Sprintf("%s:%s:%d:%s:%dms", item.table, item.col, item.rows, item.reason, item.waitMs))
	}
	b.WriteString("]")
	return b.String()
}

func pullAllCaches() {
	for _, colCache := range synCacheMap {
		if colCache.isIdle() {
			continue
		}
		safePullData(colCache)
	}
}

func safePullData(colCache *ColSynCache) {
	gutil.TryCatch(gctx.New(), func(ctx context.Context) {
		colCache.PullData()
	}, func(ctx context.Context, exception error) {
		xrlog.ErrorWithErr(ctx, "SynDb", "PullData,table="+colCache.TbName+",col="+colCache.ColName, exception)
	})
}

func safeBatchSave(colCache *ColSynCache, maxRows int) int {
	if maxRows <= 0 {
		return 0
	}
	var flushedRows int
	gutil.TryCatch(gctx.New(), func(ctx context.Context) {
		var flushed bool
		flushedRows, flushed = colCache.batchSave(maxRows)
		if !flushed {
			flushedRows = 0
		}
	}, func(ctx context.Context, exception error) {
		xrlog.ErrorWithErr(ctx, "SynDb", "batchSave,table="+colCache.TbName+",col="+colCache.ColName, exception)
	})
	return flushedRows
}

func (colCache *ColSynCache) isIdle() bool {
	return colCache.DataQueue.Len() == 0 && len(colCache.Pending) == 0
}

// PullData 拉取队列数据到 Pending(按主键去重,只保留最新值).
func (colCache *ColSynCache) PullData() {
	for {
		select {
		case data := <-colCache.DataQueue.C:
			if len(colCache.Pending) == 0 {
				colCache.firstPendingTime = time.Now()
			}
			colCache.Pending[data.IdVal] = data
		default:
			return
		}
	}
}

func (colCache *ColSynCache) batchSave(maxRows int) (int, bool) {
	if len(colCache.Pending) == common.Zero || maxRows <= 0 {
		return 0, false
	}

	dataMap := make([]map[string]interface{}, 0, maxRows)
	savedKeys := make([]any, 0, maxRows)
	for idVal, val := range colCache.Pending {
		dataMap = append(dataMap, g.Map{
			colCache.IdName:  idVal,
			colCache.ColName: val.ColVal,
		})
		savedKeys = append(savedKeys, idVal)
		if len(dataMap) >= maxRows {
			break
		}
	}

	rows := len(dataMap)
	_, err := g.DB().Model(colCache.TbName).Data(dataMap).Batch(rows).Save()
	saved := rows
	if err != nil {
		saved = colCache.saveRowsDiscardOnError(dataMap)
	}

	for _, idVal := range savedKeys {
		delete(colCache.Pending, idVal)
	}
	if len(colCache.Pending) == 0 {
		colCache.firstPendingTime = time.Time{}
	}

	if saved == 0 && err != nil {
		xrlog.DetailLog.Warningf(gctx.New(), "syndb batchSave失败已全部丢弃,table=%v,col=%v,rows=%v,err=%v",
			colCache.TbName, colCache.ColName, rows, err)
		return 0, true
	}
	if saved < rows && err != nil {
		xrlog.DetailLog.Warningf(gctx.New(), "syndb batchSave部分落库,table=%v,col=%v,saved=%v,discarded=%v,err=%v",
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
				xrlog.DetailLog.Warningf(gctx.New(), "syndb落库冲突已丢弃,table=%v,col=%v,id=%v,err=%v",
					colCache.TbName, colCache.ColName, idVal, err)
			} else {
				xrlog.DetailLog.Warningf(gctx.New(), "syndb落库失败已丢弃,table=%v,col=%v,id=%v,err=%v",
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
