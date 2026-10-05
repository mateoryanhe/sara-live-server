package game

import (
	stdmath "math"
	"time"

	"xr-game-server/dao/liveroomdao"
	"xr-game-server/module/wallet"
)

// RecordLiveRoomGameBetIncome 直播间内游戏下注收益:写入结算体系并累加单场直播游戏消费
func RecordLiveRoomGameBetIncome(liveRoomId, liveRecordId uint64, goldAmount float64) {
	if liveRoomId == 0 || goldAmount <= 0 || stdmath.IsNaN(goldAmount) || stdmath.IsInf(goldAmount, 0) {
		return
	}
	// 游戏流水保留金币原值；展示用总收益按本次下注时的汇率折算为钻石。
	incomeDelta := wallet.CalcGoldToDiamond(goldAmount)
	liveroomdao.GetLiveRoomIncomeUnsettled(liveRoomId).AddGameEarn(goldAmount, incomeDelta)
	liveroomdao.GetLiveRoomIncomeTotal(liveRoomId).AddGameEarn(goldAmount, incomeDelta)
	liveroomdao.MirrorGuildGameEarn(liveRoomId, goldAmount, incomeDelta)
	liveroomdao.MirrorDailyAnchorGameEarn(liveRoomId, time.Now(), goldAmount, incomeDelta)
	if liveRecordId == 0 {
		return
	}
	if liveRecord := liveroomdao.GetLiveRecordById(liveRecordId); liveRecord != nil {
		liveRecord.AddGameEarn(goldAmount, incomeDelta)
		liveroomdao.PublishLiveRecord(liveRecord)
	}
}
