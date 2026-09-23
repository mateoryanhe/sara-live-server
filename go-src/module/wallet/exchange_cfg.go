package wallet

import (
	"xr-game-server/dao/cfgdao"
)

const (
	DefaultGoldToDiamondRate  = 100
	DefaultExchangeFeePercent = 3.0
	DefaultUsdToGoldRate      = 100
)

type ExchangeCfgSnapshot struct {
	GoldToDiamondRate  int
	ExchangeFeePercent float64
	UsdToGoldRate      int
}

func GetExchangeCfgSnapshot() ExchangeCfgSnapshot {
	row := cfgdao.GetWalletExchangeCfgCached()
	if row == nil {
		return ExchangeCfgSnapshot{
			GoldToDiamondRate:  DefaultGoldToDiamondRate,
			ExchangeFeePercent: DefaultExchangeFeePercent,
			UsdToGoldRate:      DefaultUsdToGoldRate,
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
	usdToGold := row.UsdToGoldRate
	if usdToGold <= 0 {
		usdToGold = DefaultUsdToGoldRate
	}
	return ExchangeCfgSnapshot{
		GoldToDiamondRate:  rate,
		ExchangeFeePercent: fee,
		UsdToGoldRate:      usdToGold,
	}
}

func GetGoldToDiamondRate() int {
	return GetExchangeCfgSnapshot().GoldToDiamondRate
}

// CalcDiamondToUsd 钻石折算美金(按钱包兑换配置: 1USD=N金币, 1金币=M钻石)
func CalcDiamondToUsd(diamondAmount float64) float64 {
	if diamondAmount <= 0 {
		return 0
	}
	snap := GetExchangeCfgSnapshot()
	return CalcDiamondToUsdWithSnapshot(diamondAmount, snap)
}

// CalcDiamondToUsdWithSnapshot 使用指定的钱包汇率快照换算，供结算审计复用。
func CalcDiamondToUsdWithSnapshot(diamondAmount float64, snap ExchangeCfgSnapshot) float64 {
	if diamondAmount <= 0 {
		return 0
	}
	denom := float64(snap.GoldToDiamondRate) * float64(snap.UsdToGoldRate)
	if denom <= 0 {
		return 0
	}
	return diamondAmount / denom
}

// CalcSettlementUsdWithSnapshot 按结算时汇率把金币分佣换算为钻石，
// 再与已有钻石金额合并计算审核用 USD。原始金币流水仍按金币保存。
func CalcSettlementUsdWithSnapshot(baseDiamond, gameGold float64, snap ExchangeCfgSnapshot) (gameDiamond, totalDiamond, usd float64) {
	if snap.GoldToDiamondRate <= 0 || snap.UsdToGoldRate <= 0 {
		return 0, 0, 0
	}
	gameDiamond = gameGold * float64(snap.GoldToDiamondRate)
	totalDiamond = baseDiamond + gameDiamond
	usd = CalcDiamondToUsdWithSnapshot(totalDiamond, snap)
	return gameDiamond, totalDiamond, usd
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
