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
	writeIncomeAmountLocked(tb, LiveRoomIncomeTotalLiveDuration, id, a.TotalLiveDuration)
	writeIncomeAmountLocked(tb, LiveRoomIncomeSettlementSalary, id, salary)
}
