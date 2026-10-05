package liveroom

import (
	"math"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/syndb"
	"xr-game-server/dao/anchorgamesharecfgdao"
	"xr-game-server/dao/anchornosalarysharecfgdao"
	"xr-game-server/dao/anchorsalarycfgdao"
	"xr-game-server/dao/anchorsalarysocialsharecfgdao"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/entity/live"
	"xr-game-server/module/wallet"
)

type anchorWeeklySettlementCfg struct {
	salaryCfgs          []*entity.AnchorSalaryCfg
	salarySocialTiers   []*entity.AnchorSalarySocialShareCfg
	noSalarySocialTiers []*entity.AnchorNoSalaryShareCfg
	withSalaryGameTiers []*entity.AnchorGameShareCfg
	noSalaryGameTiers   []*entity.AnchorGameShareCfg
}

type tieredAnchorSettlementOutcome uint8

const (
	tieredAnchorSettlementFailed tieredAnchorSettlementOutcome = iota
	tieredAnchorSettlementNoData
	tieredAnchorSettlementCreated
)

const settlementPersistTimeout = 15 * time.Second

func loadAnchorWeeklySettlementCfg() *anchorWeeklySettlementCfg {
	return &anchorWeeklySettlementCfg{
		salaryCfgs:          anchorsalarycfgdao.ListAllOrderBySalaryDesc(),
		salarySocialTiers:   anchorsalarysocialsharecfgdao.ListAllOrderByLevelAsc(),
		noSalarySocialTiers: anchornosalarysharecfgdao.ListAllOrderByLevelAsc(),
		withSalaryGameTiers: anchorgamesharecfgdao.ListAllBySalaryTypeOrderByThresholdDesc(entity.AnchorGameShareSalaryTypeWithSalary),
		noSalaryGameTiers:   anchorgamesharecfgdao.ListAllBySalaryTypeOrderByThresholdDesc(entity.AnchorGameShareSalaryTypeNoSalary),
	}
}

// ready 要求当前主播类型涉及的全部结算配置都已维护。
// 任一配置表为空时终止该主播结算，避免用 0 或默认比例生成不可审核的结算单。
func (c *anchorWeeklySettlementCfg) ready(hasSalary bool) bool {
	if c == nil {
		return false
	}
	if hasSalary {
		return len(c.salaryCfgs) > 0 &&
			len(c.salarySocialTiers) > 0 &&
			len(c.withSalaryGameTiers) > 0
	}
	return len(c.noSalarySocialTiers) > 0 && len(c.noSalaryGameTiers) > 0
}

func normalGuildAnchorTieredReady(room *entity.LiveRoom, cfg *anchorWeeklySettlementCfg, periodEnd time.Time) bool {
	if room == nil || room.GuildId == 0 || cfg == nil {
		return false
	}
	unsettled := liveroomdao.GetLiveRoomIncomeUnsettled(room.ID)
	if unsettled == nil {
		return false
	}
	hasSalary := room.IsSalaryEffective(periodEnd.Add(-time.Nanosecond))
	return cfg.ready(hasSalary)
}

// settleNormalGuildAnchorTiered 结算普通工会主播；配置或落库失败时不扣除未结算流水。
func settleNormalGuildAnchorTiered(room *entity.LiveRoom, cfg *anchorWeeklySettlementCfg, periodEnd time.Time) bool {
	if room == nil || room.GuildId == 0 || cfg == nil {
		return false
	}
	return settleTieredAnchor(room, cfg, periodEnd, false)
}

// settlePlatformAnchorTiered 使用与普通工会主播一致的档位结算平台主播，
// 但只生成主播自己的直接代付单，不向任何工会聚合。
func settlePlatformAnchorTiered(room *entity.LiveRoom, cfg *anchorWeeklySettlementCfg, periodEnd time.Time) bool {
	if room == nil || room.ID == 0 || room.GuildId != 0 || cfg == nil {
		return false
	}
	return settleTieredAnchor(room, cfg, periodEnd, true)
}

func settleTieredAnchor(room *entity.LiveRoom, cfg *anchorWeeklySettlementCfg, periodEnd time.Time, directPayout bool) bool {
	return settleTieredAnchorWithOutcome(room, cfg, periodEnd, directPayout) != tieredAnchorSettlementFailed
}

func settlePlatformAnchorTieredWithOutcome(room *entity.LiveRoom, cfg *anchorWeeklySettlementCfg, periodEnd time.Time) tieredAnchorSettlementOutcome {
	if room == nil || room.ID == 0 || room.GuildId != 0 || cfg == nil {
		return tieredAnchorSettlementFailed
	}
	return settleTieredAnchorWithOutcome(room, cfg, periodEnd, true)
}

