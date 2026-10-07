package cfgdao

import "xr-game-server/entity/user"

func resolveEffectiveLiveMinSessionMinutes(cfg *entity.WalletExchangeCfg) int {
	if cfg == nil || cfg.EffectiveLiveMinSessionMinutes <= 0 || cfg.EffectiveLiveMinSessionMinutes > MaxEffectiveLiveMinSessionMinutes {
		return DefaultEffectiveLiveMinSessionMinutes
	}
	return cfg.EffectiveLiveMinSessionMinutes
}

func resolveEffectiveLiveDailyAccumulatedMinutes(cfg *entity.WalletExchangeCfg) int {
	if cfg == nil || cfg.EffectiveLiveDailyAccumulatedMinutes < 0 || cfg.EffectiveLiveDailyAccumulatedMinutes > MaxEffectiveLiveDailyAccumulatedMinutes {
		return DefaultEffectiveLiveDailyAccumulatedMinutes
	}
	return cfg.EffectiveLiveDailyAccumulatedMinutes
}

// EffectiveLiveMinSessionMinutes 返回统一钱包配置缓存中的有效直播门槛。
// 数据库未配置或值无效时返回30分钟。
func EffectiveLiveMinSessionMinutes() int {
	return resolveEffectiveLiveMinSessionMinutes(GetWalletExchangeCfgCached())
}

// EffectiveLiveDailyAccumulatedMinutes 返回每天累计开播时长门槛(分钟)。
// 0 表示未配置；负数或超范围时按 0 处理。
func EffectiveLiveDailyAccumulatedMinutes() int {
	return resolveEffectiveLiveDailyAccumulatedMinutes(GetWalletExchangeCfgCached())
}
