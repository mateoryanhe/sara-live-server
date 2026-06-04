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
)

const (
	// Max 批量同步数量
	Max     = 150
	ChanMax = 1000
	// CloseTime 服务器关闭,同步周期变更
	CloseTime = 50 * time.Millisecond
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
	lazyMap[string(tbName)+":"+string(tbCol)].DataQueue <- colData
}

// AddDataToQuickChan 加入变更数据到快速缓存区
func AddDataToQuickChan(tbName db.TbName, tbCol db.TbCol, colData *ColData) {
	quickMap[string(tbName)+":"+string(tbCol)].DataQueue <- colData
}

// SysExit 调整全部数据库同步组件同步时间到最小
func SysExit(sig os.Signal) {
	g.Log().Warning(gctx.New(), "准备关机,开始强制同步内存数据到数据库")
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
	dataMap := make([]map[string]interface{}, common.Zero)
	for _, val := range colCache.TempData {
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
