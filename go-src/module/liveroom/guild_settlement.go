package liveroom

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/guildsalarycfgdao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/entity/live"
	"xr-game-server/gameevent"
)

func initGuildSettlement() {
	event.Sub(gameevent.WeekEvent, onWeekGuildSettlement)
}

func onWeekGuildSettlement(_ any) {
	settleOnShelfGuilds()
}

// settleOnShelfGuilds 周一0点:结算全部上架工会薪资+未结算收益
func settleOnShelfGuilds() {
	cfgs := guildsalarycfgdao.ListAllOrderBySalaryDesc()
	guilds := guilddao.ListOnShelfGuilds()
	ctx := gctx.New()
	g.Log().Infof(ctx, "guild weekly settlement start, guilds=%d, salaryCfgs=%d", len(guilds), len(cfgs))
	for _, guild := range guilds {
		if guild == nil || guild.ID == 0 {
			continue
		}
		settleOneGuild(guild.ID, cfgs)
	}
	g.Log().Infof(ctx, "guild weekly settlement done")
}

func settleOneGuild(guildId uint64, cfgs []*entity.GuildSalaryCfg) {
	dailyRows := liveroomdao.ListRecentUnsettledDailyGuildEffectiveLives(guildId)
	unsettled := liveroomdao.GetGuildIncomeUnsettled(guildId)
	if unsettled == nil {
		return
	}
	salary := matchGuildSalaryAmount(countGuildWeeklyWorkDays(dailyRows), dailyRows, cfgs)
	hasDaily := len(dailyRows) > 0
	hasUnsettled := !unsettled.IsZero()
	if !hasDaily && !hasUnsettled && salary == 0 {
		return
	}

	snap := unsettled.SnapshotAndClear()
	settled := liveroomdao.GetGuildIncomeSettled(guildId)
	if settled != nil {
		settled.AddAmounts(&snap)
		settled.AddSettlementSalary(salary)
	}
	if salary != 0 {
		if total := liveroomdao.GetGuildIncomeTotal(guildId); total != nil {
			total.AddSettlementSalary(salary)
		}
	}
	if hasDaily {
		liveroomdao.MarkDailyGuildEffectiveLivesSettled(dailyRows)
	}
	_ = entity.NewGuildIncomeSettlementLog(guildId, &snap, salary)
}

// matchGuildSalaryAmount 按薪资降序取最高满足档(结算规则后续完善)
func matchGuildSalaryAmount(weeklyWorkDays uint64, dailyRows []*entity.DailyGuildEffectiveLive, cfgs []*entity.GuildSalaryCfg) float64 {
	for _, cfg := range cfgs {
		if cfg == nil {
			continue
		}
		if weeklyWorkDays < cfg.WeeklyWorkDays {
			continue
		}
		if !dailyGuildEffectiveLivesMeet(dailyRows, cfg.DailyLiveDurationMinutes) {
			continue
		}
		return cfg.SalaryAmount
	}
	return 0
}

func dailyGuildEffectiveLivesMeet(dailyRows []*entity.DailyGuildEffectiveLive, dailyNeedMinutes uint64) bool {
	if dailyNeedMinutes == 0 {
		return true
	}
	if len(dailyRows) == 0 {
		return false
	}
	needSec := float64(dailyNeedMinutes) * 60
	for _, row := range dailyRows {
		if row == nil || row.LiveDuration < needSec {
			return false
		}
	}
	return true
}

func countGuildWeeklyWorkDays(dailyRows []*entity.DailyGuildEffectiveLive) uint64 {
	var n uint64
	for _, row := range dailyRows {
		if row != nil && row.LiveDuration > 0 {
			n++
		}
	}
	return n
}
