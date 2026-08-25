package liveroom

import (
	"github.com/gogf/gf/v2/container/gmap"
	"time"
)

var onlineMap = gmap.NewKVMap[uint64, *gmap.KVMap[uint64, time.Time]](true)

func getRoomOnlineMap(roomId uint64, create bool) *gmap.KVMap[uint64, time.Time] {
	if create {
		return onlineMap.GetOrSetFuncLock(roomId, func() *gmap.KVMap[uint64, time.Time] {
			return gmap.NewKVMap[uint64, time.Time](true)
		})
	}
	return onlineMap.Get(roomId)
}

func addToOnline(userId uint64, roomId uint64) {
	getRoomOnlineMap(roomId, true).Set(userId, time.Now())
}

func getLenForRoom(roomId uint64) int {
	roomMap := getRoomOnlineMap(roomId, false)
	if roomMap == nil {
		return 0
	}
	return roomMap.Size()
}

func initRoomOnline(roomId uint64) {
	getRoomOnlineMap(roomId, true)
}

func removeOnline(userId uint64, roomId uint64) {
	roomMap := getRoomOnlineMap(roomId, false)
	if roomMap == nil {
		return
	}
	roomMap.Remove(userId)
}

// getOnline 获取房间在线用户ID列表;key 不存在时返回空切片,避免 nil
func getOnline(roomId uint64) []uint64 {
	roomMap := getRoomOnlineMap(roomId, false)
	if roomMap == nil {
		return make([]uint64, 0)
	}
	keys := roomMap.Keys()
	if keys == nil {
		return make([]uint64, 0)
	}
	return keys
}

func findOnlineRoomIdsByUser(userId uint64) []uint64 {
	roomIds := make([]uint64, 0)
	if userId == 0 {
		return roomIds
	}
	for _, roomId := range onlineMap.Keys() {
		if isUserInOnlineMap(userId, roomId) {
			roomIds = append(roomIds, roomId)
		}
	}
	return roomIds
}

func isUserInOnlineMap(userId, roomId uint64) bool {
	roomMap := getRoomOnlineMap(roomId, false)
	return roomMap != nil && roomMap.Contains(userId)
}

func countAudienceInRoom(roomId uint64) int {
	return len(getOnline(roomId))
}

// appAudienceOnlineCount App 端展示的在线观众数; 主播下播时固定返回 0
func appAudienceOnlineCount(roomId, liveRecordId uint64) int {
	if roomId == 0 || liveRecordId == 0 {
		return 0
	}
	return countAudienceInRoom(roomId)
}

// clearRoomOnlineMap 下播时清空该房在线名单内存
func clearRoomOnlineMap(roomId uint64) {
	if roomId == 0 {
		return
	}
	onlineMap.Remove(roomId)
}
