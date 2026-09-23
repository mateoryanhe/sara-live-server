package liveroom

import (
	"math"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/syndb"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/entity/live"
	"xr-game-server/module/wallet"
)

// settleOneNormalGuildTiered 汇总普通工会分项；落库失败时保留工会未结算流水和本轮主播聚合。
func settleOneNormalGuildTiered(guild *entity.LiveGuild) {
	if guild == nil || guild.ID == 0 {
		return
	}
	if liveroomdao.IsGuildWeeklyAnchorSettlementBlocked(guild.ID) {
		g.Log().Warningf(gctx.New(), "guild weekly settlement retained: anchor settlement incomplete guildId=%d", guild.ID)
		return
	}
	dailyRows := liveroomdao.ListRecentUnsettledDailyGuildEffectiveLives(guild.ID)
	unsettled := liveroomdao.GetGuildIncomeUnsettled(guild.ID)
	if unsettled == nil {
		return
	}
	weekly := liveroomdao.GetGuildWeeklyAnchorSettlement(guild.ID)
	snap := unsettled.Snapshot()
	if len(dailyRows) == 0 && snap.IsZero() && weekly.IsZero() {
		return
	}
	exchangeCfg := wallet.GetExchangeCfgSnapshot()
	gameGold := weekly.AnchorGameShareGold + weekly.GuildGameShareGold
	baseDiamond := weekly.SalaryDiamond + weekly.AnchorSocialShareDiamond + weekly.GuildSocialShareDiamond
	gameDiamond, totalDiamond, receivableUsd := wallet.CalcSettlementUsdWithSnapshot(baseDiamond, gameGold, exchangeCfg)
	receivableUsd = math.Round(receivableUsd*10000) / 10000
	row := entity.NewGuildIncomeSettlementLogWithBreakdown(
		guild.ID,
		&snap,
		weekly.SalaryDiamond,
		weekly.GuildSocialShareDiamond,
		0,
		receivableUsd,
		0,
		&entity.GuildIncomeSettlementBreakdown{
			SettlementRuleType:        entity.GuildIncomeSettlementRuleTiered,
			AnchorSocialShareAmount:   weekly.AnchorSocialShareDiamond,
			GuildSocialShareAmount:    weekly.GuildSocialShareDiamond,
			AnchorGameShareAmountGold: weekly.AnchorGameShareGold,
			GuildGameShareAmountGold:  weekly.GuildGameShareGold,
		},
	)
	row.SetPayoutConversion(exchangeCfg.GoldToDiamondRate, exchangeCfg.UsdToGoldRate, gameDiamond, totalDiamond, receivableUsd)
	if !syndb.FlushUntilIdle(settlementPersistTimeout) ||
		!liveroomdao.VerifyGuildIncomeSettlementLogPersisted(row) ||
		!liveroomdao.VerifyGuildPayoutConversionPersisted(row) {
		liveroomdao.MarkGuildWeeklyAnchorSettlementBlocked(guild.ID)
		g.Log().Errorf(gctx.New(), "guild weekly settlement retained: log persist failed guildId=%d logId=%d", guild.ID, row.ID)
		return
	}

	unsettled.ConsumeSnapshot(&snap)
	liveroomdao.RemoveGuildWeeklyAnchorSettlement(guild.ID)
	if settled := liveroomdao.GetGuildIncomeSettled(guild.ID); settled != nil {
		settled.AddAmounts(&snap)
		settled.AddSettlementSalary(weekly.SalaryDiamond)
		settled.AddSettlementShareAmount(weekly.GuildSocialShareDiamond)
	}
	if total := liveroomdao.GetGuildIncomeTotal(guild.ID); total != nil {
		total.AddSettlementSalary(weekly.SalaryDiamond)
		total.AddSettlementShareAmount(weekly.GuildSocialShareDiamond)
	}
	if len(dailyRows) > 0 {
		liveroomdao.MarkDailyGuildEffectiveLivesSettled(dailyRows)
	}
}
