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

// RegLazyWithSmall 单行数据同步 角色表
func RegLazyWithSmall(tbName db.TbName, tbCol db.TbCol) {
	colKey := string(tbName) + ":" + string(tbCol)
	chanSize := ChanMax
	if cfg.LazyDbBufferCfg.Small > common.Zero {
		chanSize = cfg.LazyDbBufferCfg.Small
	}
	period := LazySynPeriod
	if cfg.LazyDbBufferCfg.Period > common.Zero {
		period = cfg.LazyDbBufferCfg.Period
	}
	lazyMap[colKey] = &ColSynCache{
		ColName:   string(tbCol),
		TbName:    string(tbName),
		DataQueue: make(chan *ColData, chanSize),
		TempData:  make([]*ColData, 0),
		Period:    time.Duration(period) * time.Millisecond,
		LastTime:  time.Now(),
		IdName:    string(db.IdName),
	}
}

// RegLazyWithMiddle 多行数据同步 背包,邮件
func RegLazyWithMiddle(tbName db.TbName, tbCol db.TbCol) {
	colKey := string(tbName) + ":" + string(tbCol)
	chanSize := ChanMax
	if cfg.LazyDbBufferCfg.Middle > common.Zero {
		chanSize = cfg.LazyDbBufferCfg.Middle
	}
	period := LazySynPeriod
	if cfg.LazyDbBufferCfg.Period > common.Zero {
		period = cfg.LazyDbBufferCfg.Period
	}
	lazyMap[colKey] = &ColSynCache{
		ColName:   string(tbCol),
		TbName:    string(tbName),
		DataQueue: make(chan *ColData, chanSize),
		TempData:  make([]*ColData, 0),
		Period:    time.Duration(period) * time.Millisecond,
		LastTime:  time.Now(),
		IdName:    string(db.IdName),
	}
}

func RegLazyWithLarge(tbName db.TbName, tbCol db.TbCol) {
	colKey := string(tbName) + ":" + string(tbCol)
	chanSize := ChanMax
	if cfg.LazyDbBufferCfg.Large > common.Zero {
		chanSize = cfg.LazyDbBufferCfg.Large
	}
	period := LazySynPeriod
	if cfg.LazyDbBufferCfg.Period > common.Zero {
		period = cfg.LazyDbBufferCfg.Period
	}
	lazyMap[colKey] = &ColSynCache{
		ColName:   string(tbCol),
		TbName:    string(tbName),
		DataQueue: make(chan *ColData, chanSize),
		TempData:  make([]*ColData, 0),
		Period:    time.Duration(period) * time.Millisecond,
		LastTime:  time.Now(),
		IdName:    string(db.IdName),
	}
}
