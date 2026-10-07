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

// withoutGuildShare 平台主播没有所属工会，结算单不得保留工会分佣比例或金额。
func (r anchorSettlementTierResult) withoutGuildShare() anchorSettlementTierResult {
	r.GuildSocialSharePercent = 0
	r.GuildGameSharePercent = 0
	r.GuildSocialShareDiamond = 0
	r.GuildGameShareGold = 0
	return r
}

func resolveAnchorSettlementTier(
	hasSalary bool,
	effectiveLiveDays uint64,
	socialDiamond, gameGold float64,
	salarySocialTiers []*entity.AnchorSalarySocialShareCfg,
	noSalarySocialTiers []*entity.AnchorNoSalaryShareCfg,
	gameTiers []*entity.AnchorGameShareCfg,
) anchorSettlementTierResult {
	ret := anchorSettlementTierResult{HasSalary: hasSalary}
	if hasSalary && qualifiesForSalarySocialTrack(effectiveLiveDays, socialDiamond, salarySocialTiers) {
		matched := matchSalarySocialShareTier(effectiveLiveDays, socialDiamond, salarySocialTiers)
		if matched != nil {
			ret.AnchorSocialSharePercent = matched.AnchorSocialSharePercent
			ret.GuildSocialSharePercent = matched.GuildSocialSharePercent
		}
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

// resolveBaseSalaryDiamond 有底薪主播按社交流水分佣档的有效天+流水门槛发放直播底薪(钻石)。
func resolveBaseSalaryDiamond(
	hasSalary bool,
	effectiveLiveDays uint64,
	socialDiamond float64,
	salarySocialTiers []*entity.AnchorSalarySocialShareCfg,
) float64 {
	if !hasSalary || !qualifiesForSalarySocialTrack(effectiveLiveDays, socialDiamond, salarySocialTiers) {
		return 0
	}
	matched := matchSalarySocialShareTier(effectiveLiveDays, socialDiamond, salarySocialTiers)
	if matched == nil || effectiveLiveDays < matched.EffectiveLiveDays {
		return 0
	}
	return matched.LiveBaseSalaryDiamond
}

// qualifiesForSalarySocialTrack 是否启用有底薪社交流水分佣表(否则社交提成走无底薪表)。
// 进入条件只看 1 档有效直播天数；社交流水边界仅用于档内升级。
func qualifiesForSalarySocialTrack(effectiveLiveDays uint64, _ float64, tiers []*entity.AnchorSalarySocialShareCfg) bool {
	validTiers := filterSalarySocialTiers(tiers)
	if len(validTiers) == 0 {
		return false
	}
	return effectiveLiveDays >= validTiers[0].EffectiveLiveDays
}

func filterSalarySocialTiers(tiers []*entity.AnchorSalarySocialShareCfg) []*entity.AnchorSalarySocialShareCfg {
	validTiers := make([]*entity.AnchorSalarySocialShareCfg, 0, len(tiers))
	for _, tier := range tiers {
		if tier != nil {
			validTiers = append(validTiers, tier)
		}
	}
	return validTiers
}

// matchSalarySocialShareTier 有底薪档:流水与有效天同时达到当前档边界才升到下一档。
func matchSalarySocialShareTier(effectiveLiveDays uint64, socialDiamond float64, tiers []*entity.AnchorSalarySocialShareCfg) *entity.AnchorSalarySocialShareCfg {
	validTiers := filterSalarySocialTiers(tiers)
	if len(validTiers) == 0 {
		return nil
	}
	matched := validTiers[0]
	for i := 0; i < len(validTiers)-1; i++ {
		cur := validTiers[i]
		if socialDiamond < cur.SocialTotalDiamondRevenue || effectiveLiveDays < cur.EffectiveLiveDays {
			break
		}
		matched = validTiers[i+1]
	}
	return matched
}

func matchSalarySocialSharePercent(totalDiamond float64, tiers []*entity.AnchorSalarySocialShareCfg) (float64, float64) {
	matched := matchSalarySocialShareTier(0, totalDiamond, tiers)
	if matched == nil {
		return 0, 0
	}
	return matched.AnchorSocialSharePercent, matched.GuildSocialSharePercent
}

func matchNoSalarySocialSharePercent(totalDiamond float64, tiers []*entity.AnchorNoSalaryShareCfg) (float64, float64) {
	validTiers := make([]*entity.AnchorNoSalaryShareCfg, 0, len(tiers))
	for _, tier := range tiers {
		if tier != nil {
			validTiers = append(validTiers, tier)
		}
	}
	if len(validTiers) == 0 {
		return 0, 0
	}
	matched := validTiers[0]
	for i := 0; i < len(validTiers)-1; i++ {
		if totalDiamond < validTiers[i].SocialTotalDiamondRevenue {
			break
		}
		matched = validTiers[i+1]
	}
	return matched.AnchorSocialSharePercent, matched.GuildSocialSharePercent
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
