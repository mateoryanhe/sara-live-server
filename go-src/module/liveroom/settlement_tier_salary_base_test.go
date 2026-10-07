package liveroom

import (
	"testing"

	"xr-game-server/entity/live"
)

func TestResolveBaseSalaryDiamondUsesMatchedTier(t *testing.T) {
	tiers := []*entity.AnchorSalarySocialShareCfg{
		{Level: 1, EffectiveLiveDays: 3, LiveBaseSalaryDiamond: 100, SocialTotalDiamondRevenue: 1000, AnchorSocialSharePercent: 10, GuildSocialSharePercent: 4},
		{Level: 2, EffectiveLiveDays: 5, LiveBaseSalaryDiamond: 300, SocialTotalDiamondRevenue: 5000, AnchorSocialSharePercent: 20, GuildSocialSharePercent: 8},
	}
	if got := resolveBaseSalaryDiamond(true, 3, 500, tiers); got != 100 {
		t.Fatalf("expected level 1 base salary 100, got %v", got)
	}
	if got := resolveBaseSalaryDiamond(true, 5, 1500, tiers); got != 300 {
		t.Fatalf("expected level 2 base salary 300, got %v", got)
	}
	if got := resolveBaseSalaryDiamond(true, 2, 1500, tiers); got != 0 {
		t.Fatalf("insufficient effective days must not pay base salary, got %v", got)
	}
}

func TestHasSalaryAnchorFallsBackToNoSalarySocialWhenTrackNotQualified(t *testing.T) {
	salaryTiers := []*entity.AnchorSalarySocialShareCfg{
		{Level: 1, EffectiveLiveDays: 5, LiveBaseSalaryDiamond: 100, SocialTotalDiamondRevenue: 0, AnchorSocialSharePercent: 99, GuildSocialSharePercent: 99},
	}
	noSalaryTiers := []*entity.AnchorNoSalaryShareCfg{
		{Level: 1, SocialTotalDiamondRevenue: 100, AnchorSocialSharePercent: 12, GuildSocialSharePercent: 3},
	}
	got := resolveAnchorSettlementTier(true, 1, 500, 0, salaryTiers, noSalaryTiers, nil)
	if got.AnchorSocialSharePercent != 12 || got.GuildSocialSharePercent != 3 {
		t.Fatalf("expected no-salary social tier percents, got %+v", got)
	}
	if got.AnchorSocialShareDiamond != 60 {
		t.Fatalf("expected no-salary social share amount 60, got %v", got.AnchorSocialShareDiamond)
	}
}
