package livecfg

import (
	"sync/atomic"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/live"
)

const DefaultAudienceListRefreshSeconds uint32 = 300

const DefaultOneToOneDailyFreeSeconds uint32 = 30

const DefaultMaxAnchorPhotos uint32 = 100

const DefaultMaxAnchorVideos uint32 = 100

const DefaultMaxAnchorVideoFileSizeMB uint32 = 100

type liveCfgSnapshot struct {
	PaidDanmakuPrice                float64
	VideoCallTicketEnabled          bool
	AudienceListRefreshSeconds      uint32
	OneToOneDailyFreeSeconds        uint32
	LiveRoomCallMinTotalRechargeUsd float64
	ShowcaseLiveRoomEnabled         bool
	AudienceGameListHidden          bool
	AnchorGameListHidden            bool
	GameLiveRoomEnabled             bool
	VoiceChatLiveRoomEnabled        bool
	OneToOneLiveRoomEnabled         bool
	MaxAnchorPhotos                 uint32
	MaxAnchorVideos                 uint32
	MaxAnchorVideoFileSizeMB        uint32
}

var (
	liveCfgCache         atomic.Value // *liveCfgSnapshot
	emptyLiveCfgSnapshot = &liveCfgSnapshot{
		VideoCallTicketEnabled:          true,
		AudienceListRefreshSeconds:      DefaultAudienceListRefreshSeconds,
		OneToOneDailyFreeSeconds:        DefaultOneToOneDailyFreeSeconds,
		LiveRoomCallMinTotalRechargeUsd: entity.LiveRoomCallMinTotalRechargeUSD,
		MaxAnchorPhotos:                 DefaultMaxAnchorPhotos,
		MaxAnchorVideos:                 DefaultMaxAnchorVideos,
		MaxAnchorVideoFileSizeMB:        DefaultMaxAnchorVideoFileSizeMB,
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
		VideoCallTicketEnabled:          row.VideoCallTicketEnabled,
		AudienceListRefreshSeconds:      NormalizeAudienceListRefreshSeconds(row.AudienceListRefreshSeconds),
		OneToOneDailyFreeSeconds:        NormalizeOneToOneDailyFreeSeconds(row.OneToOneDailyFreeSeconds),
		LiveRoomCallMinTotalRechargeUsd: NormalizeLiveRoomCallMinTotalRechargeUsd(row.LiveRoomCallMinTotalRechargeUsd),
		ShowcaseLiveRoomEnabled:         row.ShowcaseLiveRoomEnabled,
		AudienceGameListHidden:          row.AudienceGameListHidden,
		AnchorGameListHidden:            row.AnchorGameListHidden,
		GameLiveRoomEnabled:             row.GameLiveRoomEnabled,
		VoiceChatLiveRoomEnabled:        row.VoiceChatLiveRoomEnabled,
		OneToOneLiveRoomEnabled:         row.OneToOneLiveRoomEnabled,
		MaxAnchorPhotos:                 NormalizeMaxAnchorPhotos(row.MaxAnchorPhotos),
		MaxAnchorVideos:                 NormalizeMaxAnchorVideos(row.MaxAnchorVideos),
		MaxAnchorVideoFileSizeMB:        NormalizeMaxAnchorVideoFileSizeMB(row.MaxAnchorVideoFileSizeMB),
	}
}

func NormalizeMaxAnchorPhotos(v uint32) uint32 {
	if v == 0 {
		return DefaultMaxAnchorPhotos
	}
	return v
}

func NormalizeMaxAnchorVideos(v uint32) uint32 {
	if v == 0 {
		return DefaultMaxAnchorVideos
	}
	return v
}

func NormalizeMaxAnchorVideoFileSizeMB(v uint32) uint32 {
	if v == 0 {
		return DefaultMaxAnchorVideoFileSizeMB
	}
	return v
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

// GetLiveRoomCallMinTotalRechargeUSD CMS 保留字段：历史为 allowCallIcon 充值门槛，当前未用于通话按钮展示。
func GetLiveRoomCallMinTotalRechargeUSD() float64 {
	return getLiveCfgCache().LiveRoomCallMinTotalRechargeUsd
}

// IsShowcaseLiveRoomEnabled App 是否展示秀场直播间入口；未配置时默认关闭。
func IsShowcaseLiveRoomEnabled() bool {
	return getLiveCfgCache().ShowcaseLiveRoomEnabled
}

// IsAudienceGameListHidden 开启后非主播用户查询 App 游戏列表与直播间推荐游戏为空。
func IsAudienceGameListHidden() bool {
	return getLiveCfgCache().AudienceGameListHidden
}

// IsAnchorGameListHidden 开启后主播用户查询 App 游戏列表与直播间推荐游戏为空。
func IsAnchorGameListHidden() bool {
	return getLiveCfgCache().AnchorGameListHidden
}

// IsGameLiveRoomEnabled App 是否展示游戏直播间入口；未配置时默认关闭。
func IsGameLiveRoomEnabled() bool {
	return getLiveCfgCache().GameLiveRoomEnabled
}

// IsVoiceChatLiveRoomEnabled App 是否展示语聊房入口；未配置时默认关闭。
func IsVoiceChatLiveRoomEnabled() bool {
	return getLiveCfgCache().VoiceChatLiveRoomEnabled
}

// IsOneToOneLiveRoomEnabled App 1v1 总开关：1v1 房间页/入口及直播间 allowCallIcon；未配置时默认关闭。
func IsOneToOneLiveRoomEnabled() bool {
	return getLiveCfgCache().OneToOneLiveRoomEnabled
}

// GetMaxAnchorPhotos 主播相册照片上限；未配置或为 0 时默认 100。
func GetMaxAnchorPhotos() uint32 {
	return getLiveCfgCache().MaxAnchorPhotos
}

// GetMaxAnchorVideos 主播相册视频上限；未配置或为 0 时默认 100。
func GetMaxAnchorVideos() uint32 {
	return getLiveCfgCache().MaxAnchorVideos
}

// GetMaxAnchorVideoFileSizeMB 主播相册单个视频大小上限(MB)；未配置或为 0 时默认 100。
func GetMaxAnchorVideoFileSizeMB() uint32 {
	return getLiveCfgCache().MaxAnchorVideoFileSizeMB
}

// GetMaxAnchorVideoFileSizeBytes 主播相册单个视频大小上限(字节)。
func GetMaxAnchorVideoFileSizeBytes() int64 {
	return int64(GetMaxAnchorVideoFileSizeMB()) * 1024 * 1024
}
