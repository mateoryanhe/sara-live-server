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

// calcExchangeDiamond 计算兑换钻石: 毛钻石=金币*比例; App手动兑换时手续费从钻石扣(如100钻扣3%得97钻)
func calcExchangeDiamond(goldAmount float64, snap ExchangeCfgSnapshot, applyAppFee bool) (gross, fee, net float64) {
	gross = goldAmount * float64(snap.GoldToDiamondRate)
	if !applyAppFee || snap.ExchangeFeePercent <= 0 {
		return gross, 0, gross
	}
	fee = gross * (snap.ExchangeFeePercent / 100)
	net = gross - fee
	if net < 0 {
		net = 0
	}
	return gross, fee, net
}
