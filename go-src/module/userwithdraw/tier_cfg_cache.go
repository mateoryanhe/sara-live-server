package userwithdraw

import (
	"sync/atomic"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/userwithdrawtiercfgdto"
	userentity "xr-game-server/entity/user"
)

type userWithdrawTierCfgSnapshot struct {
	byId map[uint64]*userwithdrawtiercfgdto.AppUserWithdrawTierItem
	list []*userwithdrawtiercfgdto.AppUserWithdrawTierItem
}

var (
	userWithdrawTierCfgCache     atomic.Value // *userWithdrawTierCfgSnapshot
	emptyUserWithdrawTierCfgList = make([]*userwithdrawtiercfgdto.AppUserWithdrawTierItem, 0)
)

func toAppUserWithdrawTierItem(r *userentity.UserWithdrawTierCfg) *userwithdrawtiercfgdto.AppUserWithdrawTierItem {
	return &userwithdrawtiercfgdto.AppUserWithdrawTierItem{
		ID:             r.ID,
		Name:           r.Name,
		WithdrawUsd:    r.WithdrawUsd,
		MinGameConsumeGold: r.MinGameConsumeGold,
	}
}

func loadUserWithdrawTierCfgCache() []*userwithdrawtiercfgdto.AppUserWithdrawTierItem {
	rows := cfgdao.GetOnShelfUserWithdrawTierCfgs()
	byId := make(map[uint64]*userwithdrawtiercfgdto.AppUserWithdrawTierItem, len(rows))
	list := make([]*userwithdrawtiercfgdto.AppUserWithdrawTierItem, 0, len(rows))
	for _, r := range rows {
		item := toAppUserWithdrawTierItem(r)
		byId[r.ID] = item
		list = append(list, item)
	}
	userWithdrawTierCfgCache.Store(&userWithdrawTierCfgSnapshot{byId: byId, list: list})
	return list
}

func reloadUserWithdrawTierCfgCache() {
	loadUserWithdrawTierCfgCache()
}

// ReloadUserWithdrawTierCfgCache 供同步等模块调用
func ReloadUserWithdrawTierCfgCache() {
	reloadUserWithdrawTierCfgCache()
}

func getUserWithdrawTierCfgSnapshot() *userWithdrawTierCfgSnapshot {
	v := userWithdrawTierCfgCache.Load()
	if v == nil {
		return &userWithdrawTierCfgSnapshot{
			byId: make(map[uint64]*userwithdrawtiercfgdto.AppUserWithdrawTierItem),
			list: emptyUserWithdrawTierCfgList,
		}
	}
	return v.(*userWithdrawTierCfgSnapshot)
}

func GetOnShelfUserWithdrawTiers() []*userwithdrawtiercfgdto.AppUserWithdrawTierItem {
	if userWithdrawTierCfgCache.Load() == nil {
		return loadUserWithdrawTierCfgCache()
	}
	return getUserWithdrawTierCfgSnapshot().list
}

func GetUserWithdrawTierFromCacheById(id uint64) *userwithdrawtiercfgdto.AppUserWithdrawTierItem {
	if userWithdrawTierCfgCache.Load() == nil {
		loadUserWithdrawTierCfgCache()
	}
	return getUserWithdrawTierCfgSnapshot().byId[id]
}
