package liveroom

import (
	"math"

	"xr-game-server/entity/live"
)

// anchorSettlementTierResult 使用原始单位保存单个主播的档位计算结果。
type anchorSettlementTierResult struct {
	HasSalary                bool
	AnchorSocialSharePercent float64
	GuildSocialSharePercent  float64
	AnchorGameSharePercent   float64
	GuildGameSharePercent    float64
	AnchorSocialShareDiamond float64
	GuildSocialShareDiamond  float64
	AnchorGameShareGold      float64
	GuildGameShareGold       float64
}

func resolveAnchorSettlementTier(
	hasSalary bool,
	socialDiamond, gameGold float64,
	salarySocialTiers []*entity.AnchorSalarySocialShareCfg,
	noSalarySocialTiers []*entity.AnchorNoSalaryShareCfg,
	gameTiers []*entity.AnchorGameShareCfg,
) anchorSettlementTierResult {
	ret := anchorSettlementTierResult{HasSalary: hasSalary}
	if hasSalary {
		ret.AnchorSocialSharePercent, ret.GuildSocialSharePercent = matchSalarySocialSharePercent(socialDiamond, salarySocialTiers)
	} else {
		ret.AnchorSocialSharePercent, ret.GuildSocialSharePercent = matchNoSalarySocialSharePercent(socialDiamond, noSalarySocialTiers)
	}
	ret.AnchorGameSharePercent, ret.GuildGameSharePercent = matchGameSharePercent(gameGold, gameTiers)
	ret.AnchorSocialShareDiamond = calcTierShare(socialDiamond, ret.AnchorSocialSharePercent)
	ret.GuildSocialShareDiamond = calcTierShare(socialDiamond, ret.GuildSocialSharePercent)
	ret.AnchorGameShareGold = calcTierShare(gameGold, ret.AnchorGameSharePercent)
	ret.GuildGameShareGold = calcTierShare(gameGold, ret.GuildGameSharePercent)
	return ret
}

func matchSalarySocialSharePercent(totalDiamond float64, tiers []*entity.AnchorSalarySocialShareCfg) (float64, float64) {
	for _, tier := range tiers {
		if tier != nil && totalDiamond >= tier.SocialTotalDiamondRevenue {
			return tier.AnchorSocialSharePercent, tier.GuildSocialSharePercent
		}
	}
	return 0, 0
}

func matchNoSalarySocialSharePercent(totalDiamond float64, tiers []*entity.AnchorNoSalaryShareCfg) (float64, float64) {
	for _, tier := range tiers {
		if tier != nil && totalDiamond >= tier.SocialTotalDiamondRevenue {
			return tier.AnchorSocialSharePercent, tier.GuildSocialSharePercent
		}
	}
	return 0, 0
}

func matchGameSharePercent(totalGold float64, tiers []*entity.AnchorGameShareCfg) (float64, float64) {
	for _, tier := range tiers {
		if tier != nil && totalGold >= tier.GameTotalGoldRevenue {
			return tier.AnchorGameSharePercent, tier.GuildGameSharePercent
		}
	}
	return 0, 0
}

func calcTierShare(amount, percent float64) float64 {
	if amount <= 0 || percent <= 0 {
		return 0
	}
	return math.Round(amount*percent) / 100
}
