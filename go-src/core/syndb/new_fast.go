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
	colKey := string(tbName) + ":" + string(tbCol)
	chanSize := ChanMax
	if cfg.FastDbBufferCfg.Small > common.Zero {
		chanSize = cfg.FastDbBufferCfg.Small
	}
	period := QuickSynPeriod
	if cfg.FastDbBufferCfg.Period > common.Zero {
		period = cfg.FastDbBufferCfg.Period
	}
	quickMap[colKey] = &ColSynCache{
		ColName:   string(tbCol),
		TbName:    string(tbName),
		DataQueue: make(chan *ColData, chanSize),
		TempData:  make([]*ColData, 0),
		Period:    time.Duration(period) * time.Millisecond,
		LastTime:  time.Now(),
		IdName:    string(db.IdName),
	}
}

func RegQuickWithMiddle(tbName db.TbName, tbCol db.TbCol) {
	colKey := string(tbName) + ":" + string(tbCol)
	chanSize := ChanMax
	if cfg.FastDbBufferCfg.Middle > common.Zero {
		chanSize = cfg.FastDbBufferCfg.Middle
	}
	period := QuickSynPeriod
	if cfg.FastDbBufferCfg.Period > common.Zero {
		period = cfg.FastDbBufferCfg.Period
	}
	quickMap[colKey] = &ColSynCache{
		ColName:   string(tbCol),
		TbName:    string(tbName),
		DataQueue: make(chan *ColData, chanSize),
		TempData:  make([]*ColData, 0),
		Period:    time.Duration(period) * time.Millisecond,
		LastTime:  time.Now(),
		IdName:    string(db.IdName),
	}
}
func RegQuickWithLarge(tbName db.TbName, tbCol db.TbCol) {
	colKey := string(tbName) + ":" + string(tbCol)
	chanSize := ChanMax
	if cfg.FastDbBufferCfg.Large > common.Zero {
		chanSize = cfg.FastDbBufferCfg.Large
	}
	period := QuickSynPeriod
	if cfg.FastDbBufferCfg.Period > common.Zero {
		period = cfg.FastDbBufferCfg.Period
	}
	quickMap[colKey] = &ColSynCache{
		ColName:   string(tbCol),
		TbName:    string(tbName),
		DataQueue: make(chan *ColData, chanSize),
		TempData:  make([]*ColData, 0),
		Period:    time.Duration(period) * time.Millisecond,
		LastTime:  time.Now(),
		IdName:    string(db.IdName),
	}
}
