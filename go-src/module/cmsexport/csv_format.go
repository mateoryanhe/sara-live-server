package cmsexport

import (
	"math"
	"strconv"
	"time"

	"xr-game-server/dto/accountdto"
	liveentity "xr-game-server/entity/live"
)

func formatCSVTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func formatCSVTimePtr(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func formatLiveDurationMinutes(seconds float64) string {
	if seconds <= 0 {
		return ""
	}
	minutes := math.Round(seconds/60*10) / 10
	return strconv.FormatFloat(minutes, 'f', -1, 64)
}

func parseUint64Filter(val string) uint64 {
	if val == "" {
		return 0
	}
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

func formatSettledText(settled bool, yesText, noText string) string {
	if settled {
		return yesText
	}
	return noText
}

func incomeAmountCSVCells(amounts *liveentity.LiveRoomIncomeAmounts) []string {
	if amounts == nil {
		return []string{"", "", "", "", "", "", "", "", ""}
	}
	return []string{
		formatCSVFloat(amounts.TotalIncome),
		formatCSVFloat(amounts.TotalGiftIncome),
		formatCSVFloat(amounts.TotalPaidDanmakuIncome),
		formatCSVFloat(amounts.TotalPrivateRoomTicketIncome),
		formatCSVFloat(amounts.TotalPrivateRoomWatchIncome),
		formatCSVFloat(amounts.TotalVideoCallIncome),
		formatCSVFloat(amounts.TotalVideoCallTicketIncome),
		formatCSVFloat(amounts.TotalVideoCallBillingIncome),
		formatLiveDurationMinutes(amounts.TotalLiveDuration),
	}
}

func incomeAmountItemCSVCells(item accountdto.LiveRoomIncomeAmountsItem) []string {
	return []string{
		formatCSVFloat(item.TotalIncome),
		formatCSVFloat(item.TotalGiftIncome),
		formatCSVFloat(item.TotalPaidDanmakuIncome),
		formatCSVFloat(item.TotalPrivateRoomTicketIncome),
		formatCSVFloat(item.TotalPrivateRoomWatchIncome),
		formatCSVFloat(item.TotalVideoCallIncome),
		formatCSVFloat(item.TotalVideoCallTicketIncome),
		formatCSVFloat(item.TotalVideoCallBillingIncome),
		formatLiveDurationMinutes(item.TotalLiveDuration),
	}
}

func settlementIncomeCSVCells(item *accountdto.LiveRoomIncomeAmountsItem) []string {
	if item == nil {
		return incomeAmountCSVCells(nil)
	}
	return incomeAmountItemCSVCells(*item)
}
