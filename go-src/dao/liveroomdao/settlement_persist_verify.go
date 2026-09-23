package liveroomdao

import (
	"math"

	"github.com/gogf/gf/v2/frame/g"
	live "xr-game-server/entity/live"
)

func settlementFloatEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.0001
}

// VerifyAnchorIncomeSettlementLogPersisted 结算清零前确认关键分项已完整落库。
func VerifyAnchorIncomeSettlementLogPersisted(expected *live.AnchorIncomeSettlementLog) bool {
	if expected == nil || expected.ID == 0 {
		return false
	}
	var row live.AnchorIncomeSettlementLog
	if err := g.DB().Model(string(live.TbAnchorIncomeSettlementLog)).WherePri(expected.ID).Scan(&row); err != nil || row.ID == 0 {
		return false
	}
	return row.RoomId == expected.RoomId &&
		row.DirectPayout == expected.DirectPayout &&
		row.SettlementRuleType == expected.SettlementRuleType &&
		row.HasSalary == expected.HasSalary &&
		settlementFloatEqual(row.TotalSocialIncome, expected.TotalSocialIncome) &&
		settlementFloatEqual(row.TotalGameIncome, expected.TotalGameIncome) &&
		settlementFloatEqual(row.SettlementSalary, expected.SettlementSalary) &&
		settlementFloatEqual(row.AnchorSocialShareAmount, expected.AnchorSocialShareAmount) &&
		settlementFloatEqual(row.GuildSocialShareAmount, expected.GuildSocialShareAmount) &&
		settlementFloatEqual(row.AnchorGameShareAmountGold, expected.AnchorGameShareAmountGold) &&
		settlementFloatEqual(row.GuildGameShareAmountGold, expected.GuildGameShareAmountGold)
}

// VerifyAnchorPayoutConversionPersisted 代付外呼前确认平台主播换算快照与USD金额已落库。
func VerifyAnchorPayoutConversionPersisted(expected *live.AnchorIncomeSettlementLog) bool {
	if expected == nil || expected.ID == 0 {
		return false
	}
	var row live.AnchorIncomeSettlementLog
	if err := g.DB().Model(string(live.TbAnchorIncomeSettlementLog)).WherePri(expected.ID).Scan(&row); err != nil || row.ID == 0 {
		return false
	}
	return row.DirectPayout == expected.DirectPayout &&
		row.GoldToDiamondRate == expected.GoldToDiamondRate &&
		row.UsdToGoldRate == expected.UsdToGoldRate &&
		settlementFloatEqual(row.GameShareAmountDiamond, expected.GameShareAmountDiamond) &&
		settlementFloatEqual(row.TotalSettlementDiamond, expected.TotalSettlementDiamond) &&
		settlementFloatEqual(row.SettlementReceivableUsd, expected.SettlementReceivableUsd)
}

// VerifyGuildIncomeSettlementLogPersisted 工会未结算清零前确认关键分项已完整落库。
func VerifyGuildIncomeSettlementLogPersisted(expected *live.GuildIncomeSettlementLog) bool {
	if expected == nil || expected.ID == 0 {
		return false
	}
	var row live.GuildIncomeSettlementLog
	if err := g.DB().Model(string(live.TbGuildIncomeSettlementLog)).WherePri(expected.ID).Scan(&row); err != nil || row.ID == 0 {
		return false
	}
	var detail live.GuildIncomeSettlementDetail
	if err := g.DB().Model(string(live.TbGuildIncomeSettlementDetail)).
		Where(string(live.GuildIncomeSettlementDetailSettlementId)+" = ?", expected.ID).
		Scan(&detail); err != nil || detail.SettlementId == 0 {
		return false
	}
	return row.GuildId == expected.GuildId &&
		detail.SettlementRuleType == expected.SettlementRuleType &&
		settlementFloatEqual(detail.TotalSocialIncome, expected.TotalSocialIncome) &&
		settlementFloatEqual(detail.TotalGameIncome, expected.TotalGameIncome) &&
		settlementFloatEqual(detail.SettlementSalary, expected.SettlementSalary) &&
		settlementFloatEqual(detail.AnchorSocialShareAmount, expected.AnchorSocialShareAmount) &&
		settlementFloatEqual(detail.GuildSocialShareAmount, expected.GuildSocialShareAmount) &&
		settlementFloatEqual(detail.AnchorGameShareAmountGold, expected.AnchorGameShareAmountGold) &&
		settlementFloatEqual(detail.GuildGameShareAmountGold, expected.GuildGameShareAmountGold)
}

// VerifyGuildPayoutConversionPersisted 代付外呼前确认换算快照与USD金额已落库。
func VerifyGuildPayoutConversionPersisted(expected *live.GuildIncomeSettlementLog) bool {
	if expected == nil || expected.ID == 0 {
		return false
	}
	var row live.GuildIncomeSettlementLog
	if err := g.DB().Model(string(live.TbGuildIncomeSettlementLog)).WherePri(expected.ID).Scan(&row); err != nil || row.ID == 0 {
		return false
	}
	var detail live.GuildIncomeSettlementDetail
	if err := g.DB().Model(string(live.TbGuildIncomeSettlementDetail)).
		Where(string(live.GuildIncomeSettlementDetailSettlementId)+" = ?", expected.ID).
		Scan(&detail); err != nil || detail.SettlementId == 0 {
		return false
	}
	return detail.GoldToDiamondRate == expected.GoldToDiamondRate &&
		detail.UsdToGoldRate == expected.UsdToGoldRate &&
		settlementFloatEqual(detail.GameShareAmountDiamond, expected.GameShareAmountDiamond) &&
		settlementFloatEqual(detail.TotalSettlementDiamond, expected.TotalSettlementDiamond) &&
		settlementFloatEqual(row.SettlementReceivableUsd, expected.SettlementReceivableUsd)
}
