package roomstatus

// 直播间状态
const (
	LiveRoomStatusClosed uint8 = 0 // 未开播/已下播
	LiveRoomStatusLive   uint8 = 1
	Empty                      = 2 // 直播中
)
