package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/syndb"
)

func writeIncomeSettlementLogAmounts(tb db.TbName, id uint64, a *LiveRoomIncomeAmounts, salary float64) {
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalIncome, id, a.TotalIncome)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalSocialIncome, id, a.TotalSocialIncome)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalGiftIncome, id, a.TotalGiftIncome)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalPaidDanmakuIncome, id, a.TotalPaidDanmakuIncome)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalVideoCallIncome, id, a.TotalVideoCallIncome)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalVideoCallTicketIncome, id, a.TotalVideoCallTicketIncome)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalVideoCallBillingIncome, id, a.TotalVideoCallBillingIncome)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalShortVideoIncome, id, a.TotalShortVideoIncome)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalGameIncome, id, a.TotalGameIncome)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalLiveDuration, id, a.TotalLiveDuration)
	writeIncomeAmountLocked(tb, LiveRoomIncomeSettlementSalary, id, salary)
}

func writeGuildIncomeSettlementLogAmounts(id uint64, a *LiveRoomIncomeAmounts, salary, shareAmount, shareAmountUsd, guildSharePercent float64) {
	writeIncomeSettlementLogAmounts(TbGuildIncomeSettlementDetail, id, a, salary)
	writeIncomeAmountLocked(TbGuildIncomeSettlementDetail, LiveRoomIncomeSettlementShareAmount, id, shareAmount)
	writeIncomeAmountLocked(TbGuildIncomeSettlementDetail, LiveRoomIncomeSettlementShareAmountUsd, id, shareAmountUsd)
	writeIncomeAmountLocked(TbGuildIncomeSettlementDetail, GuildIncomeSettlementLogGuildSharePercent, id, guildSharePercent)
}

func writeGuildIncomeSettlementBreakdown(id uint64, b *GuildIncomeSettlementBreakdown) {
	if b == nil {
		return
	}
	syndb.AddData(TbGuildIncomeSettlementDetail, GuildIncomeSettlementLogSettlementRuleType, &syndb.ColData{IdVal: id, ColVal: b.SettlementRuleType})
	writeIncomeAmountLocked(TbGuildIncomeSettlementDetail, GuildIncomeSettlementLogAnchorSocialShareAmount, id, b.AnchorSocialShareAmount)
	writeIncomeAmountLocked(TbGuildIncomeSettlementDetail, GuildIncomeSettlementLogGuildSocialShareAmount, id, b.GuildSocialShareAmount)
	writeIncomeAmountLocked(TbGuildIncomeSettlementDetail, GuildIncomeSettlementLogAnchorGameShareAmountGold, id, b.AnchorGameShareAmountGold)
	writeIncomeAmountLocked(TbGuildIncomeSettlementDetail, GuildIncomeSettlementLogGuildGameShareAmountGold, id, b.GuildGameShareAmountGold)
}

func writeAnchorIncomeSettlementLogAmounts(id uint64, a *LiveRoomIncomeAmounts, salary, shareAmount, shareAmountUsd, anchorSharePercent float64) {
	writeIncomeSettlementLogAmounts(TbAnchorIncomeSettlementLog, id, a, salary)
	writeIncomeAmountLocked(TbAnchorIncomeSettlementLog, LiveRoomIncomeSettlementShareAmount, id, shareAmount)
	writeIncomeAmountLocked(TbAnchorIncomeSettlementLog, LiveRoomIncomeSettlementShareAmountUsd, id, shareAmountUsd)
	writeIncomeAmountLocked(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogAnchorSharePercent, id, anchorSharePercent)
}

func writeAnchorIncomeSettlementBreakdown(id uint64, b *AnchorIncomeSettlementBreakdown) {
	if b == nil {
		return
	}
	syndb.AddData(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogSettlementRuleType, &syndb.ColData{IdVal: id, ColVal: b.SettlementRuleType})
	syndb.AddData(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogHasSalary, &syndb.ColData{IdVal: id, ColVal: b.HasSalary})
	writeIncomeAmountLocked(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogAnchorSocialSharePercent, id, b.AnchorSocialSharePercent)
	writeIncomeAmountLocked(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogGuildSocialSharePercent, id, b.GuildSocialSharePercent)
	writeIncomeAmountLocked(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogAnchorGameSharePercent, id, b.AnchorGameSharePercent)
	writeIncomeAmountLocked(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogGuildGameSharePercent, id, b.GuildGameSharePercent)
	writeIncomeAmountLocked(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogAnchorSocialShareAmount, id, b.AnchorSocialShareAmount)
	writeIncomeAmountLocked(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogGuildSocialShareAmount, id, b.GuildSocialShareAmount)
	writeIncomeAmountLocked(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogAnchorGameShareAmountGold, id, b.AnchorGameShareAmountGold)
	writeIncomeAmountLocked(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogGuildGameShareAmountGold, id, b.GuildGameShareAmountGold)
}
