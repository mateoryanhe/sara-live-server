package luckygiftpoolcfg

import (
	"math/rand"

	liveentity "xr-game-server/entity/live"
)

// DrawResult 单次抽奖结果(入池+判奖+扣池后)
type DrawResult struct {
	WinProbability float64
	Multiplier     uint32
	Payout         float64
	TierType       string
}

func processBatchDrawLocked(unitPrice float64, count int) []DrawResult {
	if count <= 0 {
		return nil
	}
	runtimeMu.Lock()
	defer runtimeMu.Unlock()

	out := make([]DrawResult, 0, count)
	balance := runtimePoolBalanceLocked()
	startBalance := balance
	contributePct := runtimeContributeLocked()
	winProb := runtimeWinProbLocked()
	forceNoWin := isSettleNoWin()
	var pendingCredit float64

	var payoutSum float64
	for i := 0; i < count; i++ {
		d, next, credit := drawOnceOnBalance(unitPrice, balance, contributePct, winProb, runtime.tiers, forceNoWin)
		balance = next
		pendingCredit += credit
		payoutSum = payoutSum + d.Payout
		out = append(out, d)
	}
	accumulateLuckyGiftTotalsLocked(unitPrice, count, payoutSum)
	if forceNoWin {
		// 进池金额进临时表，不改 pool_balance
		if pendingCredit > 0 {
			enqueuePendingPoolCredit(pendingCredit)
		}
		return out
	}
	if balance != startBalance {
		setRuntimePoolBalanceLocked(balance)
		storePublicSnapshotLocked()
	}
	return out
}

// accumulateLuckyGiftTotalsLocked 开奖后 syndb 累计消费/返奖(须已持 runtimeMu)
func accumulateLuckyGiftTotalsLocked(unitPrice float64, count int, payoutSum float64) {
	if count <= 0 {
		return
	}
	cfg := ensureRuntimeCfgLocked()
	if unitPrice > 0 {
		cfg.AddTotalConsumeAmount(unitPrice * float64(count))
	}
	if payoutSum > 0 {
		cfg.AddTotalPayoutAmount(payoutSum)
	}
}

// drawOnceOnBalance 单次抽奖。
// forceNoWin 时强制不中奖，进池金额经 credit 返回由上层写入 pending。
// 正常模式 credit 恒为 0，进池直接加在 balance 上。
func drawOnceOnBalance(
	unitPrice, balance, contributePct, winProb float64,
	tiers []*liveentity.LuckyGiftMultiplierTier,
	forceNoWin bool,
) (DrawResult, float64, float64) {
	res := DrawResult{WinProbability: winProb}
	var credit float64

	if unitPrice > 0 && contributePct > 0 {
		contribute := unitPrice * contributePct / 100
		if contribute > 0 {
			if forceNoWin {
				credit = contribute
			} else {
				balance += contribute
			}
		}
	}

	if forceNoWin || !rollWin(winProb) {
		return res, balance, credit
	}

	tier := pickValidWeightedTier(unitPrice, balance, tiers)
	if tier == nil {
		return res, balance, credit
	}
	multiplier := tier.Multiplier
	if multiplier == 0 {
		return res, balance, credit
	}
	payout := unitPrice * float64(multiplier)
	if payout <= 0 || payout > balance {
		return res, balance, credit
	}

	balance -= payout
	res.Multiplier = multiplier
	res.Payout = payout
	res.TierType = tier.Type
	if !liveentity.IsValidLuckyGiftMultiplierType(res.TierType) {
		res.TierType = liveentity.LuckyGiftMultiplierTypeLow
	}
	return res, balance, credit
}

func rollWin(winProbability float64) bool {
	p := normalizePercent(winProbability)
	if p <= 0 {
		return false
	}
	if p >= 100 {
		return true
	}
	return rand.Float64()*100 < p
}

// pickValidWeightedTier 仅「倍率>0、权重>0、类型合法、返奖不超过奖池」的挡位参与权重随机
func pickValidWeightedTier(unitPrice, poolBalance float64, tiers []*liveentity.LuckyGiftMultiplierTier) *liveentity.LuckyGiftMultiplierTier {
	type candidate struct {
		tier   *liveentity.LuckyGiftMultiplierTier
		weight uint32
	}
	cands := make([]candidate, 0, len(tiers))
	var totalWeight uint64
	for _, tier := range tiers {
		if !isTierEligible(tier, unitPrice, poolBalance) {
			continue
		}
		cands = append(cands, candidate{tier: tier, weight: tier.Weight})
		totalWeight += uint64(tier.Weight)
	}
	if len(cands) == 0 || totalWeight == 0 {
		return nil
	}
	roll := rand.Uint64()%totalWeight + 1
	var acc uint64
	for _, c := range cands {
		acc += uint64(c.weight)
		if roll <= acc {
			return c.tier
		}
	}
	return cands[len(cands)-1].tier
}

func isTierEligible(tier *liveentity.LuckyGiftMultiplierTier, unitPrice, poolBalance float64) bool {
	if tier == nil {
		return false
	}
	if tier.Multiplier == 0 || tier.Weight == 0 {
		return false
	}
	if !liveentity.IsValidLuckyGiftMultiplierType(tier.Type) {
		return false
	}
	if unitPrice <= 0 || poolBalance <= 0 {
		return false
	}
	payout := unitPrice * float64(tier.Multiplier)
	return payout > 0 && payout <= poolBalance
}

func processApplyConfigLocked(p *applyConfigPayload) (uint64, error) {
	if p == nil {
		return 0, errInvalidParam()
	}
	rows := make([]*liveentity.LuckyGiftMultiplierTier, 0, requiredGlobalTierCount)
	if len(p.tiers) != requiredGlobalTierCount {
		return 0, errInvalidParam()
	}
	for i, item := range p.tiers {
		if item == nil || item.Multiplier == 0 || item.Weight == 0 {
			return 0, errInvalidParam()
		}
		if !liveentity.IsValidLuckyGiftMultiplierType(item.Type) {
			return 0, errInvalidParam()
		}
		rows = append(rows, &liveentity.LuckyGiftMultiplierTier{
			GiftID:     liveentity.GlobalLuckyGiftMultiplierGiftID,
			Multiplier: item.Multiplier,
			Weight:     item.Weight,
			Type:       item.Type,
			Sort:       i,
		})
	}

	runtimeMu.Lock()
	defer runtimeMu.Unlock()
	return applyConfigLocked(p.cfgID, p.poolBalance, p.contributePercent, p.winProbability, rows)
}
