package currency

// Reason 货币变动原因枚举
type Reason uint8

const (
	// ReasonUnknown 未指定原因
	ReasonUnknown Reason = 0
	// ReasonRecharge 充值
	ReasonRecharge Reason = 1
	// ReasonShopBuy 商城购买消耗
	ReasonShopBuy Reason = 2
	// ReasonTaskReward 任务奖励
	ReasonTaskReward Reason = 3
	// ReasonRankReward 排行榜奖励
	ReasonRankReward Reason = 4
	// ReasonMailAttachment 邮件附件
	ReasonMailAttachment Reason = 5
	// ReasonGmAdjust GM 后台调整
	ReasonGmAdjust Reason = 6
	// ReasonGuildContribute 工会贡献
	ReasonGuildContribute Reason = 7
	// ReasonRefund 退款返还
	ReasonRefund Reason = 8
	// ReasonSystemGrant 系统发放(活动/补偿)
	ReasonSystemGrant Reason = 9
)

// String 返回枚举的英文标识(用于日志/调试,不参与多语言展示;
// 面向用户的本地化文案请使用 Text(lang))
func (r Reason) String() string {
	switch r {
	case ReasonRecharge:
		return "Recharge"
	case ReasonShopBuy:
		return "ShopBuy"
	case ReasonTaskReward:
		return "TaskReward"
	case ReasonRankReward:
		return "RankReward"
	case ReasonMailAttachment:
		return "MailAttachment"
	case ReasonGmAdjust:
		return "GmAdjust"
	case ReasonGuildContribute:
		return "GuildContribute"
	case ReasonRefund:
		return "Refund"
	case ReasonSystemGrant:
		return "SystemGrant"
	default:
		return "Unknown"
	}
}
