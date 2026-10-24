package stat

import (
	"time"

	"xr-game-server/core/event"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/statdao"
	"xr-game-server/entity/stat"
	"xr-game-server/gameevent"
	"xr-game-server/module/auth"
)

func initLoginEvent() {
	event.Sub(gameevent.LoginEvent, onLoginEvent)
}

func onLoginEvent(data any) {
	enqueue(&statJob{Kind: jobLogin, At: time.Now(), Payload: data})
}

func consumeLoginJob(job *statJob) {
	userId, ok := job.Payload.(uint64)
	if !ok || userId == 0 {
		return
	}
	now := job.At
	if now.IsZero() {
		now = time.Now()
	}
	clientChannel := accountClientLoginChannel(userId)
	recordDailyLogin(now, userId, clientChannel)
	recordWeeklyLogin(now, userId, clientChannel)
	recordMonthlyLogin(now, userId, clientChannel)
}

// accountClientLoginChannel 返回 App/H5 设备码渠道;其它渠道返回 0
func accountClientLoginChannel(userId uint64) uint {
	acc := accountdao.GetAccountById(userId)
	if acc == nil {
		return 0
	}
	switch acc.Channel {
	case auth.DeviceChannel, auth.H5DeviceChannel:
		return acc.Channel
	default:
		return 0
	}
}

func recordDailyLogin(now time.Time, userId uint64, clientChannel uint) {
	date := entity.FormatDailyLoginStatDate(now)
	if !statdao.TryRecordDailyLogin(date, userId) {
		return
	}
	stat := statdao.GetDailyLoginStatByDate(date)
	stat.AddLoginCount(1)
	switch clientChannel {
	case auth.DeviceChannel:
		stat.AddAppDeviceLoginCount(1)
	case auth.H5DeviceChannel:
		stat.AddH5DeviceLoginCount(1)
	}
	statdao.PublishDailyLoginStat(stat)
}

func recordWeeklyLogin(now time.Time, userId uint64, clientChannel uint) {
	week := entity.FormatWeeklyLoginStatKey(now)
	if !statdao.TryRecordWeeklyLogin(week, userId) {
		return
	}
	stat := statdao.GetWeeklyLoginStatByWeek(week)
	stat.AddLoginCount(1)
	switch clientChannel {
	case auth.DeviceChannel:
		stat.AddAppDeviceLoginCount(1)
	case auth.H5DeviceChannel:
		stat.AddH5DeviceLoginCount(1)
	}
	statdao.PublishWeeklyLoginStat(stat)
}

func recordMonthlyLogin(now time.Time, userId uint64, clientChannel uint) {
	month := entity.FormatMonthlyLoginStatKey(now)
	if !statdao.TryRecordMonthlyLogin(month, userId) {
		return
	}
	stat := statdao.GetMonthlyLoginStatByMonth(month)
	stat.AddLoginCount(1)
	switch clientChannel {
	case auth.DeviceChannel:
		stat.AddAppDeviceLoginCount(1)
	case auth.H5DeviceChannel:
		stat.AddH5DeviceLoginCount(1)
	}
	statdao.PublishMonthlyLoginStat(stat)
}
