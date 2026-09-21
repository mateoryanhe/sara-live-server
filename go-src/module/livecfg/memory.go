package livecfg

import (
	"sync/atomic"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/live"
)

type liveCfgSnapshot struct {
	PaidDanmakuPrice       float64
	VideoCallTicketEnabled bool
}

var (
	liveCfgCache         atomic.Value // *liveCfgSnapshot
	emptyLiveCfgSnapshot = &liveCfgSnapshot{VideoCallTicketEnabled: true}
)

func reloadLiveCfgMemory() {
	liveCfgCache.Store(toLiveCfgSnapshot(cfgdao.LoadLiveCfg()))
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
		PaidDanmakuPrice:       row.PaidDanmakuPrice,
		VideoCallTicketEnabled: row.VideoCallTicketEnabled,
	}
}

// GetPaidDanmakuPrice 获取付费弹幕价格(钻石)
func GetPaidDanmakuPrice() float64 {
	return getLiveCfgCache().PaidDanmakuPrice
}

// IsVideoCallTicketEnabled 返回直播间来源视频通话接通时是否扣门票。
// 未配置直播参数时默认开启，保持历史扣费行为。
func IsVideoCallTicketEnabled() bool {
	return getLiveCfgCache().VideoCallTicketEnabled
}
