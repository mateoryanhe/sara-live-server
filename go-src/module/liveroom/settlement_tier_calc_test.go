package liveroom

import (
	"testing"

	"xr-game-server/entity/live"
)

func TestResolveAnchorSettlementTierUsesHighestMatchedThreshold(t *testing.T) {
	socialTiers := []*entity.AnchorSalarySocialShareCfg{
		{SocialTotalDiamondRevenue: 1000, AnchorSocialSharePercent: 20, GuildSocialSharePercent: 8},
		{SocialTotalDiamondRevenue: 100, AnchorSocialSharePercent: 10, GuildSocialSharePercent: 4},
	}
	gameTiers := []*entity.AnchorGameShareCfg{
		{GameTotalGoldRevenue: 500, AnchorGameSharePercent: 12, GuildGameSharePercent: 6},
		{GameTotalGoldRevenue: 50, AnchorGameSharePercent: 5, GuildGameSharePercent: 2},
	}
	got := resolveAnchorSettlementTier(true, 1500, 700, socialTiers, nil, gameTiers)
	if got.AnchorSocialShareDiamond != 300 || got.GuildSocialShareDiamond != 120 {
		t.Fatalf("unexpected social shares: %+v", got)
	}
	if got.AnchorGameShareGold != 84 || got.GuildGameShareGold != 42 {
		t.Fatalf("unexpected game shares: %+v", got)
	}
}

func TestResolveAnchorSettlementTierBelowMinimumIsValidZero(t *testing.T) {
	got := resolveAnchorSettlementTier(true, 10, 10,
		[]*entity.AnchorSalarySocialShareCfg{{SocialTotalDiamondRevenue: 100}}, nil,
		[]*entity.AnchorGameShareCfg{{GameTotalGoldRevenue: 100}},
	)
	if got.AnchorSocialShareDiamond != 0 || got.GuildSocialShareDiamond != 0 || got.AnchorGameShareGold != 0 || got.GuildGameShareGold != 0 {
		t.Fatalf("below-minimum turnover must resolve to zero shares: %+v", got)
	}
}

func TestResolveNoSalarySocialShareUsesHighestMatchedThreshold(t *testing.T) {
	tiers := []*entity.AnchorNoSalaryShareCfg{
		{SocialTotalDiamondRevenue: 1000, AnchorSocialSharePercent: 15, GuildSocialSharePercent: 5},
		{SocialTotalDiamondRevenue: 100, AnchorSocialSharePercent: 8, GuildSocialSharePercent: 3},
	}
	got := resolveAnchorSettlementTier(false, 1500, 0, nil, tiers, nil)
	if got.AnchorSocialSharePercent != 15 || got.GuildSocialSharePercent != 5 {
		t.Fatalf("unexpected no-salary social tier: %+v", got)
	}
	if got.AnchorSocialShareDiamond != 225 || got.GuildSocialShareDiamond != 75 {
		t.Fatalf("unexpected no-salary social shares: %+v", got)
	}
}

func TestResolveNoSalarySocialShareBelowMinimumIsValidZero(t *testing.T) {
	got := resolveAnchorSettlementTier(false, 99, 0, nil,
		[]*entity.AnchorNoSalaryShareCfg{{SocialTotalDiamondRevenue: 100, AnchorSocialSharePercent: 15, GuildSocialSharePercent: 5}},
		nil,
	)
	if got.AnchorSocialShareDiamond != 0 || got.GuildSocialShareDiamond != 0 {
		t.Fatalf("below-minimum no-salary turnover must resolve to zero shares: %+v", got)
	}
}

func TestAnchorWeeklySettlementCfgReadyOnlyRequiresUsedRevenueConfig(t *testing.T) {
	cfg := &anchorWeeklySettlementCfg{
		salaryCfgs:          []*entity.AnchorSalaryCfg{{}},
		noSalarySocialTiers: []*entity.AnchorNoSalaryShareCfg{{}},
	}
	if !cfg.ready(true, &entity.LiveRoomIncomeAmounts{}) {
		t.Fatal("zero-flow salaried settlement should only require salary config")
	}
	if !cfg.ready(false, &entity.LiveRoomIncomeAmounts{TotalSocialIncome: 1}) {
		t.Fatal("no-salary social flow should accept the social tier configuration")
	}
	if cfg.ready(true, &entity.LiveRoomIncomeAmounts{TotalSocialIncome: 1}) {
		t.Fatal("salaried social flow must require a social tier")
	}
	if cfg.ready(false, &entity.LiveRoomIncomeAmounts{TotalGameIncome: 1}) {
		t.Fatal("no-salary game flow must require a game tier")
	}
}