func settleTieredAnchorWithOutcome(room *entity.LiveRoom, cfg *anchorWeeklySettlementCfg, periodEnd time.Time, directPayout bool) tieredAnchorSettlementOutcome {
	dailyRows := liveroomdao.ListRecentUnsettledDailyEffectiveLives(room.ID)
	unsettled := liveroomdao.GetLiveRoomIncomeUnsettled(room.ID)
	if unsettled == nil {
		return tieredAnchorSettlementFailed
	}
	hasSalary := room.IsSalaryEffective(periodEnd.Add(-time.Nanosecond))
	snap := unsettled.Snapshot()
	if !cfg.ready(hasSalary) {
		g.Log().Warningf(gctx.New(), "anchor weekly settlement skipped: config missing roomId=%d guildId=%d hasSalary=%v", room.ID, room.GuildId, hasSalary)
		return tieredAnchorSettlementFailed
	}

	salary := float64(0)
	if hasSalary {
		salary = matchAnchorSalaryAmount(countAnchorWeeklyWorkDays(dailyRows), dailyRows, cfg.salaryCfgs)
	}
	if len(dailyRows) == 0 && !snap.HasSettlementFlow() && salary == 0 {
		return tieredAnchorSettlementNoData
	}
	gameTiers := cfg.noSalaryGameTiers
	if hasSalary {
		gameTiers = cfg.withSalaryGameTiers
	}
	result := resolveAnchorSettlementTier(hasSalary, snap.TotalSocialIncome, snap.TotalGameIncome, cfg.salarySocialTiers, cfg.noSalarySocialTiers, gameTiers)
	if directPayout {
		result = result.withoutGuildShare()
	}
	breakdown := &entity.AnchorIncomeSettlementBreakdown{
		SettlementRuleType:        entity.AnchorIncomeSettlementRuleTiered,
		HasSalary:                 result.HasSalary,
		AnchorSocialSharePercent:  result.AnchorSocialSharePercent,
		GuildSocialSharePercent:   result.GuildSocialSharePercent,
		AnchorGameSharePercent:    result.AnchorGameSharePercent,
		GuildGameSharePercent:     result.GuildGameSharePercent,
		AnchorSocialShareAmount:   result.AnchorSocialShareDiamond,
		GuildSocialShareAmount:    result.GuildSocialShareDiamond,
		AnchorGameShareAmountGold: result.AnchorGameShareGold,
		GuildGameShareAmountGold:  result.GuildGameShareGold,
	}
	var (
		exchangeCfg   wallet.ExchangeCfgSnapshot
		gameDiamond   float64
		totalDiamond  float64
		receivableUsd float64
	)
	if directPayout {
		exchangeCfg = wallet.GetExchangeCfgSnapshot()
		baseDiamond := salary + result.AnchorSocialShareDiamond
		gameDiamond, totalDiamond, receivableUsd = wallet.CalcSettlementUsdWithSnapshot(baseDiamond, result.AnchorGameShareGold, exchangeCfg)
		receivableUsd = math.Round(receivableUsd*10000) / 10000
		minimumSettlementUsd := cfgdao.PlatformAnchorMinimumSettlementUsd()
		if !platformAnchorSettlementAmountEligible(receivableUsd, minimumSettlementUsd) {
			g.Log().Infof(gctx.New(), "platform anchor settlement deferred: amount not over minimum roomId=%d amountUsd=%.4f minimumUsd=%.2f",
				room.ID, receivableUsd, minimumSettlementUsd)
			return tieredAnchorSettlementNoData
		}
	}
	var logRow *entity.AnchorIncomeSettlementLog
	if directPayout {
		logRow = entity.NewPlatformAnchorIncomeSettlementLogWithBreakdown(
			room.ID, &snap, salary, result.AnchorSocialShareDiamond, result.AnchorSocialSharePercent, breakdown,
		)
		logRow.SetPayoutConversion(exchangeCfg.GoldToDiamondRate, exchangeCfg.UsdToGoldRate, gameDiamond, totalDiamond, receivableUsd)
	} else {
		logRow = entity.NewAnchorIncomeSettlementLogWithBreakdown(
			room.ID, &snap, salary, result.AnchorSocialShareDiamond, 0, result.AnchorSocialSharePercent, breakdown,
		)
	}
	if !syndb.FlushUntilIdle(settlementPersistTimeout) ||
		!liveroomdao.VerifyAnchorIncomeSettlementLogPersisted(logRow) ||
		(directPayout && !liveroomdao.VerifyAnchorPayoutConversionPersisted(logRow)) {
		g.Log().Errorf(gctx.New(), "anchor weekly settlement retained: log persist failed roomId=%d guildId=%d directPayout=%v logId=%d", room.ID, room.GuildId, directPayout, logRow.ID)
		return tieredAnchorSettlementFailed
	}

	// 只扣除已记录的快照，保留持久化期间产生的新流水。
	unsettled.ConsumeSnapshot(&snap)
	if settled := liveroomdao.GetLiveRoomIncomeSettled(room.ID); settled != nil {
		settled.AddAmounts(&snap)
		settled.AddSettlementSalary(salary)
		settled.AddSettlementShareAmount(result.AnchorSocialShareDiamond)
	}
	if total := liveroomdao.GetLiveRoomIncomeTotal(room.ID); total != nil {
		total.AddSettlementSalary(salary)
		total.AddSettlementShareAmount(result.AnchorSocialShareDiamond)
	}
	if !directPayout {
		liveroomdao.AddGuildWeeklyAnchorSettlement(room.GuildId, liveroomdao.GuildWeeklyAnchorSettlement{
			SalaryDiamond:            salary,
			AnchorSocialShareDiamond: result.AnchorSocialShareDiamond,
			GuildSocialShareDiamond:  result.GuildSocialShareDiamond,
			AnchorGameShareGold:      result.AnchorGameShareGold,
			GuildGameShareGold:       result.GuildGameShareGold,
		})
	}
	if len(dailyRows) > 0 {
		liveroomdao.MarkDailyEffectiveLivesSettled(dailyRows)
	}
	return tieredAnchorSettlementCreated
}

func platformAnchorSettlementAmountEligible(receivableUsd, minimumSettlementUsd float64) bool {
	return receivableUsd > minimumSettlementUsd
}
