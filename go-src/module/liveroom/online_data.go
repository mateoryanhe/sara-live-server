package liveroom

import (
	"github.com/gogf/gf/v2/container/gmap"
	"time"
)

var onlineMap = gmap.NewKVMap[uint64, *gmap.KVMap[uint64, time.Time]](false)

func addToOnline(userId uint64, roomId uint64) {
	roomMap := onlineMap.Get(roomId)
	if roomMap == nil {
		initRoomOnline(roomId)
		roomMap = onlineMap.Get(roomId)
	}
	roomMap.Set(userId, time.Now())
}

func getLenForRoom(roomId uint64) int {
	roomMap := onlineMap.Get(roomId)
	if roomMap == nil {
		return 0
	}
	return roomMap.Size()
}

func initRoomOnline(roomId uint64) {
	if onlineMap.Contains(roomId) {
		return
	}
	onlineMap.Set(roomId, gmap.NewKVMap[uint64, time.Time](true))
}

func removeOnline(userId uint64, roomId uint64) {
	roomMap := onlineMap.Get(roomId)
	if roomMap == nil {
		return
	}
	roomMap.Remove(userId)
}

func getOnline(roomId uint64) []uint64 {
	roomMap := onlineMap.Get(roomId)
	if roomMap == nil {
		initRoomOnline(roomId)
		return make([]uint64, 0)
	}
	return roomMap.Keys()
}
