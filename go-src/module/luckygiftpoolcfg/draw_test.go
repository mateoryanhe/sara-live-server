package luckygiftpoolcfg

import (
	"testing"

	liveentity "xr-game-server/entity/live"
)

func TestIsTierEligible(t *testing.T) {
	ok := &liveentity.LuckyGiftMultiplierTier{
		Multiplier: 10,
		Weight:     5,
		Type:       liveentity.LuckyGiftMultiplierTypeLow,
	}
	if !isTierEligible(ok, 100, 1000) {
		t.Fatal("expected eligible")
	}
	if isTierEligible(ok, 100, 999) {
		t.Fatal("payout 1000 exceeds pool 999")
	}
	badType := &liveentity.LuckyGiftMultiplierTier{Multiplier: 2, Weight: 1, Type: "x"}
	if isTierEligible(badType, 10, 1000) {
		t.Fatal("invalid type must be excluded")
	}
	zeroWeight := &liveentity.LuckyGiftMultiplierTier{
		Multiplier: 2,
		Weight:     0,
		Type:       liveentity.LuckyGiftMultiplierTypeLow,
	}
	if isTierEligible(zeroWeight, 10, 1000) {
		t.Fatal("zero weight must be excluded")
	}
}
