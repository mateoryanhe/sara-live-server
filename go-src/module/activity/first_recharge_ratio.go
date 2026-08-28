package activity

const defaultFirstRechargeRatio = 20

func normalizeFirstRechargeRatio(ratio float64) float64 {
	if ratio < 0 {
		return 0
	}
	return ratio
}

// ConfiguredFirstRechargeRatio 首充活动配置比例(无配置时返回0;不受活动开关影响,开关仅 App 展示用)
func ConfiguredFirstRechargeRatio() float64 {
	snap := getCfgCache()
	if snap == nil || snap.ID == 0 {
		return 0
	}
	return normalizeFirstRechargeRatio(snap.FirstRechargeRatio)
}

// ApplyFirstRechargeBonus 按配置比例加赠金币,如 20% 则 baseGold*1.2;不受活动开关影响
func ApplyFirstRechargeBonus(baseGold float64) float64 {
	if baseGold <= 0 {
		return baseGold
	}
	snap := getCfgCache()
	if snap == nil || snap.ID == 0 || snap.FirstRechargeRatio <= 0 {
		return baseGold
	}
	return baseGold * (1 + snap.FirstRechargeRatio/100)
}
