package liveroom

import (
	"math"
	"testing"

	"xr-game-server/entity/live"
	"xr-game-server/module/wallet"
)

func TestCalcCoinMerchantGuildSettlementUsesGiftFlowOnly(t *testing.T) {
	snap := &entity.LiveRoomIncomeAmounts{
		TotalIncome:       999999,
		TotalSocialIncome: 500,
		TotalGiftIncome:   120,
		TotalGameIncome:   2.5,
	}
	exchangeCfg := wallet.ExchangeCfgSnapshot{GoldToDiamondRate: 100, UsdToGoldRate: 100}
	socialShare, gameShareGold, gameShareDiamond, totalShareDiamond, receivableUsd :=
		calcCoinMerchantGuildSettlement(snap, 20, exchangeCfg)

	if socialShare != 24 {
		t.Fatalf("gift share diamond = %v, want 24", socialShare)
	}
	if gameShareGold != 0 || gameShareDiamond != 0 {
		t.Fatalf("game share must be zero, got gold=%v diamond=%v", gameShareGold, gameShareDiamond)
	}
	if totalShareDiamond != 24 {
		t.Fatalf("total share diamond = %v, want 24", totalShareDiamond)
	}
	if math.Abs(receivableUsd-0.0024) > 0.0000001 {
		t.Fatalf("receivable usd = %v, want 0.0024", receivableUsd)
	}
}
