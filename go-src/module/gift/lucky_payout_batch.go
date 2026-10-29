package gift

import (
	"sort"

	liveentity "xr-game-server/entity/live"
)

// LuckyGiftTierMultiplierLine 同一挡位内某一倍数的累加
type LuckyGiftTierMultiplierLine struct {
	Multiplier  uint32
	WinCount    int
	TotalPayout float64
}

// LuckyGiftTierWinSummary 某一挡位连抽累加结果
type LuckyGiftTierWinSummary struct {
	WinCount      int
	TotalPayout   float64
	MaxMultiplier uint32
	byMultiplier  map[uint32]*LuckyGiftTierMultiplierLine
}

// LuckyGiftBatchResult 幸运礼物按数量独立判奖后的累加结果
type LuckyGiftBatchResult struct {
	Rtp               float64
	TotalPayout       float64
	MaxWinMultiplier  uint32
	HighestWinTier    string
	Low, Medium, High *LuckyGiftTierWinSummary
}

// ResolveLuckyGiftPayoutBatch 按单价独立判奖 count 次并累加各挡位返奖
func ResolveLuckyGiftPayoutBatch(giftID uint64, unitPrice float64, count int) LuckyGiftBatchResult {
	if count <= 0 {
		return LuckyGiftBatchResult{}
	}
	rtp := GetLuckyGiftRtpByGiftID(giftID)
	draws := make([]LuckyGiftDrawResult, 0, count)
	for i := 0; i < count; i++ {
		draws = append(draws, ResolveLuckyGiftPayout(giftID, unitPrice))
	}
	return AggregateLuckyGiftDraws(rtp, draws)
}

// AggregateLuckyGiftDraws 将多次判奖累加为挡位汇总(供测试与组装)
func AggregateLuckyGiftDraws(rtp float64, draws []LuckyGiftDrawResult) LuckyGiftBatchResult {
	out := LuckyGiftBatchResult{Rtp: rtp}
	if len(draws) == 0 {
		return out
	}
	var low, medium, high *LuckyGiftTierWinSummary
	addWin := func(tier string, payout float64, multiplier uint32) {
		if payout <= 0 || tier == "" {
			return
		}
		out.TotalPayout += payout
		if multiplier > out.MaxWinMultiplier {
			out.MaxWinMultiplier = multiplier
		}
		var slot **LuckyGiftTierWinSummary
		switch tier {
		case liveentity.LuckyGiftMultiplierTypeMedium:
			slot = &medium
		case liveentity.LuckyGiftMultiplierTypeHigh:
			slot = &high
		default:
			slot = &low
		}
		if *slot == nil {
			*slot = &LuckyGiftTierWinSummary{}
		}
		(*slot).WinCount++
		(*slot).TotalPayout += payout
		(*slot).recordMultiplierWin(multiplier, payout)
	}
	for _, d := range draws {
		addWin(d.TierType, d.Payout, d.Multiplier)
	}
	out.Low, out.Medium, out.High = low, medium, high
	out.HighestWinTier = highestLuckyGiftWinTier(out.Low, out.Medium, out.High)
	return out
}

func highestLuckyGiftWinTier(low, medium, high *LuckyGiftTierWinSummary) string {
	if tierHasWins(high) {
		return liveentity.LuckyGiftMultiplierTypeHigh
	}
	if tierHasWins(medium) {
		return liveentity.LuckyGiftMultiplierTypeMedium
	}
	if tierHasWins(low) {
		return liveentity.LuckyGiftMultiplierTypeLow
	}
	return ""
}

func tierHasWins(s *LuckyGiftTierWinSummary) bool {
	return s != nil && s.WinCount > 0 && s.TotalPayout > 0
}

func (s *LuckyGiftTierWinSummary) recordMultiplierWin(multiplier uint32, payout float64) {
	if s == nil || multiplier == 0 || payout <= 0 {
		return
	}
	if multiplier > s.MaxMultiplier {
		s.MaxMultiplier = multiplier
	}
	if s.byMultiplier == nil {
		s.byMultiplier = make(map[uint32]*LuckyGiftTierMultiplierLine)
	}
	line := s.byMultiplier[multiplier]
	if line == nil {
		line = &LuckyGiftTierMultiplierLine{Multiplier: multiplier}
		s.byMultiplier[multiplier] = line
	}
	line.WinCount++
	line.TotalPayout += payout
}

// LuckyGiftTierMultiplierLines 挡位内按倍数升序的明细(供 DTO/推送)
func LuckyGiftTierMultiplierLines(s *LuckyGiftTierWinSummary) []*LuckyGiftTierMultiplierLine {
	if s == nil || len(s.byMultiplier) == 0 {
		return nil
	}
	muls := make([]uint32, 0, len(s.byMultiplier))
	for m := range s.byMultiplier {
		muls = append(muls, m)
	}
	sort.Slice(muls, func(i, j int) bool { return muls[i] < muls[j] })
	out := make([]*LuckyGiftTierMultiplierLine, 0, len(muls))
	for _, m := range muls {
		if line := s.byMultiplier[m]; line != nil && line.WinCount > 0 {
			out = append(out, line)
		}
	}
	return out
}

// PushRouteTierType 飘屏推送路由挡位: 有高推全服, 否则有中推本房, 否则低仅本人
func (b LuckyGiftBatchResult) PushRouteTierType() string {
	if b.HighestWinTier != "" {
		return b.HighestWinTier
	}
	return highestLuckyGiftWinTier(b.Low, b.Medium, b.High)
}
