package syndb

import (
	"time"
	"xr-game-server/constants/common"
	"xr-game-server/constants/db"
	"xr-game-server/core/cfg"
)

const (
	// QuickSynPeriod 快速缓冲默认 500ms 落库(配合 500ms consume,目标约 1 秒内入库)
	QuickSynPeriod = 500
)

var quickMap = make(map[string]*ColSynCache)

// RegQuick 注册快速缓冲列(约 1 秒内落库,周期见 bufferSize.db.fast.period)
func RegQuick(tbName db.TbName, tbCol db.TbCol) {
	colKey := string(tbName) + ":" + string(tbCol)
	period := QuickSynPeriod
	if cfg.FastDbBufferCfg.Period > common.Zero {
		period = cfg.FastDbBufferCfg.Period
	}
	quickMap[colKey] = newColSynCache(string(tbName), string(tbCol), time.Duration(period)*time.Millisecond)
}
