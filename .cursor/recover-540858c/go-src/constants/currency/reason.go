package currency

// Reason 货币变动原因枚举
type Reason uint8

const (
	// ReasonUnknown 未指定原因
	ReasonUnknown Reason = 0
	// ReasonRecharge 充值(历史数据兼容,新流水请使用下方细分枚举)
	ReasonRecharge Reason = 1
	// ReasonShopBuy 商城购买消耗
	ReasonShopBuy Reason = 2
	// ReasonTaskReward 任务奖励
	ReasonTaskReward Reason = 3
	// ReasonRankReward 排行榜奖励
	ReasonRankReward Reason = 4
	// ReasonMailAttachment 邮件附件
	ReasonMailAttachment Reason = 5
	// ReasonGmAdjustTest CMS 后台调整-测试
	ReasonGmAdjustTest Reason = 6
	// ReasonGmAdjustCompensation CMS 后台调整-GM补偿
	ReasonGmAdjustCompensation Reason = 7
	// ReasonGuildContribute 工会贡献
	ReasonGuildContribute Reason = 8
	// ReasonRefund 退款返还
	ReasonRefund Reason = 9
	// ReasonSystemGrant 系统发放(活动/补偿)
	ReasonSystemGrant Reason = 10
	// ReasonGiftSend 直播间送礼消耗
	ReasonGiftSend Reason = 11
	// ReasonGoldExchange 金币兑换钻石
	ReasonGoldExchange Reason = 12
	// ReasonShortVideoWatch 短视频观看扣费
	ReasonShortVideoWatch Reason = 13
	// ReasonAnchorGiftRevenue 主播收到礼物收益
	ReasonAnchorGiftRevenue Reason = 14
	// ReasonPrivateRoomTicket 私密直播间门票
	ReasonPrivateRoomTicket Reason = 15
	// ReasonPrivateRoomBilling 私密直播间按分钟计费
	ReasonPrivateRoomBilling Reason = 16
	// ReasonPaidDanmaku 直播间付费弹幕消耗
	ReasonPaidDanmaku Reason = 17
	// ReasonLiveRoomVideoCallTicket 直播间视频通话门票
	ReasonLiveRoomVideoCallTicket Reason = 18
	// ReasonLiveRoomVideoCallBilling 直播间视频通话计费
	ReasonLiveRoomVideoCallBilling Reason = 19
	// ReasonGameBet 游戏消费扣款
	ReasonGameBet Reason = 20
	// ReasonGameBetWin 游戏奖励
	ReasonGameBetWin Reason = 21

	// ReasonRechargeFirstBonus 档位首充(含加赠)
	ReasonRechargeFirstBonus Reason = 22
	// ReasonRechargeTier 档位充值
	ReasonRechargeTier Reason = 23
	// ReasonRechargeCustom 自定义金额充值
	ReasonRechargeCustom Reason = 24
	// ReasonRechargeGoogle Google Play 充值
	ReasonRechargeGoogle Reason = 25
	// ReasonRechargeIOS iOS 充值
	ReasonRechargeIOS Reason = 26
	// ReasonRechargeChannel 渠道充值
	ReasonRechargeChannel Reason = 27
	// ReasonRechargeWhitelist 白名单充值(测试)
	ReasonRechargeWhitelist Reason = 28
	// ReasonRechargeManual 后台手动充值
	ReasonRechargeManual Reason = 29
	// ReasonShortVideoAuthorSettlement 非主播作者短视频周结算到账
	ReasonShortVideoAuthorSettlement Reason = 30
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
	case ReasonGmAdjustTest:
		return "GmAdjustTest"
	case ReasonGmAdjustCompensation:
		return "GmAdjustCompensation"
	case ReasonGuildContribute:
		return "GuildContribute"
	case ReasonRefund:
		return "Refund"
	case ReasonSystemGrant:
		return "SystemGrant"
	case ReasonGiftSend:
		return "GiftSend"
	case ReasonGoldExchange:
		return "GoldExchange"
	case ReasonShortVideoWatch:
		return "ShortVideoWatch"
	case ReasonAnchorGiftRevenue:
		return "AnchorGiftRevenue"
	case ReasonPrivateRoomTicket:
		return "PrivateRoomTicket"
	case ReasonPrivateRoomBilling:
		return "PrivateRoomBilling"
	case ReasonPaidDanmaku:
		return "PaidDanmaku"
	case ReasonLiveRoomVideoCallTicket:
		return "LiveRoomVideoCallTicket"
	case ReasonLiveRoomVideoCallBilling:
		return "LiveRoomVideoCallBilling"
	case ReasonGameBet:
		return "GameBet"
	case ReasonGameBetWin:
		return "GameBetWin"
	case ReasonRechargeFirstBonus:
		return "RechargeFirstBonus"
	case ReasonRechargeTier:
		return "RechargeTier"
	case ReasonRechargeCustom:
		return "RechargeCustom"
	case ReasonRechargeGoogle:
		return "RechargeGoogle"
	case ReasonRechargeIOS:
		return "RechargeIOS"
	case ReasonRechargeChannel:
		return "RechargeChannel"
	case ReasonRechargeWhitelist:
		return "RechargeWhitelist"
	case ReasonRechargeManual:
		return "RechargeManual"
	case ReasonShortVideoAuthorSettlement:
		return "ShortVideoAuthorSettlement"
	default:
		return "Unknown"
	}
}
