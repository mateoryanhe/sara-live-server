package liveroom

import (
	"sort"
	"strconv"
	"sync/atomic"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity"
)

type roomTagSnapshot struct {
	byID map[uint64]*entity.LiveRoomTag
	list []*liveroomdto.AppLiveRoomTagItem
}

var (
	roomTagCache     atomic.Value // *roomTagSnapshot
	emptyRoomTagList = make([]*liveroomdto.AppLiveRoomTagItem, 0)
)

func initRoomTagMemory() {
	reloadRoomTagMemory()
}

func reloadRoomTagMemory() {
	rows := liveroomdao.GetAllRoomTags()
	byID := make(map[uint64]*entity.LiveRoomTag, len(rows))
	list := make([]*liveroomdto.AppLiveRoomTagItem, 0, len(rows))
	for _, r := range rows {
		byID[r.ID] = r
		list = append(list, toAppLiveRoomTagItem(r))
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].Sort != list[j].Sort {
			return list[i].Sort > list[j].Sort
		}
		return list[i].ID > list[j].ID
	})
	roomTagCache.Store(&roomTagSnapshot{
		byID: byID,
		list: list,
	})
}

func getRoomTagSnapshot() *roomTagSnapshot {
	v := roomTagCache.Load()
	if v == nil {
		return &roomTagSnapshot{
			byID: make(map[uint64]*entity.LiveRoomTag),
			list: emptyRoomTagList,
		}
	}
	return v.(*roomTagSnapshot)
}

func getAppRoomTagList() []*liveroomdto.AppLiveRoomTagItem {
	return getRoomTagSnapshot().list
}

func GetRoomTagFromMemoryById(id uint64) *entity.LiveRoomTag {
	return getRoomTagSnapshot().byID[id]
}

func getRoomTagName(tagId uint64) string {
	if tag := GetRoomTagFromMemoryById(tagId); tag != nil {
		return tag.Name
	}
	return ""
}

func toAppLiveRoomTagItem(tag *entity.LiveRoomTag) *liveroomdto.AppLiveRoomTagItem {
	return &liveroomdto.AppLiveRoomTagItem{
		ID:   strconv.FormatUint(tag.ID, 10),
		Name: tag.Name,
		Sort: tag.Sort,
	}
}
