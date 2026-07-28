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

func getOnline(roomId uint64) []uint64 {
	roomMap := getRoomOnlineMap(roomId, false)
	if roomMap == nil {
		return make([]uint64, 0)
	}
	return roomMap.Keys()
}

func isUserInOnlineMap(userId, roomId uint64) bool {
	roomMap := getRoomOnlineMap(roomId, false)
	return roomMap != nil && roomMap.Contains(userId)
}

func countAudienceInRoom(roomId uint64) int {
	count := len(getOnline(roomId))
	return count
}
