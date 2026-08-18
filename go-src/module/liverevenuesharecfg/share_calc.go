package liverevenuesharecfg

import (
	"xr-game-server/core/math"
	"xr-game-server/dao/cfgdao"
)

// ResolveAnchorSharePercent 读取主播流水分佣比例(未配置时用默认值)
func ResolveAnchorSharePercent() float64 {
	cfg := cfgdao.GetLiveRevenueShareCfgCached()
	if cfg == nil {
		return DefaultAnchorSharePercent
	}
	return cfg.AnchorSharePercent
}

// ResolveGuildSharePercent 读取工会流水分佣比例(未配置时用默认值)
func ResolveGuildSharePercent() float64 {
	cfg := cfgdao.GetLiveRevenueShareCfgCached()
	if cfg == nil {
		return DefaultGuildSharePercent
	}
	return cfg.GuildSharePercent
}

// CalcSettlementShareAmount 结算分佣金额 = 结算薪资 + 主播分佣比例% * 未结算总流水(TotalIncome)
func CalcSettlementShareAmount(salary, unsettledTotalIncome float64) float64 {
	commission := unsettledTotalIncome * ResolveAnchorSharePercent() / 100
	return math.AddFloat64(salary, commission)
}

// CalcGuildSettlementShareAmount 工会结算分佣金额 = 工会分佣比例% * 未结算总流水(TotalIncome)
func CalcGuildSettlementShareAmount(unsettledTotalIncome float64) float64 {
	return unsettledTotalIncome * ResolveGuildSharePercent() / 100
}
