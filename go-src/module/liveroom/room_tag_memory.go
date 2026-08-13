package liveroom

import (
	"sort"
	"strconv"
	"strings"
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

// getAppRoomTagNormalList App 端普通标签列表(仅 isSpecial=false).
func getAppRoomTagNormalList() []*liveroomdto.AppLiveRoomTagItem {
	all := getAppRoomTagList()
	if len(all) == 0 {
		return emptyRoomTagList
	}
	list := make([]*liveroomdto.AppLiveRoomTagItem, 0, len(all))
	for _, item := range all {
		if item == nil || item.IsSpecial {
			continue
		}
		list = append(list, item)
	}
	return list
}

const (
	specialRoomTagNameAll = "all"
	specialRoomTagNameMy  = "my"
)

func resolveSpecialRoomTagFilterMode(tagId uint64) string {
	if tagId == 0 {
		return ""
	}
	tag := GetRoomTagFromMemoryById(tagId)
	if tag == nil || !tag.IsSpecial {
		return ""
	}
	switch strings.ToLower(strings.TrimSpace(tag.Name)) {
	case specialRoomTagNameAll:
		return specialRoomTagNameAll
	case specialRoomTagNameMy:
		return specialRoomTagNameMy
	default:
		return ""
	}
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

func getRoomTagIdByName(name string) uint64 {
	name = strings.TrimSpace(name)
	if name == "" {
		return 0
	}
	for _, item := range getAppRoomTagList() {
		if strings.EqualFold(item.Name, name) {
			id, err := strconv.ParseUint(item.ID, 10, 64)
			if err == nil {
				return id
			}
		}
	}
	return 0
}

func toAppLiveRoomTagItem(tag *entity.LiveRoomTag) *liveroomdto.AppLiveRoomTagItem {
	return &liveroomdto.AppLiveRoomTagItem{
		ID:        strconv.FormatUint(tag.ID, 10),
		Name:      tag.Name,
		Sort:      tag.Sort,
		IsSpecial: tag.IsSpecial,
	}
}
