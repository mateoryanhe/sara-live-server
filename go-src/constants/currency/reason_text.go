package currency

// reasonTextMap 货币变动原因的多语言文案表
// 新增枚举值时,请同步在每个语言映射中补全对应文案
var reasonTextMap = map[Lang]map[Reason]string{
	LangZHCN: {
		ReasonUnknown:           "未知",
		ReasonRecharge:          "充值",
		ReasonShopBuy:           "商城购买",
		ReasonTaskReward:        "任务奖励",
		ReasonRankReward:        "排行榜奖励",
		ReasonMailAttachment:    "邮件附件",
		ReasonGmAdjust:          "GM调整",
		ReasonGuildContribute:   "工会贡献",
		ReasonRefund:            "退款返还",
		ReasonSystemGrant:       "系统发放",
		ReasonGiftSend:          "直播间送礼",
		ReasonGoldExchange:      "金币兑换钻石",
		ReasonShortVideoWatch:   "短视频观看",
		ReasonAnchorGiftRevenue: "主播收到礼物收益",
		ReasonPrivateRoomTicket: "私密直播间门票",
	},
	LangZHTW: {
		ReasonUnknown:           "未知",
		ReasonRecharge:          "儲值",
		ReasonShopBuy:           "商城購買",
		ReasonTaskReward:        "任務獎勵",
		ReasonRankReward:        "排行榜獎勵",
		ReasonMailAttachment:    "郵件附件",
		ReasonGmAdjust:          "GM調整",
		ReasonGuildContribute:   "公會貢獻",
		ReasonRefund:            "退款返還",
		ReasonSystemGrant:       "系統發放",
		ReasonGiftSend:          "直播間送禮",
		ReasonGoldExchange:      "金幣兌換鑽石",
		ReasonShortVideoWatch:   "短視頻觀看",
		ReasonAnchorGiftRevenue: "主播收到禮物收益",
		ReasonPrivateRoomTicket: "私密直播間門票",
	},
	LangEN: {
		ReasonUnknown:           "Unknown",
		ReasonRecharge:          "Recharge",
		ReasonShopBuy:           "Shop Purchase",
		ReasonTaskReward:        "Task Reward",
		ReasonRankReward:        "Rank Reward",
		ReasonMailAttachment:    "Mail Attachment",
		ReasonGmAdjust:          "GM Adjustment",
		ReasonGuildContribute:   "Guild Contribution",
		ReasonRefund:            "Refund",
		ReasonSystemGrant:       "System Grant",
		ReasonGiftSend:          "Live Room Gift",
		ReasonGoldExchange:      "Gold to Diamond Exchange",
		ReasonShortVideoWatch:   "Short Video Watch",
		ReasonAnchorGiftRevenue: "Anchor Gift Revenue",
		ReasonPrivateRoomTicket: "Private Room Ticket",
	},
}

// Text 返回当前枚举在指定语言下的文案;未匹配时回落到默认语言,再不匹配则返回 String()
func (r Reason) Text(lang Lang) string {
	if m, ok := reasonTextMap[lang]; ok {
		if s, ok2 := m[r]; ok2 {
			return s
		}
	}
	if m, ok := reasonTextMap[DefaultLang]; ok {
		if s, ok2 := m[r]; ok2 {
			return s
		}
	}
	return r.String()
}
