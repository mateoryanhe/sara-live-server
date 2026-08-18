package liveroom

import (
	"strconv"

	"xr-game-server/constants/cmd"
	"xr-game-server/core/push"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/liveroomdto"
	liveentity "xr-game-server/entity/live"
)

// NotifyLiveRecordTotalIncome 向房间内全体在线用户(含主播)广播本场直播总收益
func NotifyLiveRecordTotalIncome(room *liveentity.LiveRoom) {
	if room == nil || room.LiveRecordId == 0 {
		return
	}
	liveRecord := liveroomdao.GetLiveRecordById(room.LiveRecordId)
	if liveRecord == nil {
		return
	}
	broadcastLiveRecordTotalIncome(room.ID, room.LiveRecordId, liveRecord.TotalIncome)
}

func broadcastLiveRecordTotalIncome(roomId, liveRecordId uint64, totalIncome float64) {
	if roomId == 0 || liveRecordId == 0 {
		return
	}
	payload := &liveroomdto.LiveRecordTotalIncomePushItem{
		RoomId:       strconv.FormatUint(roomId, 10),
		LiveRecordId: strconv.FormatUint(liveRecordId, 10),
		TotalIncome:  totalIncome,
	}
	for _, userId := range getOnline(roomId) {
		push.Data(userId, cmd.LiveRoomTotalIncome, payload)
	}
	push.Data(roomId, cmd.LiveRoomTotalIncome, payload)
}
