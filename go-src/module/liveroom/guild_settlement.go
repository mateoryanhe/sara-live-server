package liveroom

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/core/math"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/entity/live"
	"xr-game-server/gameevent"
	"xr-game-server/module/liverevenuesharecfg"
	"xr-game-server/module/wallet"
)

func initGuildSettlement() {
	event.Sub(gameevent.WeekEvent, onWeekGuildSettlement)
}

func onWeekGuildSettlement(_ any) {
	settleOnShelfGuilds()
}

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

// resolveGuildSettlementSharePercent 币商工会用自身 sharePercent;普通工会用全局工会分佣
func resolveGuildSettlementSharePercent(guild *entity.LiveGuild) float64 {
	if guild != nil && guild.GuildType == entity.LiveGuildTypeCoinMerchant {
		return guild.SharePercent
	}
	return liverevenuesharecfg.ResolveGuildSharePercent()
}

func settleOneGuild(guild *entity.LiveGuild) {
	guildId := guild.ID
	dailyRows := liveroomdao.ListRecentUnsettledDailyGuildEffectiveLives(guildId)
	unsettled := liveroomdao.GetGuildIncomeUnsettled(guildId)
	if unsettled == nil {
		return
	}
	hasDaily := len(dailyRows) > 0
	hasUnsettled := !unsettled.IsZero()
	guildSharePercent := resolveGuildSettlementSharePercent(guild)
	if !hasDaily && !hasUnsettled {
		weeklySalary := liveroomdao.TakeGuildWeeklyAnchorSalary(guildId)
		weeklyAnchorShareAmountUsd := liveroomdao.TakeGuildWeeklyAnchorShareAmountUsd(guildId)
		if weeklySalary == 0 && weeklyAnchorShareAmountUsd == 0 {
			return
		}
		weeklySalaryUsd := wallet.CalcDiamondToUsd(weeklySalary)
		receivableUsd := math.AddFloat64(weeklyAnchorShareAmountUsd, weeklySalaryUsd)
		_ = entity.NewGuildIncomeSettlementLog(guildId, &entity.LiveRoomIncomeAmounts{}, weeklySalary, 0, 0, receivableUsd, guildSharePercent)
		return
	}

	snap := unsettled.SnapshotAndClear()
	weeklySalary := liveroomdao.TakeGuildWeeklyAnchorSalary(guildId)
	weeklyAnchorShareAmountUsd := liveroomdao.TakeGuildWeeklyAnchorShareAmountUsd(guildId)
	shareAmount := liverevenuesharecfg.CalcGuildSettlementShareAmount(snap.TotalIncome, guildSharePercent)
	shareAmountUsd := wallet.CalcDiamondToUsd(shareAmount)
	weeklySalaryUsd := wallet.CalcDiamondToUsd(weeklySalary)
	receivableUsd := math.AddFloat64(shareAmountUsd, math.AddFloat64(weeklyAnchorShareAmountUsd, weeklySalaryUsd))
	settled := liveroomdao.GetGuildIncomeSettled(guildId)
	if settled != nil {
		settled.AddAmounts(&snap)
		settled.AddSettlementShareAmount(shareAmount)
		settled.AddSettlementShareAmountUsd(shareAmountUsd)
		if shareAmountUsd != 0 {
			settled.AddSettlementReceivableUsd(shareAmountUsd)
		}
	}
	if shareAmount != 0 {
		if total := liveroomdao.GetGuildIncomeTotal(guildId); total != nil {
			total.AddSettlementShareAmount(shareAmount)
			total.AddSettlementShareAmountUsd(shareAmountUsd)
		}
	}
	if shareAmountUsd != 0 {
		if total := liveroomdao.GetGuildIncomeTotal(guildId); total != nil {
			total.AddSettlementReceivableUsd(shareAmountUsd)
		}
	}
	if hasDaily {
		liveroomdao.MarkDailyGuildEffectiveLivesSettled(dailyRows)
	}
	_ = entity.NewGuildIncomeSettlementLog(guildId, &snap, weeklySalary, shareAmount, shareAmountUsd, receivableUsd, guildSharePercent)
}
