package luckygiftpoolcfg

import (
	"sync/atomic"

	liveentity "xr-game-server/entity/live"
)

type luckyGiftPoolSnapshot struct {
	poolBalance       float64
	contributePercent float64
	winProbability    float64
	tiers             []*liveentity.LuckyGiftMultiplierTier
}

var luckyGiftPoolCache atomic.Value // *luckyGiftPoolSnapshot

func Init() {
	loadRuntimeFromDB()
	startPoolQueueWorker()
}

func getSnapshot() *luckyGiftPoolSnapshot {
	v := luckyGiftPoolCache.Load()
	if v == nil {
		return &luckyGiftPoolSnapshot{
			poolBalance:       liveentity.DefaultLuckyGiftPoolBalance,
			contributePercent: liveentity.DefaultLuckyGiftContributePercent,
			winProbability:    liveentity.DefaultLuckyGiftWinProbability,
			tiers:             defaultGlobalTierRows(),
		}
	}
	snap, ok := v.(*luckyGiftPoolSnapshot)
	if !ok || snap == nil {
		return &luckyGiftPoolSnapshot{
			poolBalance:       liveentity.DefaultLuckyGiftPoolBalance,
			contributePercent: liveentity.DefaultLuckyGiftContributePercent,
			winProbability:    liveentity.DefaultLuckyGiftWinProbability,
			tiers:             defaultGlobalTierRows(),
		}
	}
	return snap
}

// PoolBalance 当前奖金池余额(只读快照)
func PoolBalance() float64 {
	return LivePoolBalance()
}

// ContributePercent 礼物金额进入奖池比例(%)
func ContributePercent() float64 {
	return getSnapshot().contributePercent
}

// WinProbability 中奖概率(%)
func WinProbability() float64 {
	return getSnapshot().winProbability
}

// GlobalMultiplierTiers 全局默认 5 挡倍率(副本)
func GlobalMultiplierTiers() []*liveentity.LuckyGiftMultiplierTier {
	tiers := getSnapshot().tiers
	if len(tiers) == 0 {
		return defaultGlobalTierRows()
	}
	out := make([]*liveentity.LuckyGiftMultiplierTier, len(tiers))
	copy(out, tiers)
	return out
}

func normalizePoolBalance(v float64) float64 {
	if v < 0 {
		return 0
	}
	return v
}

func normalizePercent(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}
