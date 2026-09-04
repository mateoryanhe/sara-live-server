package stat

import (
	"time"

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
	enqueue(&statJob{Kind: jobRecharge, Payload: data})
}

func consumeRechargeJob(job *statJob) {
	order, ok := job.Payload.(*entity.RechargeOrder)
	if !ok || order == nil {
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
		dailyStat := statdao.GetDailyLoginStatByDate(date)
		dailyStat.AddRechargeUserCount(1)
		statdao.PublishDailyLoginStat(dailyStat)
	}
	week := statentity.FormatWeeklyLoginStatKey(statAt)
	if statdao.TryRecordWeeklyRecharge(week, userId) {
		weeklyStat := statdao.GetWeeklyLoginStatByWeek(week)
		weeklyStat.AddRechargeUserCount(1)
		statdao.PublishWeeklyLoginStat(weeklyStat)
	}
	month := statentity.FormatMonthlyLoginStatKey(statAt)
	if statdao.TryRecordMonthlyRecharge(month, userId) {
		monthlyStat := statdao.GetMonthlyLoginStatByMonth(month)
		monthlyStat.AddRechargeUserCount(1)
		statdao.PublishMonthlyLoginStat(monthlyStat)
	}
}
