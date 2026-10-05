package cfgdao

import "xr-game-server/entity/user"

func resolveEffectiveLiveMinSessionMinutes(cfg *entity.WalletExchangeCfg) int {
	if cfg == nil || cfg.EffectiveLiveMinSessionMinutes <= 0 || cfg.EffectiveLiveMinSessionMinutes > MaxEffectiveLiveMinSessionMinutes {
		return DefaultEffectiveLiveMinSessionMinutes
	}
	return cfg.EffectiveLiveMinSessionMinutes
}

// EffectiveLiveMinSessionMinutes 返回统一钱包配置缓存中的有效直播门槛。
// 数据库未配置或值无效时返回30分钟。
func EffectiveLiveMinSessionMinutes() int {
	return resolveEffectiveLiveMinSessionMinutes(GetWalletExchangeCfgCached())
}
