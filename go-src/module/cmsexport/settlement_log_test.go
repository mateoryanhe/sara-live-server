package cmsexport

import (
	"testing"
	"time"

	liveentity "xr-game-server/entity/live"
)

func TestAnchorSettlementLogToCSVRowIncludesFullBreakdown(t *testing.T) {
	row := testAnchorSettlementLog()
	got := anchorSettlementLogToCSVRow(row, map[uint64]string{row.RoomId: "anchor"})
	if len(got) != 27 {
		t.Fatalf("unexpected column count: got %d want 27", len(got))
	}
	checks := map[int]string{
		0:  "\t9007199254740993",
		1:  "6001000001",
		2:  "anchor",
		13: "11",
		14: "12",
		15: "13",
		16: "14",
		17: "true",
		18: "15",
		19: "16",
		24: "21",
		25: "22",
		26: "2026-09-24 10:11:12",
	}
	for index, want := range checks {
		if got[index] != want {
			t.Errorf("column %d: got %q want %q", index, got[index], want)
		}
	}
}

func TestPlatformAnchorPayoutToCSVRowIncludesPayoutAndFlowDetails(t *testing.T) {
	row := testAnchorSettlementLog()
	got := platformAnchorPayoutToCSVRow(row, map[uint64]string{row.RoomId: "anchor"})
	if len(got) != 35 {
		t.Fatalf("unexpected column count: got %d want 35", len(got))
	}
	checks := map[int]string{
		0:  "\t9007199254740993",
		1:  "2026-09-24 10:11:12",
		2:  "6001000001",
		3:  "anchor",
		4:  "23",
		5:  "2",
		6:  "USD",
		7:  "24",
		8:  "order-1",
		9:  "platform-1",
		10: "2026-09-24 12:13:14",
		11: "failed",
		22: "11",
		23: "12",
		24: "13",
		25: "14",
		26: "true",
		34: "22",
	}
	for index, want := range checks {
		if got[index] != want {
			t.Errorf("column %d: got %q want %q", index, got[index], want)
		}
	}
}

func TestIntersectExportUint64Ids(t *testing.T) {
	got := intersectExportUint64Ids([]uint64{1, 2, 3}, []uint64{2, 4})
	if len(got) != 1 || got[0] != 2 {
		t.Fatalf("unexpected intersection: %#v", got)
	}
}

func testAnchorSettlementLog() *liveentity.AnchorIncomeSettlementLog {
	createdAt := time.Date(2026, 9, 24, 10, 11, 12, 0, time.Local)
	transferAt := time.Date(2026, 9, 24, 12, 13, 14, 0, time.Local)
	row := &liveentity.AnchorIncomeSettlementLog{
		CreatedAt: createdAt,
		RoomId:    6001000001,
		LiveRoomIncomeAmounts: liveentity.LiveRoomIncomeAmounts{
			TotalIncome:                 1,
			TotalSocialIncome:           2,
			TotalGiftIncome:             3,
			TotalPaidDanmakuIncome:      4,
			TotalVideoCallIncome:        5,
			TotalVideoCallTicketIncome:  6,
			TotalVideoCallBillingIncome: 7,
			TotalShortVideoIncome:       8,
			TotalGameIncome:             9,
			TotalLiveDuration:           600,
		},
		SettlementSalary:          11,
		SettlementShareAmount:     12,
		SettlementShareAmountUsd:  13,
		AnchorSharePercent:        14,
		HasSalary:                 true,
		AnchorSocialSharePercent:  15,
		AnchorSocialShareAmount:   16,
		GuildSocialSharePercent:   17,
		GuildSocialShareAmount:    18,
		AnchorGameSharePercent:    19,
		AnchorGameShareAmountGold: 20,
		GuildGameSharePercent:     21,
		GuildGameShareAmountGold:  22,
		SettlementReceivableUsd:   23,
		Status:                    liveentity.AnchorIncomeSettlementStatusTransferred,
		TransferAt:                &transferAt,
		TransferOrderId:           "order-1",
		TransferPlatformNo:        "platform-1",
		TransferLocalAmount:       24,
		TransferCurrency:          "USD",
		TransferFailMsg:           "failed",
	}
	row.ID = 9007199254740993
	return row
}
