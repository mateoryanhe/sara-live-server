package syndb

import (
	"time"
	"xr-game-server/constants/common"
	"xr-game-server/constants/db"
	"xr-game-server/core/cfg"
)

const (
	// LazySynPeriod 延迟同步
	LazySynPeriod = 15 * 1000
)

var lazyMap = make(map[string]*ColSynCache)

func RegLazyWithSmall(tbName db.TbName, tbCol db.TbCol) {
	regLazy(tbName, tbCol)
}

func RegLazyWithMiddle(tbName db.TbName, tbCol db.TbCol) {
	regLazy(tbName, tbCol)
}

func RegLazyWithLarge(tbName db.TbName, tbCol db.TbCol) {
	regLazy(tbName, tbCol)
}

func regLazy(tbName db.TbName, tbCol db.TbCol) {
	colKey := string(tbName) + ":" + string(tbCol)
	period := LazySynPeriod
	if cfg.LazyDbBufferCfg.Period > common.Zero {
		period = cfg.LazyDbBufferCfg.Period
	}
	lazyMap[colKey] = newColSynCache(string(tbName), string(tbCol), time.Duration(period)*time.Millisecond)
}
