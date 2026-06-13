package liverevenue

// Type 直播收益流水类型
type Type uint8

const (
	Gift        Type = 1 // 礼物
	PaidDanmaku Type = 2 // 付费弹幕
	GameBet     Type = 3 // 游戏下注
	PrivateRoom Type = 4 // 私密直播间计费
)

// IsValid 是否为合法收益类型
func IsValid(t Type) bool {
	switch t {
	case Gift, PaidDanmaku, GameBet, PrivateRoom:
		return true
	default:
		return false
	}
}

// Text 收益类型文案
func Text(t Type) string {
	switch t {
	case Gift:
		return "礼物"
	case PaidDanmaku:
		return "付费弹幕"
	case GameBet:
		return "游戏下注"
	case PrivateRoom:
		return "私密直播间计费"
	default:
		return "未知"
	}
}
