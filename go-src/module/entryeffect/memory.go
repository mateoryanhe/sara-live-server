package entryeffect

import (
	"sort"
	"sync/atomic"
	"xr-game-server/dao/entryeffectdao"
	"xr-game-server/dto/entryeffectdto"
	"xr-game-server/entity"
	"xr-game-server/module/upload"
)

type entryEffectSnapshot struct {
	byID    map[uint64]*entity.LiveEntryEffect
	onShelf []*entryeffectdto.AppEntryEffectItem
}

var (
	entryEffectCache     atomic.Value // *entryEffectSnapshot
	emptyEntryEffectList = make([]*entryeffectdto.AppEntryEffectItem, 0)
)

func Init() {
	reloadEntryEffectMemory()
}

func reloadEntryEffectMemory() {
	rows := entryeffectdao.GetAll()
	byID := make(map[uint64]*entity.LiveEntryEffect, len(rows))
	onShelf := make([]*entryeffectdto.AppEntryEffectItem, 0, len(rows))
	for _, row := range rows {
		byID[row.ID] = row
		if row.Status == entity.LiveEntryEffectStatusOnShelf {
			onShelf = append(onShelf, toAppEntryEffectItem(row))
		}
	}
	sort.Slice(onShelf, func(i, j int) bool {
		if onShelf[i].LevelStart != onShelf[j].LevelStart {
			return onShelf[i].LevelStart < onShelf[j].LevelStart
		}
		if onShelf[i].LevelEnd != onShelf[j].LevelEnd {
			return onShelf[i].LevelEnd < onShelf[j].LevelEnd
		}
		return onShelf[i].ID > onShelf[j].ID
	})
	entryEffectCache.Store(&entryEffectSnapshot{
		byID:    byID,
		onShelf: onShelf,
	})
}

func getEntryEffectSnapshot() *entryEffectSnapshot {
	v := entryEffectCache.Load()
	if v == nil {
		return &entryEffectSnapshot{
			byID:    make(map[uint64]*entity.LiveEntryEffect),
			onShelf: emptyEntryEffectList,
		}
	}
	return v.(*entryEffectSnapshot)
}

func getAppEntryEffectList() []*entryeffectdto.AppEntryEffectItem {
	return getEntryEffectSnapshot().onShelf
}

func toAppEntryEffectItem(row *entity.LiveEntryEffect) *entryeffectdto.AppEntryEffectItem {
	return &entryeffectdto.AppEntryEffectItem{
		ID:         row.ID,
		Name:       row.Name,
		LevelStart: row.LevelStart,
		LevelEnd:   row.LevelEnd,
		Animation:  upload.GetUrlByName(row.Animation),
	}
}
