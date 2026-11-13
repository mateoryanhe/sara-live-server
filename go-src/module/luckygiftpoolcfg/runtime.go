package luckygiftpoolcfg

import (
	"sync"
	"sync/atomic"
	"time"

	"xr-game-server/dao/cfgdao"
	liveentity "xr-game-server/entity/live"
)

// runtime 奖池运行态；开奖/改配置经队列串行后在持锁下读写。
// cfg 指针用于 SetXxx → syndb；业务读内存字段。
type poolRuntime struct {
	cfg   *liveentity.LuckyGiftPoolCfg
	tiers []*liveentity.LuckyGiftMultiplierTier
}

var (
	runtimeMu     sync.Mutex
	runtime       poolRuntime
	balanceAtomic atomic.Uint64 // float64 bits，供 CMS/只读快速取余额
)

func loadRuntimeFromDB() {
	cfg := cfgdao.LoadLuckyGiftPoolCfg()
	tiers := cfgdao.LoadGlobalLuckyGiftMultiplierTiers()
	if len(tiers) == 0 {
		tiers = defaultGlobalTierRows()
	}

	runtimeMu.Lock()
	defer runtimeMu.Unlock()
	if cfg == nil {
		runtime.cfg = nil
	} else {
		runtime.cfg = cfg
	}
	runtime.tiers = copyTierRows(tiers)
	publishBalance(runtimePoolBalanceLocked())
	storePublicSnapshotLocked()
}

func runtimePoolBalanceLocked() float64 {
	if runtime.cfg == nil {
		return liveentity.DefaultLuckyGiftPoolBalance
	}
	return normalizePoolBalance(runtime.cfg.PoolBalance)
}

func runtimeContributeLocked() float64 {
	if runtime.cfg == nil {
		return liveentity.DefaultLuckyGiftContributePercent
	}
	return normalizePercent(runtime.cfg.ContributePercent)
}

func runtimeWinProbLocked() float64 {
	if runtime.cfg == nil {
		return liveentity.DefaultLuckyGiftWinProbability
	}
	return normalizePercent(runtime.cfg.WinProbability)
}

func storePublicSnapshotLocked() {
	balance := runtimePoolBalanceLocked()
	luckyGiftPoolCache.Store(&luckyGiftPoolSnapshot{
		poolBalance:       balance,
		contributePercent: runtimeContributeLocked(),
		winProbability:    runtimeWinProbLocked(),
		tiers:             copyTierRows(runtime.tiers),
	})
}

func publishBalance(v float64) {
	balanceAtomic.Store(float64bits(v))
}

func LivePoolBalance() float64 {
	return float64frombits(balanceAtomic.Load())
}

func copyTierRows(src []*liveentity.LuckyGiftMultiplierTier) []*liveentity.LuckyGiftMultiplierTier {
	if len(src) == 0 {
		return defaultGlobalTierRows()
	}
	out := make([]*liveentity.LuckyGiftMultiplierTier, 0, len(src))
	for _, row := range src {
		if row == nil {
			continue
		}
		cp := *row
		out = append(out, &cp)
	}
	if len(out) == 0 {
		return defaultGlobalTierRows()
	}
	return out
}

func ensureRuntimeCfgLocked() *liveentity.LuckyGiftPoolCfg {
	if runtime.cfg != nil {
		return runtime.cfg
	}
	cfg := liveentity.NewLuckyGiftPoolCfg(
		liveentity.DefaultLuckyGiftPoolBalance,
		liveentity.DefaultLuckyGiftContributePercent,
		liveentity.DefaultLuckyGiftWinProbability,
	)
	runtime.cfg = cfg
	return cfg
}

func setRuntimePoolBalanceLocked(v float64) {
	cfg := ensureRuntimeCfgLocked()
	cfg.SetPoolBalance(v)
	publishBalance(cfg.PoolBalance)
}

func applyConfigLocked(cfgID uint64, poolBalance, contributePercent, winProbability float64, tiers []*liveentity.LuckyGiftMultiplierTier) (uint64, error) {
	_ = cfgID // 客户端 id 仅参考；开奖可能已 New 出行，且前端 Number 易丢 snowflake 精度
	poolBalance = normalizePoolBalance(poolBalance)
	contributePercent = normalizePercent(contributePercent)
	winProbability = normalizePercent(winProbability)

	// 以运行态为准，避免并发/精度导致 id 对不上直接 InvalidParam
	if runtime.cfg != nil {
		runtime.cfg.SetPoolBalance(poolBalance)
		runtime.cfg.SetContributePercent(contributePercent)
		runtime.cfg.SetWinProbability(winProbability)
	} else {
		runtime.cfg = liveentity.NewLuckyGiftPoolCfg(poolBalance, contributePercent, winProbability)
	}

	if err := cfgdao.ReplaceGlobalLuckyGiftMultiplierTiers(tiers); err != nil {
		return 0, err
	}
	if err := flushLuckyGiftPoolSyndb(); err != nil {
		return 0, err
	}
	reloaded := cfgdao.LoadGlobalLuckyGiftMultiplierTiers()
	if len(reloaded) == 0 {
		runtime.tiers = copyTierRows(tiers)
	} else {
		runtime.tiers = copyTierRows(reloaded)
	}
	publishBalance(runtime.cfg.PoolBalance)
	storePublicSnapshotLocked()
	return runtime.cfg.ID, nil
}

func snapshotForCMS() (id uint64, poolBalance, contributePercent, winProbability, totalConsume, totalPayout float64, tiers []*liveentity.LuckyGiftMultiplierTier, createdAt, updatedAt time.Time) {
	runtimeMu.Lock()
	defer runtimeMu.Unlock()
	poolBalance = runtimePoolBalanceLocked()
	contributePercent = runtimeContributeLocked()
	winProbability = runtimeWinProbLocked()
	tiers = copyTierRows(runtime.tiers)
	if runtime.cfg != nil {
		id = runtime.cfg.ID
		createdAt = runtime.cfg.CreatedAt
		updatedAt = runtime.cfg.UpdatedAt
		totalConsume = runtime.cfg.TotalConsumeAmount
		totalPayout = runtime.cfg.TotalPayoutAmount
	}
	return
}
