package entity

import "testing"

func TestLiveRoomIncomeAmountsHasSettlementFlowIgnoresDisplayTotal(t *testing.T) {
	if (&LiveRoomIncomeAmounts{TotalIncome: 999999}).HasSettlementFlow() {
		t.Fatal("display-only total income must not trigger settlement")
	}
	if !(&LiveRoomIncomeAmounts{TotalSocialIncome: 1}).HasSettlementFlow() {
		t.Fatal("social diamond flow must trigger settlement")
	}
	if !(&LiveRoomIncomeAmounts{TotalGameIncome: 1}).HasSettlementFlow() {
		t.Fatal("game gold flow must trigger settlement")
	}
}
