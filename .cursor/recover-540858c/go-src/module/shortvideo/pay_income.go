package shortvideo

import (
	"time"

	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
)

// RecordShortVideoPayIncome 短视频付费观看作者收益:主播写入结算体系,非主播写入 user_ext(暂不结算)
func RecordShortVideoPayIncome(authorId uint64, amount float64) {
	if authorId == 0 || amount <= 0 {
		return
	}
	author := userinfodao.GetUserInfoByUserId(authorId)
	if author == nil {
		return
	}
	if !author.IsAnchor() {
		recordNonAnchorShortVideoIncome(authorId, amount)
		return
	}
	room := liveroomdao.ResolveRoom(authorId)
	if room == nil {
		recordNonAnchorShortVideoIncome(authorId, amount)
		return
	}
	liveroomdao.GetLiveRoomIncomeUnsettled(room.ID).AddShortVideoEarn(amount)
	liveroomdao.GetLiveRoomIncomeTotal(room.ID).AddShortVideoEarn(amount)
	liveroomdao.MirrorGuildShortVideoEarn(room.ID, amount)
	liveroomdao.MirrorDailyAnchorShortVideoEarn(room.ID, time.Now(), amount)
}

func recordNonAnchorShortVideoIncome(userId uint64, amount float64) {
	ext := userinfodao.GetUserExtByUserId(userId)
	if ext == nil {
		return
	}
	ext.AddShortVideoUnsettledIncome(amount)
	userinfodao.PublishUserExt(ext)
}
