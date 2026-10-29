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
	PaidDanmakuPrice                float64
	VideoCallTicketEnabled          bool
	AudienceListRefreshSeconds      uint32
	OneToOneDailyFreeSeconds        uint32
	LiveRoomCallMinTotalRechargeUsd float64
	ShowcaseLiveRoomEnabled         bool
	GameLiveRoomEnabled             bool
	VoiceChatLiveRoomEnabled        bool
	OneToOneLiveRoomEnabled         bool
}

var (
	liveCfgCache         atomic.Value // *liveCfgSnapshot
	emptyLiveCfgSnapshot = &liveCfgSnapshot{
		VideoCallTicketEnabled:          true,
		AudienceListRefreshSeconds:      DefaultAudienceListRefreshSeconds,
		OneToOneDailyFreeSeconds:        DefaultOneToOneDailyFreeSeconds,
		LiveRoomCallMinTotalRechargeUsd: entity.LiveRoomCallMinTotalRechargeUSD,
	}
)

func reloadLiveCfgMemory() {
	// 启动与 CMS saveLiveCfg 后：从 live_cfgs 表加载并写入进程内快照
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
		PaidDanmakuPrice:                row.PaidDanmakuPrice,
		VideoCallTicketEnabled:            row.VideoCallTicketEnabled,
		AudienceListRefreshSeconds:        NormalizeAudienceListRefreshSeconds(row.AudienceListRefreshSeconds),
		OneToOneDailyFreeSeconds:          NormalizeOneToOneDailyFreeSeconds(row.OneToOneDailyFreeSeconds),
		LiveRoomCallMinTotalRechargeUsd: NormalizeLiveRoomCallMinTotalRechargeUsd(row.LiveRoomCallMinTotalRechargeUsd),
		ShowcaseLiveRoomEnabled:         row.ShowcaseLiveRoomEnabled,
		GameLiveRoomEnabled:             row.GameLiveRoomEnabled,
		VoiceChatLiveRoomEnabled:        row.VoiceChatLiveRoomEnabled,
		OneToOneLiveRoomEnabled:         row.OneToOneLiveRoomEnabled,
	}
}

func NormalizeLiveRoomCallMinTotalRechargeUsd(v float64) float64 {
	if v < 0 {
		return entity.LiveRoomCallMinTotalRechargeUSD
	}
	return v
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

// AppLiveRoomCallTicketPrice 下发给 App 的直播间视频通话门票价；全局关闭门票时返回 0。
func AppLiveRoomCallTicketPrice(roomTicket float64) float64 {
	if !IsVideoCallTicketEnabled() {
		return 0
	}
	if roomTicket < 0 {
		return 0
	}
	return roomTicket
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

// GetLiveRoomCallMinTotalRechargeUSD 直播间 allowCallIcon 最低累计充值(USD)
func GetLiveRoomCallMinTotalRechargeUSD() float64 {
	return getLiveCfgCache().LiveRoomCallMinTotalRechargeUsd
}

// IsShowcaseLiveRoomEnabled App 是否展示秀场直播间入口；未配置时默认关闭。
func IsShowcaseLiveRoomEnabled() bool {
	return getLiveCfgCache().ShowcaseLiveRoomEnabled
}

// IsGameLiveRoomEnabled App 是否展示游戏直播间入口；未配置时默认关闭。
func IsGameLiveRoomEnabled() bool {
	return getLiveCfgCache().GameLiveRoomEnabled
}

// IsVoiceChatLiveRoomEnabled App 是否展示语聊房入口；未配置时默认关闭。
func IsVoiceChatLiveRoomEnabled() bool {
	return getLiveCfgCache().VoiceChatLiveRoomEnabled
}

// IsOneToOneLiveRoomEnabled App 是否展示 1v1 房间页面/入口；未配置时默认关闭。
func IsOneToOneLiveRoomEnabled() bool {
	return getLiveCfgCache().OneToOneLiveRoomEnabled
}

