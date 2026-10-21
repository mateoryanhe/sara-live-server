package cfgdao

import "xr-game-server/entity/user"

func resolveLuckyGiftSettlementPercent(cfg *entity.WalletExchangeCfg) float64 {
	if cfg == nil || cfg.LuckyGiftSettlementPercent < 0 || cfg.LuckyGiftSettlementPercent > MaxLuckyGiftSettlementPercent {
		return DefaultLuckyGiftSettlementPercent
	}
	return cfg.LuckyGiftSettlementPercent
}

// LuckyGiftSettlementPercent 幸运礼物结算比例(百分比); 无效或未配置时默认 10。
func LuckyGiftSettlementPercent() float64 {
	return resolveLuckyGiftSettlementPercent(GetWalletExchangeCfgCached())
}
