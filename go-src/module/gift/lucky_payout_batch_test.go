package gift

import (
	"testing"

	liveentity "xr-game-server/entity/live"
)

func TestAggregateLuckyGiftDraws(t *testing.T) {
	batch := AggregateLuckyGiftDraws(10, []LuckyGiftDrawResult{
		{Multiplier: 2, Payout: 20, TierType: liveentity.LuckyGiftMultiplierTypeLow},
		{Multiplier: 3, Payout: 30, TierType: liveentity.LuckyGiftMultiplierTypeLow},
		{Multiplier: 5, Payout: 50, TierType: liveentity.LuckyGiftMultiplierTypeMedium},
	})
	if batch.TotalPayout != 100 {
		t.Fatalf("total=%v", batch.TotalPayout)
	}
	if batch.Low == nil || batch.Low.WinCount != 2 || batch.Low.TotalPayout != 50 {
		t.Fatalf("low=%+v", batch.Low)
	}
	if batch.Low.MaxMultiplier != 3 {
		t.Fatalf("low maxMul=%d", batch.Low.MaxMultiplier)
	}
	lowLines := LuckyGiftTierMultiplierLines(batch.Low)
	if len(lowLines) != 2 || lowLines[0].Multiplier != 2 || lowLines[1].Multiplier != 3 {
		t.Fatalf("low lines=%+v", lowLines)
	}
	if batch.Medium == nil || batch.Medium.WinCount != 1 || batch.Medium.TotalPayout != 50 {
		t.Fatalf("medium=%+v", batch.Medium)
	}
	if batch.PushRouteTierType() != liveentity.LuckyGiftMultiplierTypeMedium {
		t.Fatalf("route=%q", batch.PushRouteTierType())
	}
	if batch.MaxWinMultiplier != 5 {
		t.Fatalf("maxMul=%d", batch.MaxWinMultiplier)
	}
}
