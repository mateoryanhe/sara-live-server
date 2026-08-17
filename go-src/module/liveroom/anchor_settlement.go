package liveroom

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/dao/anchorsalarycfgdao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/entity/live"
	"xr-game-server/gameevent"
)

func initAnchorSettlement() {
	event.Sub(gameevent.WeekEvent, onWeekAnchorSettlement)
}

func onWeekAnchorSettlement(_ any) {
	settleOnShelfAnchors()
}

// settleOnShelfAnchors 周一0点:结算全部上架主播薪资+未结算收益
func settleOnShelfAnchors() {
	cfgs := anchorsalarycfgdao.ListAllOrderBySalaryDesc()
	rooms := liveroomdao.GetAllLiveRoom()
	ctx := gctx.New()
	g.Log().Infof(ctx, "anchor weekly settlement start, rooms=%d, salaryCfgs=%d", len(rooms), len(cfgs))
	for _, room := range rooms {
		if room == nil || room.ID == 0 {
			continue
		}
		settleOneAnchor(room.ID, cfgs)
	}
	g.Log().Infof(ctx, "anchor weekly settlement done")
}

func settleOneAnchor(roomId uint64, cfgs []*entity.AnchorSalaryCfg) {
	dailyRows := liveroomdao.ListRecentUnsettledDailyEffectiveLives(roomId)
	unsettled := liveroomdao.GetLiveRoomIncomeUnsettled(roomId)
	if unsettled == nil {
		return
	}
	salary := matchAnchorSalaryAmount(unsettled.EffectiveLiveCount, dailyRows, cfgs)
	hasDaily := len(dailyRows) > 0
	hasUnsettled := !unsettled.IsZero()
	if !hasDaily && !hasUnsettled && salary == 0 {
		return
	}

	snap := unsettled.SnapshotAndClear()
	settled := liveroomdao.GetLiveRoomIncomeSettled(roomId)
	if settled != nil {
		settled.AddAmounts(&snap)
		settled.AddSettlementSalary(salary)
	}
	if salary != 0 {
		if total := liveroomdao.GetLiveRoomIncomeTotal(roomId); total != nil {
			total.AddSettlementSalary(salary)
		}
	}
	if hasDaily {
		liveroomdao.MarkDailyEffectiveLivesSettled(dailyRows)
	}
	_ = entity.NewAnchorIncomeSettlementLog(roomId, &snap, salary)
}

// matchAnchorSalaryAmount 按薪资降序取最高满足档;日表任一天未达日门槛则该档不达标;周次数用未结算 EffectiveLiveCount
func matchAnchorSalaryAmount(weeklyEffectiveLiveCount uint64, dailyRows []*entity.DailyAnchorEffectiveLive, cfgs []*entity.AnchorSalaryCfg) float64 {
	for _, cfg := range cfgs {
		if cfg == nil {
			continue
		}
		if weeklyEffectiveLiveCount < cfg.WeeklyEffectiveLiveCount {
			continue
		}
		if !dailyEffectiveLivesMeet(dailyRows, cfg.DailyEffectiveLiveCount) {
			continue
		}
		return cfg.SalaryAmount
	}
	return 0
}

func dailyEffectiveLivesMeet(dailyRows []*entity.DailyAnchorEffectiveLive, dailyNeed uint64) bool {
	if dailyNeed == 0 {
		return true
	}
	if len(dailyRows) == 0 {
		return false
	}
	for _, row := range dailyRows {
		if row == nil || row.EffectiveLiveCount < dailyNeed {
			return false
		}
	}
	return true
}
