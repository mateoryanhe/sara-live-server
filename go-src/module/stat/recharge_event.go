package stat

import (
	"time"

	"xr-game-server/core/event"
	"xr-game-server/dao/statdao"
	statentity "xr-game-server/entity/stat"
	"xr-game-server/gameevent"
)

func initRechargeEvent() {
	event.Sub(gameevent.UsdIncomeArrivedEvent, onUsdIncomeArrivedEvent)
}

func onUsdIncomeArrivedEvent(data any) {
	payload, ok := data.(*gameevent.UsdIncomeArrivedEventData)
	if !ok || payload == nil || !payload.Granted {
		return
	}
	enqueue(&statJob{Kind: jobRecharge, Payload: data})
}

func consumeRechargeJob(job *statJob) {
	data, ok := job.Payload.(*gameevent.UsdIncomeArrivedEventData)
	if !ok || data == nil || data.Order == nil {
		return
	}
	if !data.Granted {
		return
	}
	order := data.Order
	amount := data.UsdAmount
	if amount <= 0 {
		amount = order.Price
	}
	if amount <= 0 {
		return
	}

	// 充值白名单 → 虚拟美金累计,不进真实充值报表
	if data.IsVirtual() {
		if stat := statdao.GetSysStat(); stat != nil {
			stat.AddTotalVirtualRecharge(amount)
		}
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
		stat.AddTotalRecharge(amount)
	}
	recordPeriodRecharge(statAt, amount)
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
