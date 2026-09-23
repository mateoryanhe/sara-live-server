package incomesettlement

import (
	"testing"

	"xr-game-server/entity/live"
	"xr-game-server/module/wallet"
)

func TestCalculateTieredGuildPayoutKeepsGoldUntilTransfer(t *testing.T) {
	row := &entity.GuildIncomeSettlementLog{
		SettlementSalary:          1000,
		AnchorSocialShareAmount:   2000,
		GuildSocialShareAmount:    500,
		AnchorGameShareAmountGold: 40,
		GuildGameShareAmountGold:  10,
	}
	gameDiamond, totalDiamond, usd, err := calculateTieredGuildPayout(row, wallet.ExchangeCfgSnapshot{
		GoldToDiamondRate: 100,
		UsdToGoldRate:     100,
	})
	if err != nil {
		t.Fatal(err)
	}
	if gameDiamond != 5000 || totalDiamond != 8500 || usd != 0.85 {
		t.Fatalf("unexpected conversion gameDiamond=%v totalDiamond=%v usd=%v", gameDiamond, totalDiamond, usd)
	}
}

func TestCalculateTieredPlatformAnchorPayoutOnlyPaysAnchorShare(t *testing.T) {
	row := &entity.AnchorIncomeSettlementLog{
		SettlementSalary:          1000,
		AnchorSocialShareAmount:   2000,
		GuildSocialShareAmount:    500,
		AnchorGameShareAmountGold: 40,
		GuildGameShareAmountGold:  10,
	}
	gameDiamond, totalDiamond, usd, err := calculateTieredAnchorPayout(row, wallet.ExchangeCfgSnapshot{
		GoldToDiamondRate: 100,
		UsdToGoldRate:     100,
	})
	if err != nil {
		t.Fatal(err)
	}
	if gameDiamond != 4000 || totalDiamond != 7000 || usd != 0.7 {
		t.Fatalf("unexpected platform anchor conversion gameDiamond=%v totalDiamond=%v usd=%v", gameDiamond, totalDiamond, usd)
	}
}

func TestPrepareTieredGuildPayoutKeepsReviewedAmount(t *testing.T) {
	row := &entity.GuildIncomeSettlementLog{
		SettlementRuleType:      entity.GuildIncomeSettlementRuleTiered,
		SettlementReceivableUsd: 12.3456,
		SettlementSalary:        999999,
	}
	if err := prepareTieredGuildPayout(row); err != nil {
		t.Fatal(err)
	}
	if row.SettlementReceivableUsd != 12.3456 {
		t.Fatalf("reviewed amount changed: %v", row.SettlementReceivableUsd)
	}
}

func TestPrepareTieredPlatformAnchorPayoutKeepsReviewedAmount(t *testing.T) {
	row := &entity.AnchorIncomeSettlementLog{
		SettlementRuleType:      entity.AnchorIncomeSettlementRuleTiered,
		SettlementReceivableUsd: 6.789,
		SettlementSalary:        999999,
	}
	if err := prepareTieredAnchorPayout(row); err != nil {
		t.Fatal(err)
	}
	if row.SettlementReceivableUsd != 6.789 {
		t.Fatalf("reviewed amount changed: %v", row.SettlementReceivableUsd)
	}
}
