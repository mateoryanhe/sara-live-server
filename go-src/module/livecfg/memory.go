package livecfg

import (
	"sync/atomic"
	"time"
	"xr-game-server/dao/livecfgdao"
	"xr-game-server/entity"
)

const defaultPrivateRoomFreeWatchSeconds uint32 = 420

type liveCfgSnapshot struct {
	PaidDanmakuPrice            float64
	PrivateRoomFreeWatchSeconds uint32
}

var (
	liveCfgCache         atomic.Value // *liveCfgSnapshot
	emptyLiveCfgSnapshot = &liveCfgSnapshot{
		PrivateRoomFreeWatchSeconds: defaultPrivateRoomFreeWatchSeconds,
	}
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
		PaidDanmakuPrice:            row.PaidDanmakuPrice,
		PrivateRoomFreeWatchSeconds: row.PrivateRoomFreeWatchSeconds,
	}
}

// GetPaidDanmakuPrice 获取付费弹幕价格(钻石)
func GetPaidDanmakuPrice() float64 {
	return getLiveCfgCache().PaidDanmakuPrice
}

// GetPrivateRoomFreeWatchSeconds 获取私密直播间免费观看时长(秒)
func GetPrivateRoomFreeWatchSeconds() uint32 {
	return getLiveCfgCache().PrivateRoomFreeWatchSeconds
}

// GetPrivateRoomFreeWatchDuration 获取私密直播间免费观看时长
func GetPrivateRoomFreeWatchDuration() time.Duration {
	return time.Duration(GetPrivateRoomFreeWatchSeconds()) * time.Second
}
