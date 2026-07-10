package syndb

import (
	"time"
	"xr-game-server/constants/common"
	"xr-game-server/constants/db"
	"xr-game-server/core/cfg"
)

const (
	// LazySynPeriod 延迟同步默认 15 秒落库
	LazySynPeriod = 15 * 1000
)

var lazyMap = make(map[string]*ColSynCache)

// RegLazy 注册延迟缓冲列(低频落库,周期见 bufferSize.db.lazy.period)
func RegLazy(tbName db.TbName, tbCol db.TbCol) {
	colKey := string(tbName) + ":" + string(tbCol)
	period := LazySynPeriod
	if cfg.LazyDbBufferCfg.Period > common.Zero {
		period = cfg.LazyDbBufferCfg.Period
	}
	lazyMap[colKey] = newColSynCache(string(tbName), string(tbCol), time.Duration(period)*time.Millisecond)
}
