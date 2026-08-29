package activity

const defaultFirstRechargeRatio = 20

func normalizeFirstRechargeRatio(ratio float64) float64 {
	if ratio < 0 {
		return 0
	}
	return ratio
}

func effectiveFirstRechargeRatio() float64 {
	snap := getCfgCache()
	if snap == nil || snap.ID == 0 {
		return defaultFirstRechargeRatio
	}
	return normalizeFirstRechargeRatio(snap.FirstRechargeRatio)
}

// ConfiguredFirstRechargeRatio 首充赠送比例(%);无 CMS 配置时用默认值,不受活动开关影响
func ConfiguredFirstRechargeRatio() float64 {
	return effectiveFirstRechargeRatio()
}

// ApplyFirstRechargeBonus 按配置比例加赠金币,如 20% 则 baseGold*1.2;不受活动开关影响
func ApplyFirstRechargeBonus(baseGold float64) float64 {
	if baseGold <= 0 {
		return baseGold
	}
	ratio := effectiveFirstRechargeRatio()
	if ratio <= 0 {
		return baseGold
	}
	return baseGold * (1 + ratio/100)
}
