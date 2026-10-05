package wallet

import (
	"math"
	"testing"
)

func TestCalcGoldToDiamondWithSnapshot(t *testing.T) {
	snap := ExchangeCfgSnapshot{GoldToDiamondRate: 100, UsdToGoldRate: 100}
	if got := CalcGoldToDiamondWithSnapshot(2.5, snap); got != 250 {
		t.Fatalf("CalcGoldToDiamondWithSnapshot() = %v, want 250", got)
	}
	if got := CalcGoldToDiamondWithSnapshot(0.29, snap); got != 29 {
		t.Fatalf("decimal conversion = %v, want 29", got)
	}
	if got := CalcGoldToDiamondWithSnapshot(2.5, ExchangeCfgSnapshot{}); got != 0 {
		t.Fatalf("invalid rate conversion = %v, want 0", got)
	}
	if got := CalcGoldToDiamondWithSnapshot(math.NaN(), snap); got != 0 {
		t.Fatalf("NaN conversion = %v, want 0", got)
	}
	if got := CalcGoldToDiamondWithSnapshot(math.Inf(1), snap); got != 0 {
		t.Fatalf("Inf conversion = %v, want 0", got)
	}
}

func TestCalcSettlementUsdWithSnapshot(t *testing.T) {
	snap := ExchangeCfgSnapshot{GoldToDiamondRate: 100, UsdToGoldRate: 100}
	gameDiamond, totalDiamond, usd := CalcSettlementUsdWithSnapshot(50, 2.5, snap)
	if gameDiamond != 250 || totalDiamond != 300 || usd != 0.03 {
		t.Fatalf("conversion = (%v, %v, %v), want (250, 300, 0.03)", gameDiamond, totalDiamond, usd)
	}
}
