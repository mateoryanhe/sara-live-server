package entity

import "testing"

func TestCalcTieredGuildPayoutDiamondSplit(t *testing.T) {
	anchor, guild := CalcTieredGuildPayoutDiamondSplit(10, 20, 5, 1, 0.5, 100)
	if anchor != 130 || guild != 55 {
		t.Fatalf("split = (%v, %v), want (130, 55)", anchor, guild)
	}
}
