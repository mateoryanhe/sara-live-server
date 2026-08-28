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
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalGameIncome, id, a.TotalGameIncome)
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalLiveDuration, id, a.TotalLiveDuration)
	writeIncomeAmountLocked(tb, LiveRoomIncomeSettlementSalary, id, salary)
}

func writeGuildIncomeSettlementLogAmounts(id uint64, a *LiveRoomIncomeAmounts, salary, shareAmount, shareAmountUsd, receivableUsd, guildSharePercent float64) {
	writeIncomeSettlementLogAmounts(TbGuildIncomeSettlementLog, id, a, salary)
	writeIncomeAmountLocked(TbGuildIncomeSettlementLog, LiveRoomIncomeSettlementShareAmount, id, shareAmount)
	writeIncomeAmountLocked(TbGuildIncomeSettlementLog, LiveRoomIncomeSettlementShareAmountUsd, id, shareAmountUsd)
	writeIncomeAmountLocked(TbGuildIncomeSettlementLog, LiveRoomIncomeSettlementReceivableUsd, id, receivableUsd)
	writeIncomeAmountLocked(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogGuildSharePercent, id, guildSharePercent)
}

func writeAnchorIncomeSettlementLogAmounts(id uint64, a *LiveRoomIncomeAmounts, salary, shareAmount, shareAmountUsd, anchorSharePercent float64) {
	writeIncomeSettlementLogAmounts(TbAnchorIncomeSettlementLog, id, a, salary)
	writeIncomeAmountLocked(TbAnchorIncomeSettlementLog, LiveRoomIncomeSettlementShareAmount, id, shareAmount)
	writeIncomeAmountLocked(TbAnchorIncomeSettlementLog, LiveRoomIncomeSettlementShareAmountUsd, id, shareAmountUsd)
	writeIncomeAmountLocked(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogAnchorSharePercent, id, anchorSharePercent)
}
