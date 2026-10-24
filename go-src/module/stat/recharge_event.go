package stat

import (
	"time"

	"xr-game-server/core/event"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/statdao"
	rechargeentity "xr-game-server/entity/recharge"
	statentity "xr-game-server/entity/stat"
	"xr-game-server/gameevent"
	"xr-game-server/module/auth"
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

	isCoinMerchant := isCoinMerchantUsdIncome(data)

	// 真实美金先计入不区分业务类型的总额，再分别累计普通用户与币商入账。
	if stat := statdao.GetSysStat(); stat != nil {
		stat.AddTotalRecharge(amount)
		if isCoinMerchant {
			stat.AddTotalCoinMerchantRecharge(amount)
		} else {
			stat.AddTotalNormalUserRecharge(amount)
		}
	}

	statAt := order.PaidAt
	if statAt.IsZero() {
		statAt = time.Now()
	}
	recordPeriodRecharge(statAt, amount, isCoinMerchant)
	recordPeriodClientRecharge(statAt, order.UserId, amount)
	recordDailyPromotionChannelRechargeSuccess(statAt, order.UserId, amount)

	if !shouldCountUserStat(order.UserId) {
		return
	}
	recordPeriodRechargeUser(statAt, order.UserId)
}

func isCoinMerchantUsdIncome(data *gameevent.UsdIncomeArrivedEventData) bool {
	if data == nil {
		return false
	}
	if data.Kind == gameevent.UsdIncomeKindCoinMerchant {
		return true
	}
	return data.Order != nil && data.Order.PayChannel == rechargeentity.RechargeCfgTypeCoinMerchant
}

func recordPeriodRecharge(statAt time.Time, amount float64, isCoinMerchant bool) {
	daily := statdao.GetDailyLoginStatByDate(statentity.FormatDailyLoginStatDate(statAt))
	daily.AddRechargeAmount(amount)
	if isCoinMerchant {
		daily.AddCoinMerchantRechargeAmount(amount)
	} else {
		daily.AddNormalUserRechargeAmount(amount)
	}

	weekly := statdao.GetWeeklyLoginStatByWeek(statentity.FormatWeeklyLoginStatKey(statAt))
	weekly.AddRechargeAmount(amount)
	if isCoinMerchant {
		weekly.AddCoinMerchantRechargeAmount(amount)
	} else {
		weekly.AddNormalUserRechargeAmount(amount)
	}

	monthly := statdao.GetMonthlyLoginStatByMonth(statentity.FormatMonthlyLoginStatKey(statAt))
	monthly.AddRechargeAmount(amount)
	if isCoinMerchant {
		monthly.AddCoinMerchantRechargeAmount(amount)
	} else {
		monthly.AddNormalUserRechargeAmount(amount)
	}
}

// recordPeriodClientRecharge 按账号 channel 累计 App/H5 设备码用户充值(USD)
func recordPeriodClientRecharge(statAt time.Time, userId uint64, amount float64) {
	if userId == 0 || amount <= 0 {
		return
	}
	acc := accountdao.GetAccountById(userId)
	if acc == nil {
		return
	}
	switch acc.Channel {
	case auth.DeviceChannel:
		daily := statdao.GetDailyLoginStatByDate(statentity.FormatDailyLoginStatDate(statAt))
		daily.AddAppDeviceRechargeAmount(amount)
		weekly := statdao.GetWeeklyLoginStatByWeek(statentity.FormatWeeklyLoginStatKey(statAt))
		weekly.AddAppDeviceRechargeAmount(amount)
		monthly := statdao.GetMonthlyLoginStatByMonth(statentity.FormatMonthlyLoginStatKey(statAt))
		monthly.AddAppDeviceRechargeAmount(amount)
	case auth.H5DeviceChannel:
		daily := statdao.GetDailyLoginStatByDate(statentity.FormatDailyLoginStatDate(statAt))
		daily.AddH5DeviceRechargeAmount(amount)
		weekly := statdao.GetWeeklyLoginStatByWeek(statentity.FormatWeeklyLoginStatKey(statAt))
		weekly.AddH5DeviceRechargeAmount(amount)
		monthly := statdao.GetMonthlyLoginStatByMonth(statentity.FormatMonthlyLoginStatKey(statAt))
		monthly.AddH5DeviceRechargeAmount(amount)
	}
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
