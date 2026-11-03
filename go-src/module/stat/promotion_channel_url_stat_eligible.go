package stat

// 推广日报暂不校验客户端平台/Header; 下单/充值是否入账由 user_ext 渠道字段与 record 内逻辑决定.

func userEligibleForPromotionChannelURLDailyStat(_ uint64) bool {
	return true
}

func userEligibleForPromotionAppDailyStat(_ uint64) bool {
	return true
}
