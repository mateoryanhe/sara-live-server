package stat

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/core/event"
	"xr-game-server/dao/statdao"
	"xr-game-server/entity/recharge"
	statentity "xr-game-server/entity/stat"
	"xr-game-server/gameevent"
)

func initRechargeEvent() {
	event.Sub(gameevent.RechargeArrivedEvent, onRechargeArrivedEvent)
}

func onRechargeArrivedEvent(data any) {
	order, ok := data.(*entity.RechargeOrder)
	if !ok || order == nil {
		g.Log().Errorf(gctx.New(), "RechargeArrivedEvent payload type error: %T", data)
		return
	}
	if order.Price <= 0 {
		return
	}
	if !shouldCountUserStat(order.UserId) {
		return
	}

	statAt := order.PaidAt
	if statAt.IsZero() {
		statAt = time.Now()
	}

	lockName := "stat_recharge"
	gmlock.Lock(lockName)
	defer gmlock.Unlock(lockName)

	if stat := statdao.GetSysStat(); stat != nil {
		stat.AddTotalRecharge(order.Price)
	}
	recordPeriodRecharge(statAt, order.Price)
	recordPeriodRechargeUser(statAt, order.UserId)
}

func recordPeriodRecharge(statAt time.Time, amount float64) {
	daily := statdao.GetDailyLoginStatByDate(statentity.FormatDailyLoginStatDate(statAt))
	daily.AddRechargeAmount(amount)

	weekly := statdao.GetWeeklyLoginStatByWeek(statentity.FormatWeeklyLoginStatKey(statAt))
	weekly.AddRechargeAmount(amount)

	monthly := statdao.GetMonthlyLoginStatByMonth(statentity.FormatMonthlyLoginStatKey(statAt))
	monthly.AddRechargeAmount(amount)
}

func recordPeriodRechargeUser(statAt time.Time, userId uint64) {
	if !shouldCountUserStat(userId) {
		return
	}
	date := statentity.FormatDailyLoginStatDate(statAt)
	if statdao.TryRecordDailyRecharge(date, userId) {
		statdao.GetDailyLoginStatByDate(date).AddRechargeUserCount(1)
	}
	week := statentity.FormatWeeklyLoginStatKey(statAt)
	if statdao.TryRecordWeeklyRecharge(week, userId) {
		statdao.GetWeeklyLoginStatByWeek(week).AddRechargeUserCount(1)
	}
	month := statentity.FormatMonthlyLoginStatKey(statAt)
	if statdao.TryRecordMonthlyRecharge(month, userId) {
		statdao.GetMonthlyLoginStatByMonth(month).AddRechargeUserCount(1)
	}
}
