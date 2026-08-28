package liverevenue

// Type 直播间社交流水类型
type Type uint8

const (
	Gift                     Type = 1 // 礼物
	PaidDanmaku              Type = 2 // 付费弹幕
	PrivateRoom              Type = 4 // 私密直播间按时计费
	Ticket                   Type = 5 // 私密直播间门票
	LiveRoomVideoCallTicket  Type = 6 // 直播间视频通话门票
	LiveRoomVideoCallBilling Type = 7 // 直播间视频通话计费
)

// IsValid 是否为合法流水类型
func IsValid(t Type) bool {
	switch t {
	case Gift, PaidDanmaku, PrivateRoom, Ticket, LiveRoomVideoCallTicket, LiveRoomVideoCallBilling:
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
	case PrivateRoom:
		return "私密直播间计费"
	case Ticket:
		return "直播间门票"
	case LiveRoomVideoCallTicket:
		return "直播间视频通话门票"
	case LiveRoomVideoCallBilling:
		return "直播间视频通话计费"
	default:
		return "未知"
	}
}
