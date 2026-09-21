package liverevenue

// Type 直播间社交流水类型
type Type uint8

const (
	Gift                     Type = 1 // 礼物
	PaidDanmaku              Type = 2 // 付费弹幕
	LiveRoomVideoCallTicket  Type = 6 // 直播间视频通话门票
	LiveRoomVideoCallBilling Type = 7 // 直播间视频通话计费
)

// IsValid 是否为合法流水类型
func IsValid(t Type) bool {
	switch t {
	case Gift, PaidDanmaku, LiveRoomVideoCallTicket, LiveRoomVideoCallBilling:
		return true
	default:
		return false
	}
}

// Text 流水类型文案
func Text(t Type) string {
	switch t {
	case Gift:
		return "礼物"
	case PaidDanmaku:
		return "付费弹幕"
	case LiveRoomVideoCallTicket:
		return "直播间视频通话门票"
	case LiveRoomVideoCallBilling:
		return "直播间视频通话计费"
	default:
		return "未知"
	}
}
