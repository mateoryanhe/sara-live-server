package gift

import (
	"sort"
	"sync/atomic"

	"xr-game-server/dao/cfgdao"
	liveentity "xr-game-server/entity/live"
)

type luckyGiftMultiplierSnapshot struct {
	byGiftID map[uint64][]*liveentity.LuckyGiftMultiplierTier
}

var (
	luckyGiftMultiplierCache atomic.Value // *luckyGiftMultiplierSnapshot
)

func InitLuckyGiftMultiplier() {
	reloadLuckyGiftMultiplierMemory()
}

func reloadLuckyGiftMultiplierMemory() {
	rows := cfgdao.LoadAllLuckyGiftMultiplierTiers()
	byGiftID := make(map[uint64][]*liveentity.LuckyGiftMultiplierTier)
	for _, row := range rows {
		if row == nil || row.GiftID == 0 {
			continue
		}
		byGiftID[row.GiftID] = append(byGiftID[row.GiftID], row)
	}
	for giftID := range byGiftID {
		sort.Slice(byGiftID[giftID], func(i, j int) bool {
			if byGiftID[giftID][i].Sort != byGiftID[giftID][j].Sort {
				return byGiftID[giftID][i].Sort < byGiftID[giftID][j].Sort
			}
			return byGiftID[giftID][i].ID < byGiftID[giftID][j].ID
		})
	}
	luckyGiftMultiplierCache.Store(&luckyGiftMultiplierSnapshot{byGiftID: byGiftID})
}

func getLuckyGiftMultiplierSnapshot() *luckyGiftMultiplierSnapshot {
	v := luckyGiftMultiplierCache.Load()
	if v == nil {
		return &luckyGiftMultiplierSnapshot{byGiftID: map[uint64][]*liveentity.LuckyGiftMultiplierTier{}}
	}
	snap, ok := v.(*luckyGiftMultiplierSnapshot)
	if !ok || snap == nil {
		return &luckyGiftMultiplierSnapshot{byGiftID: map[uint64][]*liveentity.LuckyGiftMultiplierTier{}}
	}
	return snap
}

// GetLuckyGiftMultiplierTiersByGiftID 某幸运礼物的返奖倍率挡位(未配置时返回内置默认 5 挡)
func GetLuckyGiftMultiplierTiersByGiftID(giftID uint64) []*liveentity.LuckyGiftMultiplierTier {
	if giftID == 0 {
		return defaultLuckyGiftMultiplierTierRows()
	}
	snap := getLuckyGiftMultiplierSnapshot()
	tiers := snap.byGiftID[giftID]
	if len(tiers) == 0 {
		return defaultLuckyGiftMultiplierTierRows()
	}
	out := make([]*liveentity.LuckyGiftMultiplierTier, len(tiers))
	copy(out, tiers)
	return out
}

// GetLuckyGiftRtpByGiftID 某幸运礼物 RTP(百分比); 非幸运或未配置为 0
func GetLuckyGiftRtpByGiftID(giftID uint64) float64 {
	if giftID == 0 {
		return 0
	}
	g := cfgdao.GetGiftById(giftID)
	if g == nil || g.CategoryType != liveentity.LiveGiftCategoryTypeLucky {
		return 0
	}
	return NormalizeLuckyGiftRtp(g.LuckyGiftRtp)
}
