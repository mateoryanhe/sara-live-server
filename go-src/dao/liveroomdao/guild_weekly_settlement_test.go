package liveroomdao

import "testing"

func TestResetGuildWeeklyAnchorSettlementKeepsBlockedGuildForRetry(t *testing.T) {
	ResetGuildWeeklyAnchorSettlement()
	AddGuildWeeklyAnchorSettlement(1, GuildWeeklyAnchorSettlement{SalaryDiamond: 100})
	AddGuildWeeklyAnchorSettlement(2, GuildWeeklyAnchorSettlement{SalaryDiamond: 200})
	MarkGuildWeeklyAnchorSettlementBlocked(1)

	ResetGuildWeeklyAnchorSettlement()
	if got := GetGuildWeeklyAnchorSettlement(1); got.SalaryDiamond != 100 {
		t.Fatalf("blocked guild aggregation was not preserved: %+v", got)
	}
	if got := GetGuildWeeklyAnchorSettlement(2); !got.IsZero() {
		t.Fatalf("completed guild aggregation should be cleared: %+v", got)
	}
	if IsGuildWeeklyAnchorSettlementBlocked(1) {
		t.Fatal("blocked marker must be reset so the next run can retry")
	}
	RemoveGuildWeeklyAnchorSettlement(1)
}
