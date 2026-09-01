package liveroom

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/dao/anchorsalarycfgdao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/entity/live"
	"xr-game-server/gameevent"
	"xr-game-server/module/liverevenuesharecfg"
	"xr-game-server/module/wallet"
)

// writeAnchorSettlementUsd 是否把结算美金写入主播账面.
// 当前由工会账面归集、工会长分配;待主播转账功能上线后再打开.
const writeAnchorSettlementUsd = false

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
	liveroomdao.ResetGuildWeeklyAnchorSalary()
	g.Log().Infof(ctx, "anchor weekly settlement start, rooms=%d, salaryCfgs=%d, writeUsd=%v", len(rooms), len(cfgs), writeAnchorSettlementUsd)
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
	salary := matchAnchorSalaryAmount(countAnchorWeeklyWorkDays(dailyRows), dailyRows, cfgs)
	hasDaily := len(dailyRows) > 0
	hasUnsettled := !unsettled.IsZero()
	if !hasDaily && !hasUnsettled && salary == 0 {
		return
	}

	snap := unsettled.SnapshotAndClear()
	anchorSharePercent := liverevenuesharecfg.ResolveAnchorSharePercent()
	flowCommission := liverevenuesharecfg.CalcSettlementShareAmount(0, snap.TotalIncome)
	flowCommissionUsd := wallet.CalcDiamondToUsd(flowCommission)
	shareAmount := liverevenuesharecfg.CalcSettlementShareAmount(salary, snap.TotalIncome)
	shareAmountUsd := wallet.CalcDiamondToUsd(shareAmount)

	settled := liveroomdao.GetLiveRoomIncomeSettled(roomId)
	if settled != nil {
		settled.AddAmounts(&snap)
		settled.AddSettlementSalary(salary)
		settled.AddSettlementShareAmount(shareAmount)
		if writeAnchorSettlementUsd {
			settled.AddSettlementShareAmountUsd(shareAmountUsd)
		}
	}
	if salary != 0 {
		if total := liveroomdao.GetLiveRoomIncomeTotal(roomId); total != nil {
			total.AddSettlementSalary(salary)
		}
		// 开播底薪同步到工会(含折算美金可收)
		liveroomdao.MirrorGuildAnchorSettlementSalary(roomId, salary)
	}
	if shareAmount != 0 {
		if total := liveroomdao.GetLiveRoomIncomeTotal(roomId); total != nil {
			total.AddSettlementShareAmount(shareAmount)
			if writeAnchorSettlementUsd {
				total.AddSettlementShareAmountUsd(shareAmountUsd)
			}
		}
	}
	// 主播流水分佣美金入工会账面,由工会长分配(不写主播结算美金)
	if flowCommissionUsd != 0 {
		liveroomdao.MirrorGuildAnchorSettlementShareAmountUsd(roomId, flowCommissionUsd)
	}
	if hasDaily {
		liveroomdao.MarkDailyEffectiveLivesSettled(dailyRows)
	}

	logUsd := float64(0)
	if writeAnchorSettlementUsd {
		logUsd = shareAmountUsd
	}
	_ = entity.NewAnchorIncomeSettlementLog(roomId, &snap, salary, shareAmount, logUsd, anchorSharePercent)
}

// matchAnchorSalaryAmount 按薪资降序取最高满足档(结算规则后续按日表时长/workDays完善)
func matchAnchorSalaryAmount(weeklyWorkDays uint64, dailyRows []*entity.DailyAnchorEffectiveLive, cfgs []*entity.AnchorSalaryCfg) float64 {
	for _, cfg := range cfgs {
		if cfg == nil {
			continue
		}
		if weeklyWorkDays < cfg.WeeklyWorkDays {
			continue
		}
		if !dailyEffectiveLivesMeet(dailyRows, cfg.DailyLiveDurationMinutes) {
			continue
		}
		return cfg.SalaryAmount
	}
	return 0
}

func dailyEffectiveLivesMeet(dailyRows []*entity.DailyAnchorEffectiveLive, dailyNeedMinutes uint64) bool {
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

func countAnchorWeeklyWorkDays(dailyRows []*entity.DailyAnchorEffectiveLive) uint64 {
	var n uint64
	for _, row := range dailyRows {
		if row != nil && row.LiveDuration > 0 {
			n++
		}
	}
	return n
}
