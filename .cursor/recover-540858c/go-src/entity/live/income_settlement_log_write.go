package entity

import "xr-game-server/constants/db"

func writeIncomeSettlementLogAmounts(tb db.TbName, id uint64, a *LiveRoomIncomeAmounts, salary float64) {
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalIncome, id, a.TotalIncome)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalGiftIncome, id, a.TotalGiftIncome)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalPaidDanmakuIncome, id, a.TotalPaidDanmakuIncome)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalPrivateRoomTicketIncome, id, a.TotalPrivateRoomTicketIncome)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalPrivateRoomWatchIncome, id, a.TotalPrivateRoomWatchIncome)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalVideoCallIncome, id, a.TotalVideoCallIncome)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalVideoCallTicketIncome, id, a.TotalVideoCallTicketIncome)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalVideoCallBillingIncome, id, a.TotalVideoCallBillingIncome)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalShortVideoIncome, id, a.TotalShortVideoIncome)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalLiveDuration, id, a.TotalLiveDuration)
	writeIncomeAmountLocked(tb, LiveRoomIncomeSettlementSalary, id, salary)
}

func writeGuildIncomeSettlementLogAmounts(id uint64, a *LiveRoomIncomeAmounts, shareAmount, guildSharePercent float64) {
	writeIncomeSettlementLogAmounts(TbGuildIncomeSettlementLog, id, a, 0)
	writeIncomeAmountLocked(TbGuildIncomeSettlementLog, LiveRoomIncomeSettlementShareAmount, id, shareAmount)
	writeIncomeAmountLocked(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogGuildSharePercent, id, guildSharePercent)
}

func writeAnchorIncomeSettlementLogAmounts(id uint64, a *LiveRoomIncomeAmounts, salary, shareAmount, anchorSharePercent float64) {
	writeIncomeSettlementLogAmounts(TbAnchorIncomeSettlementLog, id, a, salary)
	writeIncomeAmountLocked(TbAnchorIncomeSettlementLog, LiveRoomIncomeSettlementShareAmount, id, shareAmount)
	writeIncomeAmountLocked(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogAnchorSharePercent, id, anchorSharePercent)
}
