package userstatus

// UserStatus 用户当前状态。
type UserStatus uint8

const (
	UserStatusOnline  UserStatus = 1 // WebSocket 在线
	UserStatusLive    UserStatus = 2 // 正在开播
	UserStatusOffline UserStatus = 3 // 离线
)

// 直播间状态
const (
	LiveRoomStatusClosed uint8 = 0 // 未开播/已下播
	LiveRoomStatusLive   uint8 = 1
	Empty                      = 2 // 直播中
)
