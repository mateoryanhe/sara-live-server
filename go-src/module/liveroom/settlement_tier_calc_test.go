package liveroom

import (
	"testing"

	"xr-game-server/entity/live"
)

func TestResolveAnchorSettlementTierUsesPromotedSocialLevelAndHighestGameThreshold(t *testing.T) {
	socialTiers := []*entity.AnchorSalarySocialShareCfg{
		{Level: 1, SocialTotalDiamondRevenue: 100, AnchorSocialSharePercent: 10, GuildSocialSharePercent: 4},
		{Level: 2, SocialTotalDiamondRevenue: 1000, AnchorSocialSharePercent: 20, GuildSocialSharePercent: 8},
	}
	gameTiers := []*entity.AnchorGameShareCfg{
		{GameTotalGoldRevenue: 500, AnchorGameSharePercent: 12, GuildGameSharePercent: 6},
		{GameTotalGoldRevenue: 50, AnchorGameSharePercent: 5, GuildGameSharePercent: 2},
	}
	got := resolveAnchorSettlementTier(true, 5, 1500, 700, socialTiers, nil, gameTiers)
	if got.AnchorSocialShareDiamond != 300 || got.GuildSocialShareDiamond != 120 {
		t.Fatalf("unexpected social shares: %+v", got)
	}
	if got.AnchorGameShareGold != 84 || got.GuildGameShareGold != 42 {
		t.Fatalf("unexpected game shares: %+v", got)
	}
}

func TestResolveAnchorSettlementTierUsesFirstSocialLevelBelowUpgradeBoundary(t *testing.T) {
	noSalarySocial := []*entity.AnchorNoSalaryShareCfg{
		{Level: 1, SocialTotalDiamondRevenue: 100, AnchorSocialSharePercent: 15, GuildSocialSharePercent: 5},
	}
	got := resolveAnchorSettlementTier(true, 0, 10, 10,
		[]*entity.AnchorSalarySocialShareCfg{{Level: 1, EffectiveLiveDays: 1, SocialTotalDiamondRevenue: 100, AnchorSocialSharePercent: 99, GuildSocialSharePercent: 99}}, noSalarySocial,
		[]*entity.AnchorGameShareCfg{{GameTotalGoldRevenue: 100}},
	)
	if got.AnchorSocialShareDiamond != 1.5 || got.GuildSocialShareDiamond != 0.5 {
		t.Fatalf("turnover below the first upgrade boundary must use level 1: %+v", got)
	}
	if got.AnchorGameShareGold != 0 || got.GuildGameShareGold != 0 {
		t.Fatalf("game turnover below its minimum threshold must still resolve to zero: %+v", got)
	}
}

func TestResolveNoSalarySocialShareUsesHighestMatchedThreshold(t *testing.T) {
	tiers := []*entity.AnchorNoSalaryShareCfg{
		{Level: 1, SocialTotalDiamondRevenue: 100, AnchorSocialSharePercent: 8, GuildSocialSharePercent: 3},
		{Level: 2, SocialTotalDiamondRevenue: 1000, AnchorSocialSharePercent: 15, GuildSocialSharePercent: 5},
	}
	got := resolveAnchorSettlementTier(false, 0, 1500, 0, nil, tiers, nil)
	if got.AnchorSocialSharePercent != 15 || got.GuildSocialSharePercent != 5 {
		t.Fatalf("unexpected no-salary social tier: %+v", got)
	}
	if got.AnchorSocialShareDiamond != 225 || got.GuildSocialShareDiamond != 75 {
		t.Fatalf("unexpected no-salary social shares: %+v", got)
	}
}

func TestResolveNoSalarySocialShareBelowUpgradeBoundaryUsesFirstLevel(t *testing.T) {
	got := resolveAnchorSettlementTier(false, 0, 99, 0, nil,
		[]*entity.AnchorNoSalaryShareCfg{{Level: 1, SocialTotalDiamondRevenue: 100, AnchorSocialSharePercent: 15, GuildSocialSharePercent: 5}},
		nil,
	)
	if got.AnchorSocialShareDiamond != 14.85 || got.GuildSocialShareDiamond != 4.95 {
		t.Fatalf("turnover below the first upgrade boundary must use level 1: %+v", got)
	}
}

