package syndb

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gutil"
	"os"
	"time"
	"xr-game-server/constants/common"
	"xr-game-server/constants/db"
	"xr-game-server/core/lambda"
)

const (
	// Max 批量同步数量
	Max     = 150
	ChanMax = 1000
	// CloseTime 服务器关闭,同步周期变更
	CloseTime = 50 * time.Millisecond
	// chanWriteSlowMs 写入缓冲通道超过该耗时(毫秒)时打印日志
	chanWriteSlowMs = 3
	// QuickSynPeriod 快速同步
)

type ColData struct {
	IdVal  any
	ColVal any
}

// ColSynCache 内存同步数据库工具
type ColSynCache struct {
	//缓冲数据
	DataQueue chan *ColData
	//准备同步数据
	TempData []*ColData
	//同步频率
	Period time.Duration
	//上次同步时间
	LastTime time.Time
	//表名
	TbName string
	//列名
	ColName string
	//主键列名
	IdName string
}

// AddDataToLazyChan 加入变更数据到延迟缓存区
func AddDataToLazyChan(tbName db.TbName, tbCol db.TbCol, colData *ColData) {
	addDataToChan("lazy", lazyMap, tbName, tbCol, colData)
}

// AddDataToQuickChan 加入变更数据到快速缓存区
func AddDataToQuickChan(tbName db.TbName, tbCol db.TbCol, colData *ColData) {
	addDataToChan("quick", quickMap, tbName, tbCol, colData)
}

func addDataToChan(queueType string, cacheMap map[string]*ColSynCache, tbName db.TbName, tbCol db.TbCol, colData *ColData) {
	key := string(tbName) + ":" + string(tbCol)
	start := time.Now()
	cacheMap[key].DataQueue <- colData
	costMs := time.Since(start).Milliseconds()
	if costMs > chanWriteSlowMs {
		g.Log().Infof(gctx.New(), "syndb写入%s缓冲耗时=%vms,tb=%v,col=%v,id=%v", queueType, costMs, tbName, tbCol, colData.IdVal)
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
	for _, val := range lazyMap {
		gutil.TryCatch(gctx.New(), func(ctx context.Context) {
			val.PullData()
			val.Syn()
		}, func(ctx context.Context, exception error) {
		})
	}
	for _, quick := range quickMap {
		gutil.TryCatch(gctx.New(), func(ctx context.Context) {
			quick.PullData()
			quick.Syn()
		}, func(ctx context.Context, exception error) {
		})
	}
}

// PullData 拉取数据到同步队列
func (colCache *ColSynCache) PullData() {
	select {
	//拉取一个数据到同步队列中
	case data := <-colCache.DataQueue:
		{
			colCache.TempData = append(colCache.TempData, data)
			//如果是第一个数据
			if len(colCache.TempData) == 1 {
				colCache.LastTime = time.Now()
			}
		}
		//防止堵塞
	case <-time.After(time.Nanosecond):
		{
			return
		}
	}
}

func (colCache *ColSynCache) Syn() {
	//时间到,强制同步
	now := time.Now()
	targetTime := colCache.LastTime.Add(colCache.Period)

	if len(colCache.TempData) > 0 && targetTime.Before(now) {
		colCache.batchSave()
	}
	//到达最大数量,强制同步
	if len(colCache.TempData) >= Max {
		colCache.batchSave()
	}
}

func (colCache *ColSynCache) batchSave() {
	if len(colCache.TempData) == common.Zero {
		return
	}
	//过滤一下重复的
	list := make([]*ColData, 0)
	for _, data := range colCache.TempData {
		if lambda.AnyMatch(list, func(val *ColData) bool {
			return val.IdVal == data.IdVal
		}) {
			list = lambda.Filter(list, func(val *ColData) bool {
				return val.IdVal != data.IdVal
			})
			list = append(list, data)
		} else {
			list = append(list, data)
		}
	}
	dataMap := make([]map[string]interface{}, common.Zero)
	for _, val := range list {
		//检查是否重复
		dataMap = append(dataMap, g.Map{
			colCache.IdName:  val.IdVal,
			colCache.ColName: val.ColVal,
		})
	}
	_, err := g.DB().Model(colCache.TbName).Data(dataMap).Batch(len(dataMap)).Save()
	if err != nil {
		return
	}
	colCache.TempData = make([]*ColData, common.Zero)
}

// ChangeSynTime 变更同步时间
func (colCache *ColSynCache) ChangeSynTime() {
	colCache.Period = CloseTime
}
