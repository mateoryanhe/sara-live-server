package livecfg

import (
	"sync/atomic"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/live"
)

const DefaultAudienceListRefreshSeconds uint32 = 300

const DefaultOneToOneDailyFreeSeconds uint32 = 30

type liveCfgSnapshot struct {
	PaidDanmakuPrice           float64
	VideoCallTicketEnabled     bool
	AudienceListRefreshSeconds uint32
	OneToOneDailyFreeSeconds   uint32
}

var (
	liveCfgCache         atomic.Value // *liveCfgSnapshot
	emptyLiveCfgSnapshot = &liveCfgSnapshot{
		VideoCallTicketEnabled:     true,
		AudienceListRefreshSeconds: DefaultAudienceListRefreshSeconds,
		OneToOneDailyFreeSeconds:   DefaultOneToOneDailyFreeSeconds,
	}
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

func NormalizeAudienceListRefreshSeconds(v uint32) uint32 {
	if v == 0 {
		return DefaultAudienceListRefreshSeconds
	}
	return v
}

func toLiveCfgSnapshot(row *entity.LiveCfg) *liveCfgSnapshot {
	if row == nil {
		return emptyLiveCfgSnapshot
	}
	return &liveCfgSnapshot{
		PaidDanmakuPrice:           row.PaidDanmakuPrice,
		VideoCallTicketEnabled:     row.VideoCallTicketEnabled,
		AudienceListRefreshSeconds: NormalizeAudienceListRefreshSeconds(row.AudienceListRefreshSeconds),
		OneToOneDailyFreeSeconds:   NormalizeOneToOneDailyFreeSeconds(row.OneToOneDailyFreeSeconds),
	}
}

func NormalizeOneToOneDailyFreeSeconds(v uint32) uint32 {
	return v
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

// GetAudienceListRefreshSeconds 在线观众列表刷新间隔(秒)，未配置时默认 300。
func GetAudienceListRefreshSeconds() uint32 {
	return getLiveCfgCache().AudienceListRefreshSeconds
}

// GetAudienceListRefreshDuration 在线观众列表刷新间隔。
func GetAudienceListRefreshDuration() time.Duration {
	return time.Duration(GetAudienceListRefreshSeconds()) * time.Second
}

// GetOneToOneDailyFreeSeconds 1v1 观众对单主播每日免费通话秒数，0 表示关闭。
func GetOneToOneDailyFreeSeconds() uint32 {
	return getLiveCfgCache().OneToOneDailyFreeSeconds
}
