package wallet

import (
	"xr-game-server/dao/cfgdao"
)

const (
	DefaultGoldToDiamondRate  = 100
	DefaultExchangeFeePercent = 3.0
)

type ExchangeCfgSnapshot struct {
	GoldToDiamondRate  int
	ExchangeFeePercent float64
}

func GetExchangeCfgSnapshot() ExchangeCfgSnapshot {
	row := cfgdao.GetWalletExchangeCfgCached()
	if row == nil {
		return ExchangeCfgSnapshot{
			GoldToDiamondRate:  DefaultGoldToDiamondRate,
			ExchangeFeePercent: DefaultExchangeFeePercent,
		}
	}
	rate := row.GoldToDiamondRate
	if rate <= 0 {
		rate = DefaultGoldToDiamondRate
	}
	fee := row.ExchangeFeePercent
	if fee < 0 {
		fee = 0
	}
	return ExchangeCfgSnapshot{
		GoldToDiamondRate:  rate,
		ExchangeFeePercent: fee,
	}
}

func GetGoldToDiamondRate() int {
	return GetExchangeCfgSnapshot().GoldToDiamondRate
}

// calcAppExchangeGoldDeduct 仅 App 手动兑换使用;手续费比例>0 时额外扣金币
func calcAppExchangeGoldDeduct(goldAmount float64, snap ExchangeCfgSnapshot) float64 {
	if goldAmount <= 0 || snap.ExchangeFeePercent <= 0 {
		return goldAmount
	}
	return goldAmount * (1 + snap.ExchangeFeePercent/100)
}
