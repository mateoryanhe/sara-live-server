package activity

import (
	"sync/atomic"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/activitydto"
	activityentity "xr-game-server/entity/activity"
)

const (
	defaultInviteRewardPercent = 5.0
	defaultInviteRewardDays    = 30
)

type inviteRewardCfgSnapshot struct {
	ID            uint64
	Enabled       bool
	RewardPercent float64
	ValidDays     int
	CreatedAt     string
	UpdatedAt     string
}

var inviteRewardCfgCache atomic.Value

func reloadInviteRewardCfgMemory() {
	row := cfgdao.LoadInviteRechargeRewardCfg()
	inviteRewardCfgCache.Store(toInviteRewardCfgSnapshot(row))
}

// ReloadInviteRechargeRewardCache 从 DB 重新加载邀请充值返还配置缓存
func ReloadInviteRechargeRewardCache() {
	reloadInviteRewardCfgMemory()
}

func getInviteRewardCfgCache() *inviteRewardCfgSnapshot {
	if inviteRewardCfgCache.Load() == nil {
		reloadInviteRewardCfgMemory()
	}
	v := inviteRewardCfgCache.Load()
	if v == nil {
		return &inviteRewardCfgSnapshot{
			RewardPercent: defaultInviteRewardPercent,
			ValidDays:     defaultInviteRewardDays,
		}
	}
	snap, ok := v.(*inviteRewardCfgSnapshot)
	if !ok || snap == nil {
		return &inviteRewardCfgSnapshot{
			RewardPercent: defaultInviteRewardPercent,
			ValidDays:     defaultInviteRewardDays,
		}
	}
	return snap
}

func toInviteRewardCfgSnapshot(row *activityentity.InviteRechargeRewardCfg) *inviteRewardCfgSnapshot {
	if row == nil {
		return &inviteRewardCfgSnapshot{
			RewardPercent: defaultInviteRewardPercent,
			ValidDays:     defaultInviteRewardDays,
		}
	}
	return &inviteRewardCfgSnapshot{
		ID:            row.ID,
		Enabled:       row.Enabled,
		RewardPercent: normalizeInviteRewardPercent(row.RewardPercent),
		ValidDays:     normalizeInviteRewardDays(row.ValidDays),
		CreatedAt:     formatTime(row.CreatedAt),
		UpdatedAt:     formatTime(row.UpdatedAt),
	}
}

func toInviteRewardCfgItem(snap *inviteRewardCfgSnapshot) *activitydto.InviteRechargeRewardCfgItem {
	if snap == nil {
		return defaultInviteRewardCfgItem()
	}
	return &activitydto.InviteRechargeRewardCfgItem{
		ID:            formatUintID(snap.ID),
		Enabled:       snap.Enabled,
		RewardPercent: snap.RewardPercent,
		ValidDays:     snap.ValidDays,
		CreatedAt:     snap.CreatedAt,
		UpdatedAt:     snap.UpdatedAt,
	}
}

func defaultInviteRewardCfgItem() *activitydto.InviteRechargeRewardCfgItem {
	return &activitydto.InviteRechargeRewardCfgItem{
		RewardPercent: defaultInviteRewardPercent,
		ValidDays:     defaultInviteRewardDays,
	}
}

func normalizeInviteRewardPercent(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}

func normalizeInviteRewardDays(v int) int {
	if v < 0 {
		return 0
	}
	return v
}
