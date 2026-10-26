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

// incomeAmountOnlyCSVCells 仅收益金额列(不含时长),与 CMS 每日流水列表导出表头一致
func incomeAmountOnlyCSVCells(amounts *liveentity.LiveRoomIncomeAmounts) []string {
	if amounts == nil {
		return []string{"", "", "", "", "", "", "", "", ""}
	}
	return []string{
		formatCSVFloat(amounts.TotalIncome),
		formatCSVFloat(amounts.TotalGameIncome),
		formatCSVFloat(amounts.TotalSocialIncome),
		formatCSVFloat(amounts.TotalGiftIncome),
		formatCSVFloat(amounts.TotalPaidDanmakuIncome),
		formatCSVFloat(amounts.TotalVideoCallIncome),
		formatCSVFloat(amounts.TotalVideoCallTicketIncome),
		formatCSVFloat(amounts.TotalVideoCallBillingIncome),
		formatCSVFloat(amounts.TotalShortVideoIncome),
	}
}

func incomeAmountCSVCells(amounts *liveentity.LiveRoomIncomeAmounts) []string {
	cells := incomeAmountOnlyCSVCells(amounts)
	if amounts == nil {
		return append(cells, "")
	}
	return append(cells, formatLiveDurationMinutes(amounts.TotalLiveDuration))
}

func incomeAmountItemCSVCells(item accountdto.LiveRoomIncomeAmountsItem) []string {
	return []string{
		formatCSVFloat(item.TotalIncome),
		formatCSVFloat(item.TotalGameIncome),
		formatCSVFloat(item.TotalSocialIncome),
		formatCSVFloat(item.TotalGiftIncome),
		formatCSVFloat(item.TotalPaidDanmakuIncome),
		formatCSVFloat(item.TotalVideoCallIncome),
		formatCSVFloat(item.TotalVideoCallTicketIncome),
		formatCSVFloat(item.TotalVideoCallBillingIncome),
		formatCSVFloat(item.TotalShortVideoIncome),
		formatLiveDurationMinutes(item.TotalLiveDuration),
	}
}

func settlementIncomeCSVCells(item *accountdto.LiveRoomIncomeAmountsItem) []string {
	if item == nil {
		return incomeAmountCSVCells(nil)
	}
	return incomeAmountItemCSVCells(*item)
}
