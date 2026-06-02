package stat

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/core/event"
	"xr-game-server/dao/dailyloginstatdao"
	"xr-game-server/dao/dailyuserrechargdao"
	"xr-game-server/dao/monthlyloginstatdao"
	"xr-game-server/dao/monthlyuserrechargdao"
	"xr-game-server/dao/statdao"
	"xr-game-server/dao/weeklyloginstatdao"
	"xr-game-server/dao/weeklyuserrechargdao"
	"xr-game-server/entity"
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
	daily := dailyloginstatdao.GetByDate(entity.FormatDailyLoginStatDate(statAt))
	daily.AddRechargeAmount(amount)

	weekly := weeklyloginstatdao.GetByWeek(entity.FormatWeeklyLoginStatKey(statAt))
	weekly.AddRechargeAmount(amount)

	monthly := monthlyloginstatdao.GetByMonth(entity.FormatMonthlyLoginStatKey(statAt))
	monthly.AddRechargeAmount(amount)
}

func recordPeriodRechargeUser(statAt time.Time, userId uint64) {
	if userId == 0 {
		return
	}
	date := entity.FormatDailyLoginStatDate(statAt)
	if dailyuserrechargdao.TryRecordRecharge(date, userId) {
		dailyloginstatdao.GetByDate(date).AddRechargeUserCount(1)
	}
	week := entity.FormatWeeklyLoginStatKey(statAt)
	if weeklyuserrechargdao.TryRecordRecharge(week, userId) {
		weeklyloginstatdao.GetByWeek(week).AddRechargeUserCount(1)
	}
	month := entity.FormatMonthlyLoginStatKey(statAt)
	if monthlyuserrechargdao.TryRecordRecharge(month, userId) {
		monthlyloginstatdao.GetByMonth(month).AddRechargeUserCount(1)
	}
}
