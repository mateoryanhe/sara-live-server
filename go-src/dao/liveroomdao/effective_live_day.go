package liveroomdao

import (
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/live"
)

// TryIncrementEffectiveLiveDays 下播且本场已计入日表后:日表未结算且本笔下播使当日累计有效时长首次达到门槛时,未结算有效开播天数 +1.
func TryIncrementEffectiveLiveDays(roomId uint64, at time.Time, sessionSec float64) {
	needMin := cfgdao.EffectiveLiveDailyAccumulatedMinutes()
	if roomId == 0 || sessionSec <= 0 || needMin <= 0 {
		return
	}
	needSec := float64(needMin * 60)
	date := entity.FormatDailyAnchorEffectiveLiveDate(at)
	row := GetDailyAnchorEffectiveLive(date, roomId)
	if row == nil || row.Settled {
		return
	}
	after := row.LiveDuration
	before := after - sessionSec
	if !firstReachDailyLiveThreshold(before, after, needSec) {
		return
	}
	unsettled := GetLiveRoomIncomeUnsettled(roomId)
	if unsettled == nil {
		return
	}
	unsettled.AddEffectiveLiveDay()
}

// firstReachDailyLiveThreshold 本笔下播前未达标、下播后达标(同一天更早下播已达标则 before 已 >= needSec,不再 +1).
func firstReachDailyLiveThreshold(beforeSec, afterSec, needSec float64) bool {
	if needSec <= 0 {
		return false
	}
	return beforeSec < needSec && afterSec >= needSec
}
