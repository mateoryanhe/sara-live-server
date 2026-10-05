package cfgdao

import "xr-game-server/entity/user"

func resolvePlatformAnchorMinimumSettlementUsd(cfg *entity.WalletExchangeCfg) float64 {
	if cfg == nil || cfg.PlatformAnchorMinimumSettlementUsd <= 0 || cfg.PlatformAnchorMinimumSettlementUsd > MaxPlatformAnchorMinimumSettlementUsd {
		return DefaultPlatformAnchorMinimumSettlementUsd
	}
	return cfg.PlatformAnchorMinimumSettlementUsd
}

// PlatformAnchorMinimumSettlementUsd 返回统一钱包配置缓存中的平台主播最低结算金额。
// 数据库未配置或值无效时使用5 USD。
func PlatformAnchorMinimumSettlementUsd() float64 {
	return resolvePlatformAnchorMinimumSettlementUsd(GetWalletExchangeCfgCached())
}
