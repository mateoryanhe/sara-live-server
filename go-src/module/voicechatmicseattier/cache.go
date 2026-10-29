package voicechatmicseattier

import (
	"sync/atomic"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/voicechatmicseattierdto"
	liveentity "xr-game-server/entity/live"
)

type tierSnapshot struct {
	list []*voicechatmicseattierdto.AppTierItem
}

var (
	tierCache     atomic.Value // *tierSnapshot
	emptyTierList = make([]*voicechatmicseattierdto.AppTierItem, 0)
)

func toAppTierItem(r *liveentity.VoiceChatMicSeatTier) *voicechatmicseattierdto.AppTierItem {
	return &voicechatmicseattierdto.AppTierItem{
		ID:        r.ID,
		SeatCount: r.SeatCount,
		Sort:      r.Sort,
	}
}

func loadTierCache() []*voicechatmicseattierdto.AppTierItem {
	rows := cfgdao.GetOnShelfVoiceChatMicSeatTiers()
	list := make([]*voicechatmicseattierdto.AppTierItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, toAppTierItem(r))
	}
	tierCache.Store(&tierSnapshot{list: list})
	return list
}

func reloadTierCache() {
	loadTierCache()
}

// ReloadTierCache 供同步等模块调用
func ReloadTierCache() {
	reloadTierCache()
}

func getTierSnapshot() *tierSnapshot {
	v := tierCache.Load()
	if v == nil {
		return &tierSnapshot{list: emptyTierList}
	}
	return v.(*tierSnapshot)
}

func GetOnShelfTiers() []*voicechatmicseattierdto.AppTierItem {
	if tierCache.Load() == nil {
		return loadTierCache()
	}
	return getTierSnapshot().list
}

func Init() {
	loadTierCache()
}
