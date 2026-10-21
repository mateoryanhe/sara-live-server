package gift

import (
	"math/rand"

	"xr-game-server/dao/cfgdao"
	liveentity "xr-game-server/entity/live"
)

// LuckyGiftDrawResult 单次幸运礼物判奖结果
type LuckyGiftDrawResult struct {
	Rtp      float64
	Multiplier uint32
	Payout   float64
	// TierType CMS 挡位类型 low/medium/high；未中奖为空
	TierType string
}

// ResolveLuckyGiftPayout 幸运礼物返奖：先 RTP 判中奖，中奖再抽倍数(含 CMS 挡位 type)
func ResolveLuckyGiftPayout(giftID uint64, totalCost float64) LuckyGiftDrawResult {
	rtp := GetLuckyGiftRtpByGiftID(giftID)
	if !RollLuckyGiftWin(rtp) {
		return LuckyGiftDrawResult{Rtp: rtp}
	}
	multiplier, tierType := DrawLuckyGiftMultiplierWithTier(giftID)
	payout := CalcLuckyGiftPayoutDiamonds(totalCost, multiplier)
	return LuckyGiftDrawResult{
		Rtp:        rtp,
		Multiplier: multiplier,
		Payout:     payout,
		TierType:   tierType,
	}
}

// DrawLuckyGiftMultiplier 按礼物挡位权重抽取返奖倍数
func DrawLuckyGiftMultiplier(giftID uint64) uint32 {
	multiplier, _ := DrawLuckyGiftMultiplierWithTier(giftID)
	return multiplier
}

// DrawLuckyGiftMultiplierWithTier 抽取返奖倍数及 CMS 配置的挡位 type(low/medium/high)
func DrawLuckyGiftMultiplierWithTier(giftID uint64) (uint32, string) {
	tiers := GetLuckyGiftMultiplierTiersByGiftID(giftID)
	if len(tiers) == 0 {
		return 1, liveentity.LuckyGiftMultiplierTypeLow
	}
	var totalWeight uint64
	for _, tier := range tiers {
		if tier == nil || tier.Weight == 0 {
			continue
		}
		totalWeight += uint64(tier.Weight)
	}
	if totalWeight == 0 {
		return pickMultiplierAndTierType(tiers[0])
	}
	roll := rand.Uint64()%totalWeight + 1
	var acc uint64
	for _, tier := range tiers {
		if tier == nil || tier.Weight == 0 {
			continue
		}
		acc += uint64(tier.Weight)
		if roll <= acc {
			return pickMultiplierAndTierType(tier)
		}
	}
	return pickMultiplierAndTierType(tiers[len(tiers)-1])
}

func pickMultiplierAndTierType(tier *liveentity.LuckyGiftMultiplierTier) (uint32, string) {
	if tier == nil {
		return 1, liveentity.LuckyGiftMultiplierTypeLow
	}
	multiplier := tier.Multiplier
	if multiplier == 0 {
		multiplier = 1
	}
	tierType := tier.Type
	if !liveentity.IsValidLuckyGiftMultiplierType(tierType) {
		tierType = liveentity.LuckyGiftMultiplierTypeLow
	}
	return multiplier, tierType
}

// CalcLuckyGiftPayoutDiamonds 返奖钻石 = 消耗钻石 * 倍数
func CalcLuckyGiftPayoutDiamonds(totalCost float64, multiplier uint32) float64 {
	if totalCost <= 0 || multiplier == 0 {
		return 0
	}
	return totalCost * float64(multiplier)
}

// CalcLuckyGiftAnchorSettlementEarn 主播礼物流水累计(折算后) = 消耗 * 结算比例%
func CalcLuckyGiftAnchorSettlementEarn(totalCost float64) float64 {
	if totalCost <= 0 {
		return 0
	}
	pct := cfgdao.LuckyGiftSettlementPercent()
	if pct <= 0 {
		return 0
	}
	return totalCost * pct / 100
}

func IsLuckyGiftCategory(categoryType uint8) bool {
	return categoryType == liveentity.LiveGiftCategoryTypeLucky
}
