package liveroomdao

import (
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/live"
)

// TryIncrementEffectiveLiveDays 下播且本场已计入日表后:若当日累计有效时长首次达到配置门槛,未结算有效开播天数 +1.
func TryIncrementEffectiveLiveDays(roomId uint64, at time.Time, sessionSec float64) {
	needMin := cfgdao.EffectiveLiveDailyAccumulatedMinutes()
	if roomId == 0 || sessionSec <= 0 || needMin <= 0 {
		return
	}
	needSec := float64(needMin * 60)
	date := entity.FormatDailyAnchorEffectiveLiveDate(at)
	row := GetDailyAnchorEffectiveLive(date, roomId)
	if row == nil {
		return
	}
	after := row.LiveDuration
	before := after - sessionSec
	if !crossedDailyAccumulatedLiveThreshold(before, after, needSec) {
		return
	}
	unsettled := GetLiveRoomIncomeUnsettled(roomId)
	if unsettled == nil {
		return
	}
	unsettled.AddEffectiveLiveDay()
}

func crossedDailyAccumulatedLiveThreshold(beforeSec, afterSec, needSec float64) bool {
	if needSec <= 0 {
		return false
	}
	return beforeSec < needSec && afterSec >= needSec
}
