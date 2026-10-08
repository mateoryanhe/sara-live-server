package liveroom

import (
	stdmath "math"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	xmath "xr-game-server/core/math"
	"xr-game-server/core/syndb"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/entity/live"
	"xr-game-server/module/liverevenuesharecfg"
	"xr-game-server/module/wallet"
)

// settleOnShelfGuilds 周一0点:结算全部上架工会未结算收益
func settleOnShelfGuilds() {
	guilds := guilddao.ListOnShelfGuilds()
	ctx := gctx.New()
	g.Log().Infof(ctx, "guild weekly settlement start, guilds=%d", len(guilds))
	for _, guild := range guilds {
		if guild == nil || guild.ID == 0 {
			continue
		}
		settleOneGuild(guild)
	}
	g.Log().Infof(ctx, "guild weekly settlement done")
}

func settleOneGuild(guild *entity.LiveGuild) bool {
	if guild == nil || guild.ID == 0 {
		return false
	}
	if guild.GuildType != entity.LiveGuildTypeCoinMerchant {
		return settleOneNormalGuildTiered(guild)
	}
	guildId := guild.ID
	dailyRows := liveroomdao.ListRecentUnsettledDailyGuildEffectiveLives(guildId)
	unsettled := liveroomdao.GetGuildIncomeUnsettled(guildId)
	if unsettled == nil {
		return false
	}
	snap := unsettled.Snapshot()
	if !snap.HasSettlementFlow() {
		liveroomdao.ResetGuildWeeklyAnchorLegacySettlement(guildId)
		if len(dailyRows) > 0 {
			liveroomdao.MarkDailyGuildEffectiveLivesSettled(dailyRows)
		}
		return false
	}

	// 币商工会按工会总流水统一结算。名下主播不参与底薪及个人社交/游戏分佣。
	guildSharePercent := guild.SharePercent
	exchangeCfg := wallet.GetExchangeCfgSnapshot()
	socialShareDiamond, gameShareGold, gameShareDiamond, totalShareDiamond, receivableUsd :=
		calcCoinMerchantGuildSettlement(&snap, guildSharePercent, exchangeCfg)
	row := entity.NewCoinMerchantGuildIncomeSettlementLogWithBreakdown(
		guildId,
		&snap,
		0,
		receivableUsd,
		guildSharePercent,
		&entity.GuildIncomeSettlementBreakdown{
			SettlementRuleType:       entity.GuildIncomeSettlementRuleLegacy,
			GuildSocialShareAmount:   socialShareDiamond,
			GuildGameShareAmountGold: gameShareGold,
		},
	)
	row.SetPayoutConversion(exchangeCfg.GoldToDiamondRate, exchangeCfg.UsdToGoldRate, gameShareDiamond, totalShareDiamond, receivableUsd)
	if !syndb.FlushUntilIdle(settlementPersistTimeout) ||
		!liveroomdao.VerifyGuildIncomeSettlementLogPersisted(row) ||
		!liveroomdao.VerifyGuildPayoutConversionPersisted(row) {
		g.Log().Errorf(gctx.New(), "coin merchant guild settlement retained: log persist failed guildId=%d logId=%d", guildId, row.ID)
		return false
	}

	// 结算单及换算快照确认落库后再扣除快照，保留落库期间的并发新增流水。
	unsettled.ConsumeSnapshot(&snap)
	liveroomdao.ResetGuildWeeklyAnchorLegacySettlement(guildId)
	settled := liveroomdao.GetGuildIncomeSettled(guildId)
	if settled != nil {
		settled.AddAmounts(&snap)
		settled.AddSettlementShareAmount(totalShareDiamond)
		settled.AddSettlementShareAmountUsd(receivableUsd)
		if receivableUsd != 0 {
			settled.AddSettlementReceivableUsd(receivableUsd)
		}
	}
	if totalShareDiamond != 0 {
		if total := liveroomdao.GetGuildIncomeTotal(guildId); total != nil {
			total.AddSettlementShareAmount(totalShareDiamond)
			total.AddSettlementShareAmountUsd(receivableUsd)
		}
	}
	if receivableUsd != 0 {
		if total := liveroomdao.GetGuildIncomeTotal(guildId); total != nil {
			total.AddSettlementReceivableUsd(receivableUsd)
		}
	}
	if len(dailyRows) > 0 {
		liveroomdao.MarkDailyGuildEffectiveLivesSettled(dailyRows)
	}
	return true
}

// calcCoinMerchantGuildSettlement 按“工会累计总流水 * 币商工会比例”计算。
// 社交流水单位为钻石；游戏分佣保留金币快照，再按同一份结算汇率换算为钻石和 USD。
func calcCoinMerchantGuildSettlement(snap *entity.LiveRoomIncomeAmounts, guildSharePercent float64, exchangeCfg wallet.ExchangeCfgSnapshot) (
	socialShareDiamond, gameShareGold, gameShareDiamond, totalShareDiamond, receivableUsd float64,
) {
	if snap == nil {
		return 0, 0, 0, 0, 0
	}
	socialShareDiamond = liverevenuesharecfg.CalcGuildSettlementShareAmount(snap.TotalSocialIncome, guildSharePercent)
	gameShareGold = liverevenuesharecfg.CalcGuildSettlementShareAmount(snap.TotalGameIncome, guildSharePercent)
	gameShareDiamond = wallet.CalcGoldToDiamondWithSnapshot(gameShareGold, exchangeCfg)
	totalShareDiamond = xmath.AddFloat64(socialShareDiamond, gameShareDiamond)
	receivableUsd = wallet.CalcDiamondToUsdWithSnapshot(totalShareDiamond, exchangeCfg)
	receivableUsd = stdmath.Round(receivableUsd*10000) / 10000
	return
}
