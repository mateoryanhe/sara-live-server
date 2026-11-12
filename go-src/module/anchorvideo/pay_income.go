package anchorvideo

import (
	"time"

	"xr-game-server/dao/liveroomdao"
)

// RecordAnchorVideoPayIncome 主播相册视频解锁收益:写入未结算/生涯/工会,并累计进社交流水
func RecordAnchorVideoPayIncome(anchorId uint64, amount float64) {
	if anchorId == 0 || amount <= 0 {
		return
	}
	room := liveroomdao.ResolveRoom(anchorId)
	if room == nil {
		return
	}
	liveroomdao.GetLiveRoomIncomeUnsettled(room.ID).AddAnchorVideoEarn(amount)
	liveroomdao.GetLiveRoomIncomeTotal(room.ID).AddAnchorVideoEarn(amount)
	liveroomdao.MirrorGuildAnchorVideoEarn(room.ID, amount)
	liveroomdao.MirrorDailyAnchorAnchorVideoEarn(room.ID, time.Now(), amount)
}
