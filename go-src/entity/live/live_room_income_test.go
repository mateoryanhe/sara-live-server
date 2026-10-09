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

func TestLiveRoomIncomeAmountsCoinMerchantSettlementUsesGiftOnly(t *testing.T) {
	if (&LiveRoomIncomeAmounts{TotalSocialIncome: 100, TotalGameIncome: 1}).HasCoinMerchantGuildSettlementFlow() {
		t.Fatal("social/game without gift must not trigger coin merchant guild settlement")
	}
	if !(&LiveRoomIncomeAmounts{TotalGiftIncome: 1}).HasCoinMerchantGuildSettlementFlow() {
		t.Fatal("gift flow must trigger coin merchant guild settlement")
	}
	snap := GiftFlowSettlementSnapshot(&LiveRoomIncomeAmounts{TotalGiftIncome: 12, TotalGameIncome: 99})
	if snap.TotalGiftIncome != 12 || snap.TotalGameIncome != 0 || snap.TotalSocialIncome != 12 {
		t.Fatalf("gift snapshot = %+v", snap)
	}
}
