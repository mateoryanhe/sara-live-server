package entity

import stdmath "math"

// CalcTieredGuildPayoutDiamondSplit 普通工会分项结算：按代付换算时金币→钻石汇率拆分代付钻石总额。
func CalcTieredGuildPayoutDiamondSplit(salary, anchorSocial, guildSocial, anchorGameGold, guildGameGold float64, goldToDiamondRate int) (anchorTotal, guildTotal float64) {
	anchorGameDiamond := calcSettlementGoldToDiamond(anchorGameGold, goldToDiamondRate)
	guildGameDiamond := calcSettlementGoldToDiamond(guildGameGold, goldToDiamondRate)
	anchorTotal = salary + anchorSocial + anchorGameDiamond
	guildTotal = guildSocial + guildGameDiamond
	anchorTotal = stdmath.Round(anchorTotal*10000) / 10000
	guildTotal = stdmath.Round(guildTotal*10000) / 10000
	return anchorTotal, guildTotal
}

func calcSettlementGoldToDiamond(gold float64, goldToDiamondRate int) float64 {
	if gold <= 0 || goldToDiamondRate <= 0 || stdmath.IsNaN(gold) || stdmath.IsInf(gold, 0) {
		return 0
	}
	converted := gold * float64(goldToDiamondRate)
	if stdmath.IsNaN(converted) || stdmath.IsInf(converted, 0) {
		return 0
	}
	return stdmath.Round(converted*10000) / 10000
}
