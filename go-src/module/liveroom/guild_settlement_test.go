package liveroom

import (
	"math"
	"testing"

	"xr-game-server/entity/live"
	"xr-game-server/module/wallet"
)

func TestCalcCoinMerchantGuildSettlementUsesGuildFlowAndIgnoresDisplayTotal(t *testing.T) {
	snap := &entity.LiveRoomIncomeAmounts{
		TotalIncome:       999999,
		TotalSocialIncome: 120,
		TotalGameIncome:   2.5,
	}
	exchangeCfg := wallet.ExchangeCfgSnapshot{GoldToDiamondRate: 100, UsdToGoldRate: 100}
	socialShare, gameShareGold, gameShareDiamond, totalShareDiamond, receivableUsd :=
		calcCoinMerchantGuildSettlement(snap, 20, exchangeCfg)

	if socialShare != 24 {
		t.Fatalf("social share = %v, want 24", socialShare)
	}
	if gameShareGold != 0.5 {
		t.Fatalf("game share gold = %v, want 0.5", gameShareGold)
	}
	if gameShareDiamond != 50 {
		t.Fatalf("game share diamond = %v, want 50", gameShareDiamond)
	}
	if totalShareDiamond != 74 {
		t.Fatalf("total share diamond = %v, want 74", totalShareDiamond)
	}
	if math.Abs(receivableUsd-0.0074) > 0.0000001 {
		t.Fatalf("receivable usd = %v, want 0.0074", receivableUsd)
	}
}