func TestSocialShareExactBoundaryPromotesToNextLevel(t *testing.T) {
	salaryTiers := []*entity.AnchorSalarySocialShareCfg{
		{Level: 1, SocialTotalDiamondRevenue: 200, AnchorSocialSharePercent: 10, GuildSocialSharePercent: 4},
		{Level: 2, SocialTotalDiamondRevenue: 400, AnchorSocialSharePercent: 20, GuildSocialSharePercent: 8},
	}
	anchorPercent, guildPercent := matchSalarySocialSharePercent(200, salaryTiers) // effectiveLiveDays=0 且档内天数门槛为 0
	if anchorPercent != 20 || guildPercent != 8 {
		t.Fatalf("salary social turnover at level 1 boundary must promote to level 2: anchor=%v guild=%v", anchorPercent, guildPercent)
	}

	noSalaryTiers := []*entity.AnchorNoSalaryShareCfg{
		{Level: 1, SocialTotalDiamondRevenue: 200, AnchorSocialSharePercent: 12, GuildSocialSharePercent: 3},
		{Level: 2, SocialTotalDiamondRevenue: 400, AnchorSocialSharePercent: 18, GuildSocialSharePercent: 6},
	}
	anchorPercent, guildPercent = matchNoSalarySocialSharePercent(200, noSalaryTiers)
	if anchorPercent != 18 || guildPercent != 6 {
		t.Fatalf("no-salary social turnover at level 1 boundary must promote to level 2: anchor=%v guild=%v", anchorPercent, guildPercent)
	}
}

func TestPlatformAnchorSettlementTierDropsGuildShare(t *testing.T) {
	got := (anchorSettlementTierResult{
		HasSalary:                true,
		AnchorSocialSharePercent: 10,
		GuildSocialSharePercent:  20,
		AnchorGameSharePercent:   30,
		GuildGameSharePercent:    40,
		AnchorSocialShareDiamond: 50,
		GuildSocialShareDiamond:  60,
		AnchorGameShareGold:      70,
		GuildGameShareGold:       80,
	}).withoutGuildShare()

	if got.GuildSocialSharePercent != 0 || got.GuildGameSharePercent != 0 ||
		got.GuildSocialShareDiamond != 0 || got.GuildGameShareGold != 0 {
		t.Fatalf("platform anchor settlement must not retain guild share: %+v", got)
	}
	if got.AnchorSocialSharePercent != 10 || got.AnchorGameSharePercent != 30 ||
		got.AnchorSocialShareDiamond != 50 || got.AnchorGameShareGold != 70 {
		t.Fatalf("platform anchor settlement must retain anchor share: %+v", got)
	}
}

func TestPlatformAnchorSettlementRequiresMoreThanFiveUsd(t *testing.T) {
	tests := []struct {
		name   string
		amount float64
		want   bool
	}{
		{name: "below minimum", amount: 4.9999, want: false},
		{name: "exact minimum", amount: 5, want: false},
		{name: "over minimum", amount: 5.0001, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := platformAnchorSettlementAmountEligible(tt.amount, 5); got != tt.want {
				t.Fatalf("platformAnchorSettlementAmountEligible(%v) = %v, want %v", tt.amount, got, tt.want)
			}
		})
	}
}

func TestAnchorWeeklySettlementCfgReadyRequiresAllRelevantConfig(t *testing.T) {
	complete := &anchorWeeklySettlementCfg{
		salarySocialTiers:   []*entity.AnchorSalarySocialShareCfg{{}},
		noSalarySocialTiers: []*entity.AnchorNoSalaryShareCfg{{}},
		withSalaryGameTiers: []*entity.AnchorGameShareCfg{{}},
		noSalaryGameTiers:   []*entity.AnchorGameShareCfg{{}},
	}
	if !complete.ready(true) || !complete.ready(false) {
		t.Fatal("complete settlement configuration should be ready")
	}

	tests := []struct {
		name      string
		hasSalary bool
		cfg       *anchorWeeklySettlementCfg
	}{
		{name: "salary social config missing", hasSalary: true, cfg: &anchorWeeklySettlementCfg{noSalarySocialTiers: []*entity.AnchorNoSalaryShareCfg{{}}, withSalaryGameTiers: []*entity.AnchorGameShareCfg{{}}}},
		{name: "no-salary social fallback missing", hasSalary: true, cfg: &anchorWeeklySettlementCfg{salarySocialTiers: []*entity.AnchorSalarySocialShareCfg{{}}, withSalaryGameTiers: []*entity.AnchorGameShareCfg{{}}}},
		{name: "salary game config missing", hasSalary: true, cfg: &anchorWeeklySettlementCfg{salarySocialTiers: []*entity.AnchorSalarySocialShareCfg{{}}, noSalarySocialTiers: []*entity.AnchorNoSalaryShareCfg{{}}}},
		{name: "no-salary social config missing", hasSalary: false, cfg: &anchorWeeklySettlementCfg{noSalaryGameTiers: []*entity.AnchorGameShareCfg{{}}}},
		{name: "no-salary game config missing", hasSalary: false, cfg: &anchorWeeklySettlementCfg{noSalarySocialTiers: []*entity.AnchorNoSalaryShareCfg{{}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cfg.ready(tt.hasSalary) {
				t.Fatal("settlement must stop when a required configuration is missing")
			}
		})
	}
}
