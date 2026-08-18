package liveroom

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/entity/live"
	"xr-game-server/gameevent"
	"xr-game-server/module/liverevenuesharecfg"
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
		settleOneGuild(guild.ID)
	}
	g.Log().Infof(ctx, "guild weekly settlement done")
}

func settleOneGuild(guildId uint64) {
	dailyRows := liveroomdao.ListRecentUnsettledDailyGuildEffectiveLives(guildId)
	unsettled := liveroomdao.GetGuildIncomeUnsettled(guildId)
	if unsettled == nil {
		return
	}
	hasDaily := len(dailyRows) > 0
	hasUnsettled := !unsettled.IsZero()
	if !hasDaily && !hasUnsettled {
		return
	}

	snap := unsettled.SnapshotAndClear()
	guildSharePercent := liverevenuesharecfg.ResolveGuildSharePercent()
	shareAmount := liverevenuesharecfg.CalcGuildSettlementShareAmount(snap.TotalIncome)
	settled := liveroomdao.GetGuildIncomeSettled(guildId)
	if settled != nil {
		settled.AddAmounts(&snap)
		settled.AddSettlementShareAmount(shareAmount)
	}
	if shareAmount != 0 {
		if total := liveroomdao.GetGuildIncomeTotal(guildId); total != nil {
			total.AddSettlementShareAmount(shareAmount)
		}
	}
	if hasDaily {
		liveroomdao.MarkDailyGuildEffectiveLivesSettled(dailyRows)
	}
	_ = entity.NewGuildIncomeSettlementLog(guildId, &snap, shareAmount, guildSharePercent)
}
