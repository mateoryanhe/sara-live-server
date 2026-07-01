package syndb

import (
	"time"
	"xr-game-server/constants/common"
	"xr-game-server/constants/db"
	"xr-game-server/core/cfg"
)

const (
	QuickSynPeriod = 10
)

var quickMap = make(map[string]*ColSynCache)

func RegQuickWithSmall(tbName db.TbName, tbCol db.TbCol) {
	regQuick(tbName, tbCol)
}

func RegQuickWithMiddle(tbName db.TbName, tbCol db.TbCol) {
	regQuick(tbName, tbCol)
}

func RegQuickWithLarge(tbName db.TbName, tbCol db.TbCol) {
	regQuick(tbName, tbCol)
}

func regQuick(tbName db.TbName, tbCol db.TbCol) {
	colKey := string(tbName) + ":" + string(tbCol)
	period := QuickSynPeriod
	if cfg.FastDbBufferCfg.Period > common.Zero {
		period = cfg.FastDbBufferCfg.Period
	}
	quickMap[colKey] = newColSynCache(string(tbName), string(tbCol), time.Duration(period)*time.Millisecond)
}
