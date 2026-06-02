package statdto

// 用户数据柱形图指标 key(前后端约定,新增指标在此补充常量)
const (
	UserStatBarMetricRechargeUser       = "rechargeUser"       // 充值人数(去重)
	UserStatBarMetricVipUser            = "vipUser"            // VIP用户(预留)
	UserStatBarMetricGoldConsumeUser    = "goldConsumeUser"    // 金币消费用户(去重)
	UserStatBarMetricDiamondConsumeUser = "diamondConsumeUser" // 钻石消费用户(去重)
)

// BuildUserStatBarMetrics 构造柱形图指标快照(便于后续扩展)
func BuildUserStatBarMetrics(rechargeUserCount, goldConsumeUserCount, diamondConsumeUserCount uint64) map[string]uint64 {
	return map[string]uint64{
		UserStatBarMetricRechargeUser:       rechargeUserCount,
		UserStatBarMetricGoldConsumeUser:    goldConsumeUserCount,
		UserStatBarMetricDiamondConsumeUser: diamondConsumeUserCount,
	}
}
