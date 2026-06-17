package livecfg

import (
	"sync/atomic"
	"xr-game-server/dao/livecfgdao"
	"xr-game-server/entity"
)

type liveCfgSnapshot struct {
	PaidDanmakuPrice float64
}

var (
	liveCfgCache         atomic.Value // *liveCfgSnapshot
	emptyLiveCfgSnapshot = &liveCfgSnapshot{}
)

func reloadLiveCfgMemory() {
	liveCfgCache.Store(toLiveCfgSnapshot(livecfgdao.Load()))
}

func getLiveCfgCache() *liveCfgSnapshot {
	v := liveCfgCache.Load()
	if v == nil {
		return emptyLiveCfgSnapshot
	}
	cfg, ok := v.(*liveCfgSnapshot)
	if !ok || cfg == nil {
		return emptyLiveCfgSnapshot
	}
	return cfg
}

func toLiveCfgSnapshot(row *entity.LiveCfg) *liveCfgSnapshot {
	if row == nil {
		return emptyLiveCfgSnapshot
	}
	return &liveCfgSnapshot{
		PaidDanmakuPrice: row.PaidDanmakuPrice,
	}
}

// GetPaidDanmakuPrice 获取付费弹幕价格(钻石)
func GetPaidDanmakuPrice() float64 {
	return getLiveCfgCache().PaidDanmakuPrice
}
